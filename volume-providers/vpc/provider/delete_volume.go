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
	"go.uber.org/zap"
)

// DeleteVolume deletes the volume
func (vpcs *VPCSession) DeleteVolume(volume *provider.Volume) (err error) {
	vpcs.Logger.Debug("Entry of DeleteVolume method...")
	defer vpcs.Logger.Debug("Exit from DeleteVolume method...")

	vpcs.Logger.Info("Validating basic inputs for DeleteVolume method...", zap.Reflect("VolumeDetails", volume))
	err = validateVolume(volume)
	if err != nil {
		return err
	}

	vpcs.Logger.Info("Deleting volume from VPC provider...")
	err = retry(func() error {
		err = vpcs.Apiclient.VolumeService().DeleteVolume(volume.VolumeID, vpcs.Logger)
		return err
	})

	if err != nil {
		return userError.GetUserError("FailedToDeleteVolume", err, volume.VolumeID)
	}

	vpcs.Logger.Info("Successfully deleted volume from VPC provider")
	return err
}

// validateVolume validating volume ID
func validateVolume(volume *provider.Volume) (err error) {
	if volume == nil {
		err = userError.GetUserError("InvalidVolumeID", nil, nil)
		return
	}

	if IsValidVolumeIDFormat(volume.VolumeID) {
		return nil
	}
	err = userError.GetUserError("InvalidVolumeID", nil, volume.VolumeID)
	return
}
