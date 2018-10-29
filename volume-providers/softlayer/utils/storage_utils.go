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
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/messages"
	"github.com/softlayer/softlayer-go/datatypes"
	"go.uber.org/zap"
)

var ENDURANCE_TIERS = map[string]int{
	"0.25": 100,
	"2":    200,
	"4":    300,
	"10":   1000,
}

var IOPS = map[string]string{"READHEAVY_TIER": "2", "WRITEHEAVY_TIER": "4", "10_IOPS_PER_GB": "10"}

// GetDataCenterID
func GetDataCenterID(logger *zap.Logger, sess backend.Session, dataCenterName string) (int, error) {
	locationObj := sess.GetLocationService()
	locations, err := locationObj.GetDatacenters()
	if err != nil {
		logger.Error("Could not find location", zap.String("locaion Name", dataCenterName), zap.Error(err))
		return 0, err
	}
	for _, location := range locations {
		if location.Name != nil && *location.Name == dataCenterName && location.Id != nil {
			logger.Info("Got location ID: ", zap.Int("ID", *location.Id))
			return *location.Id, nil
		}
	}
	logger.Error("Could not find location", zap.String("locaion Name", dataCenterName))
	return 0, err
}

// GetPackageDetails
func GetPackageDetails(logger *zap.Logger, sess backend.Session, category string) (datatypes.Product_Package, error) {
	packageFilter := fmt.Sprintf(`{ "categories":{"categoryCode":{"operation":"%s"}}, "statusCode": {"operation": "ACTIVE"}}`, category)
	packageMask := "id,name,items[prices[categories],attributes]"
	packages, packageErr := sess.GetProductPackageService().Filter(packageFilter).Mask(packageMask).GetAllObjects()
	if packageErr != nil {
		return datatypes.Product_Package{}, packageErr
	}

	if len(packages) != 1 {
		packageErr = fmt.Errorf("expected one product package with keyname %q but got %d", category, len(packages))
		return datatypes.Product_Package{}, packageErr
	}
	return packages[0], nil
}

// isCategoryPresent
func isCategoryPresent(logger *zap.Logger, categories []datatypes.Product_Item_Category, categoryName string) bool {
	for _, category := range categories {
		if category.CategoryCode != nil && *category.CategoryCode == categoryName { //TODO: need to check Name is fine or CategoryCode, so far category code works fine
			return true
		}
	}
	return false
}

// GetPriceIDFromItemByPriceCategory
func GetPriceIDFromItemByPriceCategory(logger *zap.Logger, item datatypes.Product_Item, category string, restrictionType string, restrictionValue int) int {
	for _, price := range item.Prices {
		if price.LocationGroupId != nil { //! skip the location specific price
			continue
		}

		//! if restriction type and value given
		if restrictionType != "" && restrictionValue != 0 {
			if price.CapacityRestrictionType == nil || price.CapacityRestrictionMinimum == nil || price.CapacityRestrictionMaximum == nil {
				continue
			}

			resMin := ToInt(*price.CapacityRestrictionMinimum)
			resMax := ToInt(*price.CapacityRestrictionMaximum)
			if *price.CapacityRestrictionType != restrictionType || restrictionValue < resMin || restrictionValue > resMax {
				continue
			}
		}

		// check category
		if !isCategoryPresent(logger, price.Categories, category) {
			continue
		}
		return *price.Id
	}

	return 0
}

// GetPriceIDByCategory
func GetPriceIDByCategory(logger *zap.Logger, packageDetails datatypes.Product_Package, category string) int {
	for _, item := range packageDetails.Items {
		priceId := GetPriceIDFromItemByPriceCategory(logger, item, category, "", 0)
		if priceId != 0 {
			return priceId
		}
	}
	return 0
}

// GetSaaSEnduranceSpacePrice
func GetSaaSEnduranceSpacePrice(logger *zap.Logger, packageDetails datatypes.Product_Package, size int, tier string) int {
	keyName := fmt.Sprintf(`STORAGE_SPACE_FOR_%s_IOPS_PER_GB`, tier)
	keyName = strings.Replace(keyName, ".", "_", -1)
	for _, item := range packageDetails.Items {
		if item.KeyName == nil || *item.KeyName != keyName {
			continue
		}

		capacityMin := ToInt(*item.CapacityMinimum)
		capacityMax := ToInt(*item.CapacityMaximum)
		if size < capacityMin || size > capacityMax {
			continue
		}

		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "performance_storage_space", "", 0)
		if priceId != 0 {
			return priceId
		}
	}
	return 0
}

