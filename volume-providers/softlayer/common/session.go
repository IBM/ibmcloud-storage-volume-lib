/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package common

import (
	"fmt"
	"strconv"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/utils"
	"github.com/softlayer/softlayer-go/filter"
	"go.uber.org/zap"
)

// SLSession implements lib.Session
type SLSession struct {
	SLAccountID        int
	Url                string
	Backend            backend.Session
	Logger             zap.Logger
	Config             *config.SoftlayerConfig
	ContextCredentials provider.ContextCredentials
	VolumeType         provider.VolumeType
	Provider           provider.VolumeProvider
}

// Close at present does nothing
func (*SLSession) Close() {
	// Do nothing for now
}

func (sls *SLSession) GetVolumeByRequestID(requestID string) (*provider.Volume, error) {
	/*filter := fmt.Sprintf(`{
					"networkStorage":{
									"nasType":{"operation":%s},
									"billingItem":{
													"orderItem":{"order":{"id":{"operation":%d}}}
									}
					}
	}`, sls.VolumeType, requestID)*/
	storageMask := "id,username,notes,billingItem.orderItem.order.id"

	slFilters := filter.New()
	slFilters = append(slFilters, filter.Path("networkStorage.nasType").Eq(sls.VolumeType))
	slFilters = append(slFilters, filter.Path("networkStorage.billingItem.orderItem.order.id").Eq(requestID)) // Do not fetch cancelled volumes

	sls.Logger.Info("Filterused ", zap.Reflect("slFilters", slFilters))
	accountService := sls.Backend.GetAccountService().Mask(storageMask)
	accountService = accountService.Filter(slFilters.Build())
	storages, err := accountService.GetNetworkStorage()
	sls.Logger.Info("Volumes found", zap.Reflect("Volumes", storages))
	if err != nil {
		return nil, err
	}
	switch len(storages) {
	case 0:
		return nil, fmt.Errorf("unable to find network storage associated with order %d", requestID)
	case 1:
		// double check if correct storage is found by matching requestID and fouund orderID
		orderID := strconv.Itoa(*storages[0].BillingItem.OrderItem.Order.Id)
		if orderID == requestID {
			vol := utils.ConvertToVolumeType(storages[0], sls.Logger, sls.Provider, sls.VolumeType)
			return vol, nil
		} else {
			sls.Logger.Error("Incorrect storage found", zap.String("requestID", requestID), zap.Reflect("storage", storages[0]))
			return nil, fmt.Errorf("Incorrect storage found %d associated with order %d", orderID, requestID)
		}
	default:
		return nil, fmt.Errorf("multiple storage volumes associated with order %d", requestID)
	}

}
