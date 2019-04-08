/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
)

// CreateSnapshot Create snapshot from given volume
func (vpcs *VPCSession) CreateSnapshot(volumeRequest *provider.Volume, tags map[string]string) (*provider.Snapshot, error) {
	vpcs.Logger.Info("Entry CreateSnapshot", zap.Reflect("volumeRequest", volumeRequest))
	defer vpcs.Logger.Info("Exit CreateSnapshot", zap.Reflect("volumeRequest", volumeRequest))

	var snapshot *models.Snapshot
	var err error

	// Step 1- validate input which are required
	vpcs.Logger.Info("Requested volume is:", zap.Reflect("Volume", volumeRequest))
	var volume *provider.Volume

	err = retry(func() error {
		volume, err = vpcs.GetVolume(volumeRequest.VolumeID)
		return err
	})
	if err != nil {
		vpcs.Logger.Info("FAILED: Not a valid volume ID")
		return nil, reasoncode.GetUserError("StorageFindFailedWithVolumeId", err, volumeRequest.VolumeID, "Not a valid volume ID")
	}

	err = retry(func() error {
		snapshot, err = vpcs.Apiclient.SnapshotService().CreateSnapshot(volumeRequest.VolumeID, snapshot)
		return err
	})
	if err != nil {
		vpcs.Logger.Info("FAILED: Failed to create snapshot with backend (vpcclient) call")
		return nil, reasoncode.GetUserError("SnapshotSpaceOrderFailed", err)
	}

	vpcs.Logger.Info("SUCCESS: Successfully created snapshot with backend (vpcclient) call")
	vpcs.Logger.Info("Backend created snapshot details", zap.Reflect("Snapshot", snapshot))

	respSnapshot := &provider.Snapshot{
		Volume:               *volume,
		SnapshotID:           snapshot.ID,
		SnapshotCreationTime: *snapshot.CreatedAt,
	}
	return respSnapshot, nil
}