// GetSaaSPerformanceSpacePrice
func GetSaaSPerformanceSpacePrice(logger *zap.Logger, packageDetails datatypes.Product_Package, size int) int {
	for _, item := range packageDetails.Items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != "performance_storage_space" {
			continue
		}

		if item.CapacityMinimum == nil || item.CapacityMaximum == nil {
			continue
		}

		capacityMin := ToInt(*item.CapacityMinimum)
		capacityMax := ToInt(*item.CapacityMaximum)
		if size < capacityMin || size > capacityMax {
			continue
		}

		keyName := fmt.Sprintf(`%d_%d_GBS`, capacityMin, capacityMax)
		logger.Info("GetSaaSPerformanceSpacePrice......", zap.String("keyName", keyName))
		if item.KeyName == nil || *item.KeyName != keyName {
			continue
		}
		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "performance_storage_space", "", 0)
		if priceId != 0 {
			return priceId
		}
	}
	return 0
}

// GetSaaSPerformanceIopsPrice
func GetSaaSPerformanceIopsPrice(logger *zap.Logger, packageDetails datatypes.Product_Package, size int, iops int) int {
	for _, item := range packageDetails.Items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != "performance_storage_iops" {
			continue
		}

		if item.CapacityMinimum == nil || item.CapacityMaximum == nil {
			continue
		}

		capacityMin := ToInt(*item.CapacityMinimum)
		capacityMax := ToInt(*item.CapacityMaximum)
		if iops < capacityMin || iops > capacityMax {
			continue
		}

		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "performance_storage_iops", "STORAGE_SPACE", size)
		if priceId != 0 {
			return priceId
		}
	}
	return 0
}

//GetSaaSEnduranceTierPrice
func GetSaaSEnduranceTierPrice(logger *zap.Logger, packageDetails datatypes.Product_Package, tier string) int {
	target_capacity := ENDURANCE_TIERS[tier]
	for _, item := range packageDetails.Items {
		if item.ItemCategory == nil || (item.ItemCategory != nil && *item.ItemCategory.CategoryCode != "storage_tier_level") {
			continue
		}

		if item.Capacity == nil || int(*item.Capacity) != target_capacity {
			continue
		}

		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "storage_tier_level", "", 0)
		if priceId != 0 {
			return priceId
		}
	}
	return 0
}

// GetSaaSSnapshotSpacePrice
func GetSaaSSnapshotSpacePrice(logger *zap.Logger, packageDetails datatypes.Product_Package, size int, tier string, iops int) int {
	targetRestrictionType := ""
	targetRestrictionValue := 0
	if tier != "" {
		targetRestrictionType = "STORAGE_TIER_LEVEL"
		targetRestrictionValue = ENDURANCE_TIERS[tier]
	} else {
		targetRestrictionType = "IOPS"
		targetRestrictionValue = iops
	}

	for _, item := range packageDetails.Items {
		if item.Capacity == nil || *item.Capacity != datatypes.Float64(size) {
			continue
		}

		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "storage_snapshot_space", targetRestrictionType, targetRestrictionValue)
		if priceId != 0 {
			return priceId
		}
	}
	return 0
}

// GetSaaSSnapshotOrderSpacePrice
func GetSaaSSnapshotOrderSpacePrice(logger *zap.Logger, packageDetails datatypes.Product_Package, size int, restrictionType string, restrictionValue int) int {
	for _, item := range packageDetails.Items {
		if item.Capacity == nil || *item.Capacity != datatypes.Float64(size) {
			continue
		}

		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "storage_snapshot_space", restrictionType, restrictionValue)
		if priceId != 0 {
			return priceId
		}
	}
	//fmt.Errorf("Could not find price")
	return 0
}

//GetEnduranceTierIopsPerGB(originalVolume)
func GetEnduranceTierIopsPerGB(logger *zap.Logger, originalVolume datatypes.Network_Storage) string {

	if originalVolume.StorageTierLevel == nil {
		return ""
	}
	tier := *originalVolume.StorageTierLevel
	iopsPerGB := ""
	if tier == "LOW_INTENSITY_TIER" {
		iopsPerGB = "0.25"
	} else if tier == "READHEAVY_TIER" {
		iopsPerGB = "2"
	} else if tier == "WRITEHEAVY_TIER" {
		iopsPerGB = "4"
	} else if tier == "10_IOPS_PER_GB" {
		iopsPerGB = "10"
	} else {
		logger.Info("Could not found iops for Tier ", zap.String("Tier", tier))
	}
	return iopsPerGB
}

