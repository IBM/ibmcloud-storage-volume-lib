/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/
package utils

import (
	"fmt"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"go.uber.org/zap"
)

func GetStorageID(sess backend.Session, volumeType string, orderID int, logger *zap.Logger, config *config.SoftlayerConfig) (int, error) {
	// Step 1- Create an inline method for retry main method, which will return (bool, error)
	// bool is for retry termination before exhausting max attempt and error is for final issue in retry
	var storageID int
	getStorageIDTempFun := func() (bool, error) {
		storage, err := GetNetworkStorageFromOrderID(sess, volumeType, orderID, logger)
		if err == nil && storage != nil {
			storageID = *storage.Id
			return true, nil
		}
		return false, err

	}

	// Step 2- Calling retry method
	if err := ProvisioningRetry(getStorageIDTempFun, logger, config.SoftlayerVolProvisionTimeout, config.SoftlayerRetryInterval); err != nil {
		return 0, err
	}
	return storageID, nil
}

//Retrieves storageID of given orderID.
func GetNetworkStorageFromOrderID(bkSession backend.Session, volumeType string, orderIDIn int, logger *zap.Logger) (*datatypes.Network_Storage, error) {

	storageMask := "id,username,notes,billingItem.orderItem.order.id"

	slFilters := filter.New()
	slFilters = append(slFilters, filter.Path("networkStorage.nasType").Eq(volumeType))
	slFilters = append(slFilters, filter.Path("networkStorage.billingItem.orderItem.order.id").Eq(orderIDIn)) // Do not fetch cancelled volumes

	logger.Info("Filterused ", zap.Reflect("slFilters", slFilters))
	accountService := bkSession.GetAccountService().Mask(storageMask)
	accountService = accountService.Filter(slFilters.Build())
	storages, err := accountService.GetNetworkStorage()
	logger.Info("Volumes found", zap.Reflect("Volumes", storages))
	if err != nil {
		return nil, err
	}
	switch len(storages) {
	case 0:
		return nil, fmt.Errorf("unable to find network storage associated with order %d", orderIDIn)
	case 1:
		// double check if correct storage is found by matching requestID and fouund orderID
		orderID := *storages[0].BillingItem.OrderItem.Order.Id
		if orderID == orderIDIn {
			return &storages[0], nil
		} else {
			logger.Error("Incorrect storage found", zap.Int("requestID", orderIDIn), zap.Reflect("storage", storages[0]))
			return nil, fmt.Errorf("Incorrect storage found %d associated with order %d", orderID, orderIDIn)
		}
	default:
		return nil, fmt.Errorf("multiple storage volumes associated with order %d", orderIDIn)
	}
}

// Makes sure there are no active transactions for volume with given ID
func WaitForTransactionsComplete(sess backend.Session, storageID int, logger *zap.Logger, config *config.SoftlayerConfig) error {
	var txnCount uint
	//slbs.Logger.Info("Waiting for transactions to complete for storage..", zap.Int("StorageID", storageID))
	WaitForTransactionsComplete := func() (bool, error) {
		var err error
		nwStorage, err := sess.GetNetworkStorageService().ID(storageID).Mask("id,activeTransactionCount").GetObject()
		if err != nil {
			return false, err
		}

		txnCount = *nwStorage.ActiveTransactionCount
		if err != nil {
			return false, err
		}
		if txnCount > 0 {
			return false, err
		}
		return true, nil
	}

	return ProvisioningRetry(WaitForTransactionsComplete, logger, config.SoftlayerVolProvisionTimeout, config.SoftlayerRetryInterval)
}
