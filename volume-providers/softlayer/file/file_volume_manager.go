/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package softlayer_file

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/utils"
	"go.uber.org/zap"
)

// Type returns the underlying volume type
func (sls *SLFileSession) Type() provider.VolumeType {
	return VolumeTypeFile
}

var ENDURANCE_TIERS = map[string]int{
	"0.25": 100,
	"2":    200,
	"4":    300,
	"10":   1000,
}

var (
	VOLUME_DELETE_REASON = "deleted by ibm-volume-lib on behalf of user request"
)

//Creates Volume along with snapshot space allocation
//TESTED: SAAS offering for endurance and Performance
//TODO: test Enterprise volumeOrdering with endurance and performance
func (sls *SLFileSession) VolumeCreate(volumeRequest provider.Volume) (*provider.Volume, error) {
	sls.Logger.Info("Creating volume as per order request .... ", zap.Reflect("Volume", volumeRequest))

	//Hard coded required values for testing TODO: Remove them and pass the values as arguments for function
	service_offering := "storage_as_a_service"
	volume_type := string(volumeRequest.VolumeType)
	hourly_billing_flag := false

	// Step 1- Validate inputs
	if volumeRequest.ProviderType != "performance" && volumeRequest.ProviderType != "endurance" {
		return nil, messages.GetUserError("E0012", nil)
	}
	datacenterID, dcError := utils.GetDataCenterID(sls.Logger, sls.Backend, volumeRequest.Az)
	if datacenterID == 0 || dcError != nil {
		return nil, messages.GetUserError("E0001", dcError, volumeRequest.Az)
	}
	if volumeRequest.Capacity == nil || *volumeRequest.Capacity == 0 {
		return nil, messages.GetUserError("E0036", nil)
	}

	if volumeRequest.ProviderType == "performance" && volumeRequest.Iops == nil {
		return nil, messages.GetUserError("E0037", nil)
	} else if volumeRequest.ProviderType == "endurance" && volumeRequest.Tier == nil {
		return nil, messages.GetUserError("E0038", nil)
	}

	// Step 2- Determine the category code to use for the order (and product package)
	order_type_is_saas, order_category_code := utils.GetOrderTypeAndCategory(service_offering, string(volumeRequest.ProviderType), volume_type)

	//Step 3- Get the product package for the given category code
	packageDetails, errPackage := utils.GetPackageDetails(sls.Logger, sls.Backend, "storage_as_a_service")
	if errPackage != nil {
		return nil, messages.GetUserError("E0017", errPackage, "storage_as_a_service")
	}

	sls.Logger.Info("Preparing order for order request .... ", zap.Reflect("Volume", volumeRequest))
	//Step 4- Based on the storage type and product package, build up the complex type
	// and array of price codes to include in the order object
	prices := []datatypes.Product_Item_Price{}
	var complex_type string
	base_type_name := "SoftLayer_Container_Product_Order_Network_"
	if order_type_is_saas {
		complex_type = base_type_name + "Storage_AsAService"
		if volumeRequest.ProviderType == "performance" {
			iops := utils.ToInt(*volumeRequest.Iops)
			sls.Logger.Info("new............. In PERFORMANCE")
			prices = []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, order_category_code))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, "storage_"+volume_type))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSPerformanceSpacePrice(sls.Logger, packageDetails, *volumeRequest.Capacity))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSPerformanceIopsPrice(sls.Logger, packageDetails, *volumeRequest.Capacity, iops))},
			}

			if volumeRequest.SnapshotSpace != nil && *volumeRequest.SnapshotSpace > 0 {
				prices = append(prices, datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSSnapshotSpacePrice(sls.Logger, packageDetails, *volumeRequest.SnapshotSpace, "", iops))})
			}
		} else { // volumeRequest.ProviderType == 'endurance'
			sls.Logger.Info("new............. In ENDURANCE")
			prices = []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, order_category_code))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, "storage_"+volume_type))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSEnduranceSpacePrice(sls.Logger, packageDetails, *volumeRequest.Capacity, *volumeRequest.Tier))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSEnduranceTierPrice(sls.Logger, packageDetails, *volumeRequest.Tier))},
			}
			if volumeRequest.SnapshotSpace != nil && *volumeRequest.SnapshotSpace > 0 {
				prices = append(prices, datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSSnapshotSpacePrice(sls.Logger, packageDetails, *volumeRequest.SnapshotSpace, *volumeRequest.Tier, 0))})
			}
		}
	} else { // offering package is enterprise or performance TODO: TEST ELSE PART
		if volumeRequest.ProviderType == "performance" {
			sls.Logger.Info("new............. In PLAIN PERFORMANCE")
			iops := utils.ToInt(*volumeRequest.Iops)
			if volume_type == "block" {
				complex_type = base_type_name + "PerformanceStorage_Iscsi"
			} else {
				complex_type = base_type_name + "PerformanceStorage_Nfs"
			}
			prices = []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, order_category_code))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPerformanceSpacePrice(sls.Logger, packageDetails, *volumeRequest.Capacity))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPerformanceIopsPrice(sls.Logger, packageDetails, *volumeRequest.Capacity, iops))},
			}
		} else { //volumeRequest.ProviderType == 'endurance'
			sls.Logger.Info("new............. In PLAIN ENDURANCE")
			complex_type = base_type_name + "Storage_Enterprise"
			prices = []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, order_category_code))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, "storage_"+volume_type))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetEnterpriseSpacePrice(sls.Logger, packageDetails, "endurance", *volumeRequest.Capacity, *volumeRequest.Tier))},
				datatypes.Product_Item_Price{Id: sl.Int(utils.GetEnterpriseEnduranceTierPrice(sls.Logger, packageDetails, *volumeRequest.Tier))},
			}

			if volumeRequest.SnapshotSpace != nil && *volumeRequest.SnapshotSpace > 0 {
				prices = append(prices, datatypes.Product_Item_Price{Id: sl.Int(utils.GetEnterpriseSpacePrice(sls.Logger, packageDetails, "snapshot", *volumeRequest.SnapshotSpace, *volumeRequest.Tier))})
			}
		}
	}

	orderContainer := &datatypes.Container_Product_Order_Network_Storage_AsAService{
		OsFormatType: &datatypes.Network_Storage_Iscsi_OS_Type{
			KeyName: sl.String("LINUX"),
		},
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType:      sl.String(complex_type),
			Location:         sl.String(strconv.Itoa(datacenterID)),
			PackageId:        sl.Int(*packageDetails.Id),
			Quantity:         sl.Int(1),
			UseHourlyPricing: &hourly_billing_flag,
			Prices:           prices,
		},
	}
	if order_type_is_saas == true {
		orderContainer.VolumeSize = sl.Int(*volumeRequest.Capacity)
		if volumeRequest.ProviderType == "performance" {
			iops := utils.ToInt(*volumeRequest.Iops)
			orderContainer.Iops = &iops
		}
	}
	// Step 5 - Place the order
	sls.Logger.Info("Order details for SL ....", zap.Reflect("orderContainer", orderContainer))
	productOrderObj := sls.Backend.GetProductOrderService()
	volOrder, err := productOrderObj.PlaceOrder(orderContainer, sl.Bool(false))
	if err != nil {
		return nil, messages.GetUserError("E0011", err, "Volume order failed")
	}
	sls.Logger.Info("Order status details ....", zap.Reflect("volOrder", volOrder), zap.Error(err)) // TODO: Will remove this

	//! Step 6- wait for provisioning completion
	volumeObj, errProv := sls.HandleProvisioning(*volOrder.OrderId)
	if errProv != nil {
		sls.Logger.Error("**Provisioning failed ....", zap.Reflect("volOrder", volOrder), zap.Error(errProv))
		return nil, messages.GetUserError("E0040", errProv, *volOrder.OrderId)
	}
	sls.Logger.Info("Volume details ... ", zap.Reflect("Volume", volumeObj))

	volumeObj.VolumeNotes = volumeRequest.VolumeNotes
	updateErr := sls.UpdateStorage(volumeObj)
	if updateErr != nil {
		sls.Logger.Error("Failed to update the storage", zap.Error(updateErr))
	}
	return volumeObj, nil
}