func GetOrderTypeAndCategory(service_offering, storage_type, volume_type string) (bool, string) {
	var order_type_is_saas bool
	var order_category_code string
	if service_offering == "storage_as_a_service" {
		order_type_is_saas = true
		order_category_code = "storage_as_a_service"
	} else if service_offering == "enterprise" {
		order_type_is_saas = false
		if storage_type == "endurance" {
			order_category_code = "storage_service_enterprise"
		} else {
			return false, ""
			//raise exceptions.SoftLayerError(
			//              "The requested offering package, '%s', is not available for "
			//"the '%s' storage type." % (service_offering, storage_type))
		}
	} else if service_offering == "performance" {
		order_type_is_saas = false
		if storage_type == "performance" {
			if volume_type == "block" {
				order_category_code = "performance_storage_iscsi"
			} else {
				order_category_code = "performance_storage_nfs"
			}
		} else {
			return false, ""
			//raise exceptions.SoftLayerError(
			//    "The requested offering package, '%s', is not available for "
			//    "the '%s' storage type." % (service_offering, storage_type))
		}
	} else {
		return false, ""
		// raise exceptions.SoftLayerError(
		//   "The requested service offering package is not valid. "
		//   "Please check the available options and try again.")
	}
	return order_type_is_saas, order_category_code
}

func GetPerformanceSpacePrice(logger *zap.Logger, packageDetails datatypes.Product_Package, size int) int {
	for _, item := range packageDetails.Items {
		if int(*item.Capacity) != size {
			continue
		}
		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "performance_storage_space", "", 0) //_find_price_id(item["prices"], "performance_storage_space")
		if priceId != 0 {
			return priceId
		}
	}
	return 0
}

func GetPerformanceIopsPrice(logger *zap.Logger, packageDetails datatypes.Product_Package, size int, iops int) int {
	for _, item := range packageDetails.Items {
		if int(*item.Capacity) != iops {
			continue
		}

		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "performance_storage_iops", "STORAGE_SPACE", size) //_find_price_id(item['prices'], 'performance_storage_iops', 'STORAGE_SPACE', size)
		if priceId != 0 {
			return priceId
		}
	}
	return 0
	//raise ValueError("Could not find price for iops for the given volume")
}

func GetEnterpriseSpacePrice(logger *zap.Logger, packageDetails datatypes.Product_Package, category string, size int, tier_level string) int {
	var category_code string
	if category == "snapshot" {
		category_code = "storage_snapshot_space"
	} else if category == "replication" {
		category_code = "performance_storage_replication"
	} else { // category == "endurance"
		category_code = "performance_storage_space"
	}
	level := ENDURANCE_TIERS[tier_level]

	for _, item := range packageDetails.Items {
		if int(*item.Capacity) != size {
			continue
		}
		priceId := GetPriceIDFromItemByPriceCategory(logger, item, category_code, "STORAGE_TIER_LEVEL", level)
		if priceId != 0 {
			return priceId
		}
	}
	return 0 //raise ValueError("Could not find price for %s storage space" % category)
}

func GetEnterpriseEnduranceTierPrice(logger *zap.Logger, packageDetails datatypes.Product_Package, tier_level string) int {
	for _, item := range packageDetails.Items {
		/*     for attribute in item.get('attributes', []){
		            if int(attribute['value']) == ENDURANCE_TIERS.get(tier_level){
		                break
		}
		}*/
		priceId := GetPriceIDFromItemByPriceCategory(logger, item, "storage_tier_level", "", 0)
		if priceId != 0 {
			return priceId
		}
	}
	return 0
}

//IsVolumeCreatedWithStaaS
func IsVolumeCreatedWithStaaS(volume datatypes.Network_Storage) bool {
	version := ToInt(*volume.StaasVersion)
	return version > 1 && volume.HasEncryptionAtRest != nil && *volume.HasEncryptionAtRest
}

//ToInt
func ToInt(valueInInt string) int {
	value, err := strconv.Atoi(valueInInt)
	if err != nil {
		return 0
	}
	return value
}

type retryFuncProv func() (bool, error)

func ProvisioningRetry(fn retryFuncProv, logger *zap.Logger, timeoutSec string, retryIntervalSec string) error {
	provisionTimeout, err := time.ParseDuration(timeoutSec)
	if err != nil {
		return err
	}
	logger.Info("provisionTimeout", zap.Reflect("provisionTimeout", provisionTimeout))
	pollingInterval, err := time.ParseDuration(retryIntervalSec)
	logger.Info("pollingInterval", zap.Reflect("pollingInterval", pollingInterval))
	if err != nil {
		return err
	}
	iterations := int(provisionTimeout / pollingInterval)
	logger.Info("iterations", zap.Int("iterations", iterations))
	var done bool
	for i := 0; i <= iterations; i++ {
		logger.Info("attempt", zap.Int("attempt", i))
		if i != 0 {
			time.Sleep(pollingInterval)
		}

		if done, err = fn(); done {
			return nil
		}
	}
	return messages.GetUserError("E0039", err)
}

