/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package softlayer_block

import (
	"strconv"
	"strings"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/messages"
	utils "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/utils"
	"go.uber.org/zap"
)

func (sls *SLBlockSession) OrderSnapshot(volumeRequest provider.Volume) error {
	// Step 1- validate input which are required
	sls.Logger.Info("Requested volume is:", zap.Reflect("Volume", volumeRequest))
	if volumeRequest.SnapshotSpace == nil {
		sls.Logger.Error("No proper input, please provide volume ID and snapshot space size")
		return messages.GetUserError("E0013", nil)
	}
	volid := utils.ToInt(volumeRequest.VolumeID)
	snapshotSize := *volumeRequest.SnapshotSpace
	if volid == 0 || snapshotSize == 0 {
		sls.Logger.Error("No proper input, please provide volume ID and snapshot space size")
		return messages.GetUserError("E0013", nil)
	}

	// Step 2- Get volume details
	mask := "id,billingItem[location,hourlyFlag],storageType[keyName],storageTierLevel,provisionedIops,staasVersion,hasEncryptionAtRest"
	storageObj := sls.Backend.GetNetworkStorageService()
	storage, err := storageObj.ID(volid).Mask(mask).GetObject()
	if err != nil {
		return messages.GetUserError("E0011", nil, volid, "Please check the volume id")
	}
	sls.Logger.Info("in OrderSnapshot Volum Object ---->", zap.Reflect("Volume", storage))

	// Step 3: verify original volume exists or not
	if storage.BillingItem == nil {
		return messages.GetUserError("E0014", nil, volid)
	}

	if storage.BillingItem.Location == nil || storage.BillingItem.Location.Id == nil {
		sls.Logger.Error("Original Volume does not have location ID", zap.Reflect("Location", storage.BillingItem.Location))
		return messages.GetUserError("E0024", nil, volid)
	}
	datacenterID := *storage.BillingItem.Location.Id

	// Step 4: Get billing item category code
	if storage.BillingItem.CategoryCode == nil {
		return messages.GetUserError("E0015", nil, volid)
	}
	billingItemCategoryCode := *storage.BillingItem.CategoryCode
	order_type_is_saas := true
	if billingItemCategoryCode == "storage_as_a_service" {
		order_type_is_saas = true
	} else if billingItemCategoryCode == "storage_service_enterprise" {
		order_type_is_saas = false
	} else {
		return messages.GetUserError("E0016", nil, volid)
	}

	// Step 5: Get the product package by using billing item category code
	packageDetails, errPackage := utils.GetPackageDetails(sls.Logger, sls.Backend, billingItemCategoryCode)
	if errPackage != nil {
		return messages.GetUserError("E0017", nil, billingItemCategoryCode)
	}
	finalPackageID := *packageDetails.Id

	// Step 6: Get required price for snapshot space as per volume type
	finalPrices := []datatypes.Product_Item_Price{}
	if order_type_is_saas {
		volume_storage_type := *storage.StorageType.KeyName
		if strings.Contains(volume_storage_type, "ENDURANCE") {
			volumeTier := utils.GetEnduranceTierIopsPerGB(sls.Logger, storage)
			finalPrices = []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSSnapshotSpacePrice(sls.Logger, packageDetails, snapshotSize, volumeTier, 0))},
			}
		} else if strings.Contains(volume_storage_type, "PERFORMANCE") {
			if !utils.IsVolumeCreatedWithStaaS(storage) {
				return messages.GetUserError("E0018", nil, volid)
			}
			iops := utils.ToInt(*storage.ProvisionedIops)
			finalPrices = []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSSnapshotSpacePrice(sls.Logger, packageDetails, snapshotSize, "", iops))},
			}
		} else {
			return messages.GetUserError("E0019", nil, volume_storage_type)
		}
	} else { // 'storage_service_enterprise' package
		volumeTier := utils.GetEnduranceTierIopsPerGB(sls.Logger, storage)
		finalPrices = []datatypes.Product_Item_Price{
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetEnterpriseSpacePrice(sls.Logger, packageDetails, "snapshot", snapshotSize, volumeTier))},
		}
	}
	/*
			if upgrade:
		        complex_type = 'SoftLayer_Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace_Upgrade'
		    else:
		        complex_type = 'SoftLayer_Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace'
	*/

	// Step 7: Create order
	cpo := datatypes.Container_Product_Order{
		ComplexType: sl.String("SoftLayer_Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace"),
		Quantity:    sl.Int(1),
		Location:    sl.String(strconv.Itoa(datacenterID)),
		PackageId:   sl.Int(finalPackageID),
		Prices:      finalPrices,
	}

	sp := &datatypes.Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace{
		VolumeId:                sl.Int(volid),
		Container_Product_Order: cpo,
	}
	sls.Logger.Info("Order deails ... ", zap.Reflect("OrderDeails", sp))
	/*orderContainer := &datatypes.Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace_Upgrade{
		Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace : sp1,
	}*/

	// Step 8: place order
	productOrderObj := sls.Backend.GetProductOrderService()
	snOrderID, snError := productOrderObj.PlaceOrder(sp, sl.Bool(false))
	if snError != nil {
		return messages.GetUserError("E0020", snError, volid, snapshotSize)
	}
	sls.Logger.Info("Successfully placed Snapshot order .... ", zap.Reflect("orderID", *snOrderID.OrderId), zap.Reflect("VolumeID", volid), zap.Reflect("Size", snapshotSize))
	sls.Logger.Info("Snapshot order details.... ", zap.Reflect("orderDetails", snOrderID))
	time.Sleep(300)
	sls.Logger.Info("Snapshot order details.... ", zap.Reflect("orderDetails", snOrderID))
	return nil
	// TODO: need to keep checking if order is ready or not
}