// Create the volume from snapshot with snapshot tags
func (sls *SLFileSession) VolumeCreateFromSnapshot(snapshot provider.Snapshot, tags map[string]string) (*provider.Volume, error) {
	//Setep 1: validate inputes
	volid := utils.ToInt(snapshot.VolumeID)
	snapshotID := utils.ToInt(snapshot.SnapshotID)
	if volid == 0 || snapshotID == 0 {
		sls.Logger.Error("No proper input, please provide volume ID and snapshot space size")
		return nil, messages.GetUserError("E0022", nil)
	}

	//! TODO: we need to get these values from user
	duplicateSnapshotSize := 0 // New snapshot space size for new volume
	duplicateVolumeSize := 0   // New volume size if needed
	duplicateVolumeTier := ""  // New Volume Tier
	duplicateIops := 0

	// Step 2- Get the original volume Details
	block_mask := "id,billingItem[location,hourlyFlag],snapshotCapacityGb,storageType[keyName],capacityGb,originalVolumeSize,provisionedIops,storageTierLevel,osType[keyName],staasVersion,hasEncryptionAtRest"
	storageObj := sls.Backend.GetNetworkStorageService()
	originalVolume, err := storageObj.ID(volid).Mask(block_mask).GetObject()
	if err != nil {
		sls.Logger.Error("While getting Original Volume", zap.Reflect("Error", err))
		return nil, messages.GetUserError("E0011", err, volid, "Not a valid volume ID")
	}

	// Step 3: verify original volume exists
	if originalVolume.BillingItem == nil {
		sls.Logger.Error("Original Volume has been deleted", zap.Reflect("BillingItem", originalVolume.BillingItem))
		return nil, messages.GetUserError("E0014", nil, volid)
	}

	// Step 4: Verify that the original volume has snapshot space (needed for duplication)
	if originalVolume.SnapshotCapacityGb == nil {
		sls.Logger.Error("Original Volume does not have Snapshot Capacity", zap.Reflect("SnapshotCapacity", originalVolume.SnapshotCapacityGb))
		return nil, messages.GetUserError("E0023", nil, volid)
	}
	originalSnapshotSize := utils.ToInt(*originalVolume.SnapshotCapacityGb)
	if originalSnapshotSize == 0 {
		return nil, messages.GetUserError("E0023", nil, volid)
	}

	// Step 4: Get the datacenter location ID for the original volume
	if originalVolume.BillingItem.Location == nil || originalVolume.BillingItem.Location.Id == nil {
		sls.Logger.Error("Original Volume does not have location ID", zap.Reflect("Location", originalVolume.BillingItem.Location))
		return nil, messages.GetUserError("E0024", nil, volid)
	}
	originalVolumeLocationID := *originalVolume.BillingItem.Location.Id

	// Step 5: Ensure the origin volume is SIaaS v2 or higher and supports Encryption at Rest
	if !utils.IsVolumeCreatedWithStaaS(originalVolume) {
		sls.Logger.Error("Original Volume does not support StaaS or not support Encryption at Rest")
		return nil, messages.GetUserError("E0018", nil, volid)
	}

	// Step 6: Check duplicate snapshot space provided or not, if not use original one
	if duplicateSnapshotSize == 0 {
		duplicateSnapshotSize = originalSnapshotSize
	}

	// Step 7: check duplicate size provided or not, if not use from original one
	if duplicateVolumeSize == 0 && originalVolume.CapacityGb != nil {
		if originalVolume.CapacityGb == nil || *originalVolume.CapacityGb == 0 {
			return nil, messages.GetUserError("E0027", nil, volid)
		}
		duplicateVolumeSize = *originalVolume.CapacityGb
	}

	// Step 8: Get the appropriate package for the order, 'storage_as_a_service' is used
	packageDetails, pkgError := utils.GetPackageDetails(sls.Logger, sls.Backend, "storage_as_a_service")
	if pkgError != nil {
		return nil, messages.GetUserError("E0017", pkgError, "storage_as_a_service")
	}

	// Step 9: get original storage type
	originalStorageType := *originalVolume.StorageType.KeyName
	isPerformanceVolume := false
	finalPrices := []datatypes.Product_Item_Price{}
	volumeCategory := fmt.Sprintf(`storage_%s`, "file")
	if strings.Contains(originalStorageType, "ENDURANCE") {
		if duplicateVolumeTier == "" {
			duplicateVolumeTier = utils.GetEnduranceTierIopsPerGB(sls.Logger, originalVolume)
		}
		finalPrices = []datatypes.Product_Item_Price{
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, "storage_as_a_service"))},
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, volumeCategory))},
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSEnduranceSpacePrice(sls.Logger, packageDetails, duplicateVolumeSize, duplicateVolumeTier))},
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSEnduranceTierPrice(sls.Logger, packageDetails, duplicateVolumeTier))},
		}

		if duplicateSnapshotSize > 0 {
			finalPrices = append(finalPrices, datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSSnapshotSpacePrice(sls.Logger, packageDetails, duplicateSnapshotSize, duplicateVolumeTier, 0))})
		}
	} else if strings.Contains(originalStorageType, "PERFORMANCE") {
		isPerformanceVolume = true
		if duplicateIops == 0 && originalVolume.ProvisionedIops != nil {
			duplicateIops = utils.ToInt(*originalVolume.ProvisionedIops)
			if duplicateIops == 0 {
				return nil, messages.GetUserError("E0025", nil, volid)
			}
		}

		finalPrices = []datatypes.Product_Item_Price{
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, "storage_as_a_service"))},
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetPriceIDByCategory(sls.Logger, packageDetails, volumeCategory))},
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSPerformanceSpacePrice(sls.Logger, packageDetails, duplicateVolumeSize))},
			datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSPerformanceIopsPrice(sls.Logger, packageDetails, duplicateVolumeSize, duplicateIops))},
		}
		if duplicateSnapshotSize > 0 {
			finalPrices = append(finalPrices, datatypes.Product_Item_Price{Id: sl.Int(utils.GetSaaSSnapshotSpacePrice(sls.Logger, packageDetails, duplicateVolumeSize, "", duplicateIops))})
		}

	} else {
		sls.Logger.Error("Origin volume does not have a valid storage type (with an appropriate keyName to indicate the volume is a PERFORMANCE or an ENDURANCE volume)")
		return nil, messages.GetUserError("E0026", nil, volid)
	}

	// Step 10: create a duplicate volume order
	duplicate_order := &datatypes.Container_Product_Order_Network_Storage_AsAService{
		OsFormatType: &datatypes.Network_Storage_Iscsi_OS_Type{
			KeyName: sl.String("LINUX"),
		},
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType:      sl.String("SoftLayer_Container_Product_Order_Network_Storage_AsAService"),
			Location:         sl.String(strconv.Itoa(originalVolumeLocationID)),
			PackageId:        packageDetails.Id,
			Quantity:         sl.Int(1),
			UseHourlyPricing: sl.Bool(false),
			Prices:           finalPrices,
		},
		DuplicateOriginVolumeId:   sl.Int(*originalVolume.Id),
		DuplicateOriginSnapshotId: sl.Int(snapshotID),
		VolumeSize:                sl.Int(duplicateVolumeSize),
	}
	if isPerformanceVolume {
		duplicate_order.Iops = sl.Int(duplicateIops)
	}

	// Step 11- Placing order
	sls.Logger.Info("Duplicate Order details ... ", zap.Reflect("Order Details->", duplicate_order))
	productOrderObj := sls.Backend.GetProductOrderService()
	newOrderID, orderError := productOrderObj.PlaceOrder(duplicate_order, sl.Bool(false))
	if orderError != nil {
		return nil, messages.GetUserError("E0028", orderError, volid, snapshotID)
	}
	sls.Logger.Info("Successfully placed order from snapshot .... ", zap.Reflect("orderID", *newOrderID.OrderId), zap.Reflect("OriginalVolumeID", volid), zap.Reflect("OriginalSnapshotID", snapshotID))

	// Step 12- handle provisioning
	volumeObj, errProv := sls.HandleProvisioning(*newOrderID.OrderId)
	if errProv != nil {
		sls.Logger.Error("**Provisioning failed ....", zap.Int("volOrder", *newOrderID.OrderId), zap.Error(errProv))
		return nil, messages.GetUserError("E0040", errProv, *newOrderID.OrderId)
	}
	sls.Logger.Info("Volume details ... ", zap.Reflect("Volume", volumeObj))
	return volumeObj, nil
}