func ConvertToVolumeType(storage datatypes.Network_Storage, logger *zap.Logger, prName provider.VolumeProvider, volType provider.VolumeType) (volume *provider.Volume) {
	logger.Info("in CovertToVolumeType")
	volume = &provider.Volume{}
	volume.VolumeID = strconv.Itoa(*storage.Id)
	var newnotes map[string]string
	if storage.Notes != nil {
		_ = json.Unmarshal([]byte(*storage.Notes), &newnotes)
		volume.Region = newnotes["region"]
	}
	volume.Provider = prName
	volume.VolumeType = volType
	if storage.StorageType != nil && storage.StorageType.KeyName != nil {
		volume.ProviderType = provider.VolumeProviderType(*storage.StorageType.KeyName)
	}
	volume.Capacity = storage.CapacityGb
	if storage.ServiceResource != nil && storage.ServiceResource.Datacenter != nil && storage.ServiceResource.Datacenter.Name != nil {
		volume.Az = *storage.ServiceResource.Datacenter.Name
	}
	if storage.SnapshotCapacityGb != nil {
		snpCapGb := ToInt(*storage.SnapshotCapacityGb)
		volume.SnapshotSpace = &snpCapGb
	}

	if storage.StorageTierLevel != nil {
		volume.Tier = storage.StorageTierLevel
		iops := IOPS[*storage.StorageTierLevel]
		volume.Iops = &iops //storage.ProvisionedIops
	}
	if storage.CreateDate != nil {
		volume.CreationTime, _ = time.Parse(time.RFC3339, storage.CreateDate.String())
	}
	volume.VolumeNotes = newnotes
	return
}

func ConvertToNetworkStorage(storage datatypes.Network_Storage_Iscsi) datatypes.Network_Storage {
	networkStorageIscsi := datatypes.Network_Storage{}
	networkStorageIscsi.Id = storage.Id
	networkStorageIscsi.Notes = storage.Notes
	networkStorageIscsi.StorageType = storage.StorageType
	networkStorageIscsi.CapacityGb = storage.CapacityGb
	networkStorageIscsi.SnapshotCapacityGb = storage.SnapshotCapacityGb
	networkStorageIscsi.StorageTierLevel = storage.StorageTierLevel
	networkStorageIscsi.CreateDate = storage.CreateDate
	networkStorageIscsi.ServiceResourceName = storage.ServiceResourceName
	return networkStorageIscsi
}

func ConverStringToMap(mapString string) map[string]string {
	mapString = strings.Replace(mapString, "{", "", -1)
	mapString = strings.Replace(mapString, "}", "", -1)
	mapString = strings.Replace(mapString, "'", "", -1)
	splitString := strings.Split(mapString, ",")
	m := make(map[string]string)
	for _, pair := range splitString {
		z := strings.Split(pair, ":")
		if len(z) < 1 {
			continue
		}
		m[z[0]] = z[1]
	}
	return m
}

func ConvertToLocalSnapshotObject(storageSnapshot datatypes.Network_Storage, pName provider.VolumeProvider, volType provider.VolumeType) *provider.Snapshot {
	snapshot := &provider.Snapshot{}
	volume := provider.Volume{}
	volume.Provider = pName
	volume.VolumeType = volType

	if storageSnapshot.ParentVolume != nil && storageSnapshot.ParentVolume.Id != nil {
		volume.VolumeID = strconv.Itoa(*storageSnapshot.ParentVolume.Id)
	}

	if storageSnapshot.ParentVolume != nil && storageSnapshot.ParentVolume.SnapshotCapacityGb != nil {
		snapspace := ToInt(*storageSnapshot.ParentVolume.SnapshotCapacityGb)
		volume.SnapshotSpace = &snapspace
	}

	if storageSnapshot.ParentVolume != nil && storageSnapshot.ParentVolume.SnapshotSizeBytes != nil {
		snapSize := ToInt(*storageSnapshot.ParentVolume.SnapshotSizeBytes)
		snapshot.SnapshotSize = &snapSize
	}
	snapshot.Volume = volume

	if storageSnapshot.CreateDate != nil {
		snapshot.SnapshotCreationTime, _ = time.Parse(time.RFC3339, storageSnapshot.CreateDate.String())
	}

	if storageSnapshot.Id != nil {
		snapshot.SnapshotID = strconv.Itoa(*storageSnapshot.Id)
	}

	if storageSnapshot.Notes != nil {
		snapshot.SnapshotTags = ConverStringToMap(*storageSnapshot.Notes)
	}

	return snapshot
}
