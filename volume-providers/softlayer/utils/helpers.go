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

	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend"
	"go.uber.org/zap"
)

func GetStorageID(sess backend.Session, orderID int, logger zap.Logger) (int, error) {
	// Step 1- Create an inline method for retry main method, which will return (bool, error)
	// bool is for retry termination before exhausting max attempt and error is for final issue in retry
	var storageID int
	getStorageIDTempFun := func() (bool, error) {
		id, err := GetNetworkStorageIDFromOrderID(sess, orderID, logger)
		if err != nil || id <= 0 {
			return false, err
		}
		storageID = id
		return true, nil
	}

	// Step 2- Calling retry method
	if err := ProvisioningRetry(getStorageIDTempFun); err != nil {
		return 0, err
	}
	return storageID, nil
}

//Retrieves storageID of given orderID.
func GetNetworkStorageIDFromOrderID(bkSession backend.Session, orderID int, logger zap.Logger) (int, error) {
	filter := fmt.Sprintf(`{"networkStorage":{"nasType":{"operation":"ISCSI"},
                        "billingItem":{"orderItem":{"order":{"id":{"operation":%d}}}
                        } } }`, orderID)

	accService := bkSession.GetAccountService()
	storage, err := accService.Filter(filter).GetNetworkStorage()
	if err != nil {
		return 0, err
	}
	switch len(storage) {
	case 0:
		return 0, fmt.Errorf("unable to find network storage associated with order %d", orderID)
	case 1:
		return *storage[0].Id, nil
	default:
		return 0, fmt.Errorf("multiple storage volumes associated with order %d", orderID)
	}
}

// Makes sure there are no active transactions for volume with given ID
func WaitForTransactionsComplete(sess backend.Session, storageID int) error {
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

	return ProvisioningRetry(WaitForTransactionsComplete)
}
