/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/client/models"
	"go.uber.org/zap"
)

func (vpcs *VPCSession) OrderSnapshot(volumeRequest provider.Volume) error {
	var snapshot *models.Snapshot
	var err error

	// Step 1- validate input which are required
	vpcs.Logger.Info("Requested volume is:", zap.Reflect("Volume", volumeRequest))
	var volume *models.Volume

	err = retry(func() error {
		volume, err = vpcs.Apiclient.Volume().GetVolume(volumeRequest.VolumeID)
		return err
	})
	if err != nil {
		return reasoncode.GetUserError("StorageFindFailedWithVolumeId", err, volumeRequest.VolumeID, "Not a valid volume ID")
	}

	err = retry(func() error {
		snapshot, err = vpcs.Apiclient.Snapshot().CreateSnapshot(volumeRequest.VolumeID, snapshot)
		return err
	})
	if err != nil {
		return reasoncode.GetUserError("SnapshotSpaceOrderFailed", err)
	}

	vpcs.Logger.Info("Backend created snapshot details", zap.Reflect("Snapshot", snapshot))

	return nil
}