func (sls *SLFileSession) HandleProvisioning(orderID int) (*provider.Volume, error) {
	storageID, err := sls.SLSession.HandleProvisioning(orderID)
	if err == nil {
		volume, err := sls.VolumeGet(strconv.Itoa(*storageID))
		return volume, err
	}
	return nil, err
}

// VolumeGet Get the volume by using ID
func (sls *SLFileSession) VolumeGet(id string) (*provider.Volume, error) {
	// Step 1: validate input
	volumeID := utils.ToInt(id)
	if volumeID == 0 {
		return nil, messages.GetUserError("E0011", nil, volumeID, "0 is not the correct volume ID")
	}

	// Step 2: Get volume details from SL
	mask := "id,username,capacityGb,serviceResourceBackendIpAddress,fileNetworkMountAddress,createDate,snapshotCapacityGb,parentVolume[snapshotSizeBytes],storageType[keyName],serviceResource[datacenter[name]],provisionedIops,originalVolumeName,storageTierLevel,notes"
	storageObj := sls.Backend.GetNetworkStorageService()
	storage, err := storageObj.ID(volumeID).Mask(mask).GetObject()
	if err != nil {
		return nil, messages.GetUserError("E0011", err, volumeID, "Not a valid volume ID")
	}

	// Step 3: Convert volume to framework based volume object
	vol := utils.ConvertToVolumeType(storage, sls.Logger, SoftLayer, VolumeTypeFile)
	return vol, err
}

// Get volume lists by using snapshot tags
func (sls *SLFileSession) VolumesList(tags map[string]string) ([]*provider.Volume, error) {
	//! TODO: we may implement
	return nil, nil
}
