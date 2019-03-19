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
	"go.uber.org/zap"
)

func (vpcs *VPCSession) OrderSnapshot(volumeRequest provider.Volume) error {
	// Step 1- validate input which are required
	vpcs.Logger.Info("Requested volume is:", zap.Reflect("Volume", volumeRequest))
	if volumeRequest.SnapshotSpace == nil {
		vpcs.Logger.Error("No proper input, please provide volume ID and snapshot space size")
		return reasoncode.GetUserError("SnapshotSpaceOrderFailed", nil)
	}

	volid := ToInt(volumeRequest.VolumeID)
	snapshotSize := *volumeRequest.SnapshotSpace
	if volid == 0 || snapshotSize == 0 {
		vpcs.Logger.Error("No proper input, please provide volume ID and snapshot space size")
		return reasoncode.GetUserError("SnapshotSpaceOrderFailed", nil)
	}
	return nil
}
