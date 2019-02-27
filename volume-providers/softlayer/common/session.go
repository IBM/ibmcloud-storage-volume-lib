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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/utils"
	"github.com/softlayer/softlayer-go/filter"
	"go.uber.org/zap"
)

type SLSessionCommonInterface interface {
	// Get the volume by using ID  //
	VolumeGet(id string) (*provider.Volume, error)
}

// SLSession implements lib.Session
type SLSession struct {
	SLAccountID        int
	Url                string
	Backend            backend.Session
	Logger             *zap.Logger
	Config             *config.SoftlayerConfig
	ContextCredentials provider.ContextCredentials
	VolumeType         provider.VolumeType
	Provider           provider.VolumeProvider
	SLSessionCommonInterface
}

const SnapshotMask = "id,username,capacityGb,createDate,snapshotCapacityGb,parentVolume[snapshotSizeBytes],parentVolume[snapshotCapacityGb],parentVolume[id],parentVolume[storageTierLevel],parentVolume[notes],storageType[keyName],serviceResource[datacenter[name]],billingItem[location,hourlyFlag],provisionedIops,lunId,originalVolumeName,storageTierLevel,notes"

var (
	VolumeDeleteReason = "deleted by ibm-volume-lib on behalf of user request"
)

//GetVolumeByRequestID getvolume by  request ID. The request ID the identifier returned after placing the order
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
		return nil, fmt.Errorf("unable to find network storage associated with order %s", requestID)
	case 1:
		// double check if correct storage is found by matching requestID and fouund orderID
		orderID := strconv.Itoa(*storages[0].BillingItem.OrderItem.Order.Id)
		if orderID == requestID {
			vol := utils.ConvertToVolumeType(storages[0], sls.Logger, sls.Provider, sls.VolumeType)
			return vol, nil
		} else {
			sls.Logger.Error("Incorrect storage found", zap.String("requestID", requestID), zap.Reflect("storage", storages[0]))
			return nil, fmt.Errorf("Incorrect storage found %s associated with order %s", orderID, requestID)
		}
	default:
		return nil, fmt.Errorf("multiple storage volumes associated with order %s", requestID)
	}

}

func (sls *SLSession) UpdateStorage(volume *provider.Volume) error {
	sls.Logger.Info("Entry UpdateStorage", zap.Reflect("Volume", volume))
	//utils.ConvertToNetworkStorage()
	volumeid, _ := strconv.Atoi(volume.VolumeID)
	networkService := sls.Backend.GetNetworkStorageService().ID(volumeid)
	storage, err := networkService.GetObject()
	notesByte, _ := json.Marshal(volume.VolumeNotes)
	notesStr := string(notesByte)
	storage.Notes = &notesStr
	networkService.EditObject(&storage)
	return err

}

// Delete the volume
func (sls *SLSession) DeleteVolume(vol *provider.Volume) error {
	//! Step 1- verify inputes
	if vol == nil {
		return messages.GetUserError("E0035", nil)
	}
	volumeID := utils.ToInt(vol.VolumeID)
	if volumeID == 0 {
		return messages.GetUserError("E0035", nil)
	}

	return sls.deleteStorage(volumeID)
}

//HandleProvisioning
func (sls *SLSession) HandleProvisioning(orderID int) (*int, error) {
	sls.Logger.Info("Handling provisioning for  ....", zap.Reflect("OrderID", orderID))
	storageID, errID := utils.GetStorageID(sls.Backend, string(sls.VolumeType), orderID, sls.Logger, sls.Config)
	if errID != nil {
		return nil, messages.GetUserError("E0011", errID, orderID)
	}
	sls.Logger.Info("Successfully got the volume ID ....", zap.Reflect("VolumeID", storageID))

	// Step 2- Cancel order if storageID not yet approved
	err := utils.WaitForTransactionsComplete(sls.Backend, storageID, sls.Logger, sls.Config)
	if err != nil {
		sls.Logger.Error("**Error while provisioning order", zap.Int("storageId", storageID), zap.Error(err))
		// cancel order
		sls.Logger.Error("**Cancelling/Cleaning up the Storage Order due to incomplete provision", zap.Int("storageId", storageID))
		errDelete := sls.deleteStorage(storageID)
		if errDelete != nil {
			sls.Logger.Error("**Failed to delete un-provisioned storage, user need to delete it manually", zap.Int("storageId", storageID), zap.Reflect("Error", errDelete))
		}
		return nil, err
	}
	sls.Logger.Info("storage details ... ", zap.Reflect("storageID", storageID))
	return &storageID, nil
}

//getBillingItemID
func (sls *SLSession) getBillingItemID(storageID int) (int, bool, error) {
	//! Step 1: Get volume details
	immediate := false
	mask := "id,billingItem[id,hourlyFlag]"
	storageObj := sls.Backend.GetNetworkStorageService()
	storage, err := storageObj.ID(storageID).Mask(mask).GetObject()
	if err != nil {
		return 0, immediate, messages.GetUserError("E0011", err, storageID, "Volume ID not found")
	}
	sls.Logger.Info("", zap.Reflect("StorageBillingInfo", storage))

	if storage.BillingItem == nil {
		return 0, immediate, messages.GetUserError("E0014", nil, storageID)
	}
	if storage.BillingItem.HourlyFlag != nil {
		immediate = true
	}
	return *storage.BillingItem.Id, immediate, nil
}

