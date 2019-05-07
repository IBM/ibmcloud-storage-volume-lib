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
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
)

// CreateSnapshot Create snapshot from given volume
func (vpcs *VPCSession) CreateSnapshot(volumeRequest *provider.Volume, tags map[string]string) (*provider.Snapshot, error) {
	vpcs.Logger.Info("Entry CreateSnapshot", zap.Reflect("volumeRequest", volumeRequest))
	defer vpcs.Logger.Info("Exit CreateSnapshot", zap.Reflect("volumeRequest", volumeRequest))

	if volumeRequest == nil {
		return nil, userError.GetUserError("StorageFindFailedWithVolumeId", nil, "Not a valid volume ID")
	}

	var snapshot *models.Snapshot
	var err error

	// Step 1- validate input which are required
	vpcs.Logger.Info("Requested volume is:", zap.Reflect("Volume", volumeRequest))

	var volume *models.Volume
	err = retry(vpcs.Logger, func() error {
		volume, err = vpcs.Apiclient.VolumeService().GetVolume(volumeRequest.VolumeID, vpcs.Logger)
		return err
	})
	if err != nil {
		return nil, userError.GetUserError("StorageFindFailedWithVolumeId", err, "Not a valid volume ID")
	}

	if volume == nil {
		return nil, userError.GetUserError("StorageFindFailedWithVolumeId", err, volumeRequest.VolumeID, "Not a valid volume ID")
	}

	err = retry(vpcs.Logger, func() error {
		snapshot, err = vpcs.Apiclient.SnapshotService().CreateSnapshot(volumeRequest.VolumeID, snapshot, vpcs.Logger)
		return err
	})
	if err != nil {
		return nil, userError.GetUserError("SnapshotSpaceOrderFailed", err)
	}

	vpcs.Logger.Info("Successfully created snapshot with backend (vpcclient) call")
	vpcs.Logger.Info("Backend created snapshot details", zap.Reflect("Snapshot", snapshot))

	// Converting volume to lib volume type
	volumeResponse := FromProviderToLibVolume(volume, vpcs.Logger)
	if volumeResponse != nil {
		respSnapshot := &provider.Snapshot{
			Volume:               *volumeResponse,
			SnapshotID:           snapshot.ID,
			SnapshotCreationTime: *snapshot.CreatedAt,
		}
		return respSnapshot, nil
	}

	return nil, userError.GetUserError("CoversionNotSuccessful", err, "Not able to prepare provider volume")
}
