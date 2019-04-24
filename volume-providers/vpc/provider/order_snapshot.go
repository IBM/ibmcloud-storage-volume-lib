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

// OrderSnapshot order snapshot
func (vpcs *VPCSession) OrderSnapshot(volumeRequest provider.Volume) error {
	vpcs.Logger.Info("Entry OrderSnapshot", zap.Reflect("volumeRequest", volumeRequest))
	defer vpcs.Logger.Info("Exit OrderSnapshot", zap.Reflect("volumeRequest", volumeRequest))

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
		return userError.GetUserError("StorageFindFailedWithVolumeId", err, volumeRequest.VolumeID, "Not a valid volume ID")
	}
	vpcs.Logger.Info("Successfully retrieved given volume details from VPC provider", zap.Reflect("VolumeDetails", volume))

	err = retry(vpcs.Logger, func() error {
		snapshot, err = vpcs.Apiclient.SnapshotService().CreateSnapshot(volumeRequest.VolumeID, snapshot, vpcs.Logger)
		return err
	})
	if err != nil {
		return userError.GetUserError("SnapshotSpaceOrderFailed", err)
	}

	vpcs.Logger.Info("Successfully created the snapshot with backend (vpcclient) call.", zap.Reflect("Snapshot", snapshot))
	return nil
}