//deleteStorage
func (sls *SLSession) deleteStorage(networkStorageID int) error {
	//! Step 1: Get volume details
	billingId, immediate, err := sls.getBillingItemID(networkStorageID)
	if err != nil || billingId <= 0 {
		return err
	}

	// Step 2: Cancel the volume from SL
	billingObj := sls.Backend.GetBillingItemService()
	cancelAssociatedBillingItems := true
	_, err = billingObj.ID(billingId).CancelItem(&immediate, &cancelAssociatedBillingItems, &VolumeDeleteReason, &VolumeDeleteReason)
	if err != nil {
		return messages.GetUserError("E0006", err, networkStorageID)
	}

	sls.Logger.Info("Successfully deleted storage  ..", zap.Int("StorageID", networkStorageID), zap.Int("BillingID", billingId))
	return nil
}

// Create the snapshot from the volume
func (sls *SLSession) CreateSnapshot(volume *provider.Volume, tags map[string]string) (*provider.Snapshot, error) {
	// Step 1: Validate input
	if volume == nil {
		return nil, messages.GetUserError("E0011", nil, nil, "nil volume struct")
	}
	volumeID := utils.ToInt(volume.VolumeID)
	if volumeID == 0 {
		return nil, messages.GetUserError("E0011", nil, volumeID, "Not a valid volume ID")
	}

	// Step 2: Get the volume details
	block_mask := "id,billingItem[location,hourlyFlag],snapshotCapacityGb,storageType[keyName],capacityGb,originalVolumeSize,provisionedIops,storageTierLevel,osType[keyName],staasVersion,hasEncryptionAtRest"
	storageObj := sls.Backend.GetNetworkStorageService()
	originalVolume, err := storageObj.ID(volumeID).Mask(block_mask).GetObject()
	if err != nil {
		return nil, messages.GetUserError("E0011", err, volumeID, "Not a valid volume ID")
	}

	// Step 3: verify original volume exists
	if originalVolume.BillingItem == nil {
		return nil, messages.GetUserError("E0014", nil, volumeID)
	}

	// Step 3: Verify that the original volume has snapshot space (needed for duplication)
	if originalVolume.SnapshotCapacityGb == nil || utils.ToInt(*originalVolume.SnapshotCapacityGb) <= 0 {
		return nil, messages.GetUserError("E0023", nil, volumeID)
	}

	newtags, _ := json.Marshal(tags)
	snapshotTags := string(newtags)
	snapshotvol, err := storageObj.ID(volumeID).CreateSnapshot(&snapshotTags)
	if err != nil {
		return nil, messages.GetUserError("E0029", err, volumeID)
	}
	sls.Logger.Info("Successfully created snapshot for given volume ... ", zap.Reflect("VolumeID", volumeID), zap.Reflect("SnapshotID", snapshotvol)) //*snapshotvol.Id

	// Setep 4: Converting to local type
	snapshot := &provider.Snapshot{}
	snapshot.SnapshotID = strconv.Itoa(*snapshotvol.Id)
	snapshot.SnapshotSpace = snapshotvol.CapacityGb
	snapshot.Volume = *volume
	snapshot.CreationTime, _ = time.Parse(time.RFC3339, snapshotvol.CreateDate.String())
	snapshot.SnapshotTags = tags
	return snapshot, err
}

// Delete the snapshot
func (sls *SLSession) DeleteSnapshot(del *provider.Snapshot) error {
	// Step 1- Validate inputes
	if del == nil {
		return messages.GetUserError("E0030", nil)
	}
	snapshotId := utils.ToInt(del.SnapshotID)
	if snapshotId == 0 {
		return messages.GetUserError("E0030", nil)
	}

	//! Step 2- Delete the snapshot from SL
	storageObj := sls.Backend.GetNetworkStorageService()
	_, err := storageObj.ID(snapshotId).DeleteObject()
	if err != nil {
		return messages.GetUserError("E0031", err, snapshotId)
	}
	sls.Logger.Info("Successfully deleted snapshot ....", zap.Reflect("SnapshotID", snapshotId))
	return nil
}

// Get the snapshot
func (sls *SLSession) GetSnapshot(snapshotId string) (*provider.Snapshot, error) {
	// Step 1- Validate inputes
	snapshotID := utils.ToInt(snapshotId)
	if snapshotID == 0 {
		return nil, messages.GetUserError("E0030", nil)
	}

	// Step 2- Get the snapshot details from SL
	filter := fmt.Sprintf(`{"networkStorage":{"nasType":{"operation":"SNAPSHOT"},"id": {"operation":%d}}}`, snapshotID)
	accService := sls.Backend.GetAccountService()
	storageSnapshot, err := accService.Filter(filter).Mask(SnapshotMask).GetNetworkStorage()
	if err != nil {
		return nil, messages.GetUserError("E0032", err, snapshotID)
	}
	if len(storageSnapshot) <= 0 {
		return nil, messages.GetUserError("E0032", err, snapshotID)
	}
	sls.Logger.Info("########======> Successfully get the snapshot details", zap.Reflect("snapshot", storageSnapshot[0]))
	// Setep 3: Converting to local type
	snapshot := utils.ConvertToLocalSnapshotObject(storageSnapshot[0], sls.Logger, sls.Provider, sls.VolumeType)
	return snapshot, nil
}

// Snapshot list by using tags
func (sls *SLSession) ListSnapshots() ([]*provider.Snapshot, error) {
	// Step 1- Get all snapshots from the SL which belongs to a IBM Infrastructure a/c
	filter := fmt.Sprintf(`{"networkStorage":{"nasType":{"operation":"SNAPSHOT"}}}`)
	accService := sls.Backend.GetAccountService()
	storageSnapshots, err := accService.Filter(filter).Mask(SnapshotMask).GetNetworkStorage()
	if err != nil {
		return nil, messages.GetUserError("E0032", err)
	}
	sls.Logger.Info("Successfully got all snapshot from SL", zap.Reflect("snapshots", storageSnapshots))

	// convert to local type
	snList := []*provider.Snapshot{}
	for _, stSnapshot := range storageSnapshots {
		snapshot := utils.ConvertToLocalSnapshotObject(stSnapshot, sls.Logger, sls.Provider, sls.VolumeType)
		snList = append(snList, snapshot)
	}
	return snList, nil
}

// List all the snapshots for a given volume
func (sls *SLSession) ListAllSnapshots(volumeID string) ([]*provider.Snapshot, error) {
	// Step 1- Validate inputs
	orderID := utils.ToInt(volumeID)
	if orderID == 0 {
		return nil, messages.GetUserError("E0011", nil, "Not a valid volume ID")
	}

	// Step 2- Get volume details
	storageObj := sls.Backend.GetNetworkStorageService()
	mask := "id,billingItem[location,hourlyFlag],storageType[keyName],storageTierLevel,provisionedIops,staasVersion,hasEncryptionAtRest"
	_, err := storageObj.ID(orderID).Mask(mask).GetObject()
	if err != nil {
		return nil, messages.GetUserError("E0011", err, orderID, "Not a valid volume ID")
	}

	// Step 3- Get all snapshots from a volume
	snapshotvol, err := storageObj.ID(orderID).Mask(SnapshotMask).GetSnapshots()
	if err != nil {
		return nil, messages.GetUserError("E0034", err, orderID)
	}
	sls.Logger.Info("Successfully got all snapshots from given volume ID .....", zap.Reflect("VolumeID", orderID), zap.Reflect("Snapshots", snapshotvol))

	// convert to local type
	snList := []*provider.Snapshot{}
	for _, stSnapshot := range snapshotvol {
		snapshot := utils.ConvertToLocalSnapshotObject(stSnapshot, sls.Logger, sls.Provider, sls.VolumeType)
		snList = append(snList, snapshot)
	}
	return snList, nil
}

// UpdateAuthorization for the volume
func (sls *SLSession) UpdateAuthorization(authorizationRequest provider.AuthorizationRequest) error {
	sls.Logger.Info("Entry UpdateAuthorization", zap.Reflect("authorizationRequest", authorizationRequest))
	volumeID, _ := strconv.Atoi(authorizationRequest.Volume.VolumeID)
	var err error
	if authorizationRequest.Subnets != nil && len(authorizationRequest.Subnets) > 0 {
		err = sls.updateSubnetAuthorization(volumeID, authorizationRequest.Subnets)
	}
	if authorizationRequest.HostIps != nil && len(authorizationRequest.HostIps) > 0 {
		err = sls.updateHostIPAuthorization(volumeID, authorizationRequest.HostIps)
	}
	sls.Logger.Info("Exit UpdateAuthorization ", zap.Error(err))
	return err

}

func (sls *SLSession) updateSubnetAuthorization(volumeID int, subnetIDs []string) error {
	sls.Logger.Info("Entry updateSubnetAuthorization", zap.Reflect("subnetIDs", subnetIDs))
	subnetList, _ := utils.GetSubnetListFromIDs(sls.Logger, sls.Backend, subnetIDs)
	storageService := sls.Backend.GetNetworkStorageService().ID(volumeID)
	result, err := storageService.AllowAccessFromSubnetList(subnetList)
	sls.Logger.Info("Exit updateSubnetAuthorization ", zap.Bool("result", result), zap.Error(err))
	return err
}

func (sls *SLSession) updateHostIPAuthorization(volumeID int, hostIPList []string) error {
	sls.Logger.Info("Entry updateHostIPAuthorization", zap.Reflect("hostIPList", hostIPList))
	subnetIPAddressList, _ := utils.GetSubnetIpAddressListFromIPs(sls.Logger, sls.Backend, hostIPList)
	storageService := sls.Backend.GetNetworkStorageService().ID(volumeID)
	result, err := storageService.AllowAccessFromIpAddressList(subnetIPAddressList)
	sls.Logger.Info("Exit updateHostIPAuthorization ", zap.Reflect("result", result), zap.Error(err))
	return err
}
