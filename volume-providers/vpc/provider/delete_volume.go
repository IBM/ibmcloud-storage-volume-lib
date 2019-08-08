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
	err = retry(vpcs.Logger, func() error {
		err = vpcs.Apiclient.VolumeService().DeleteVolume(volume.VolumeID, vpcs.Logger)
		return err
	})
	if err != nil {
		return userError.GetUserError("FailedToDeleteVolume", err, volume.VolumeID)
	}

	err = WaitForVolumeDeletion(vpcs, volume.VolumeID)
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

// WaitForVolumeDeletion checks the volume for valid status
func WaitForVolumeDeletion(vpcs *VPCSession, volumeID string) (err error) {
	vpcs.Logger.Debug("Entry of WaitForVolumeDeletion method...")
	defer vpcs.Logger.Debug("Exit from WaitForVolumeDeletion method...")

	vpcs.Logger.Info("Getting volume details from VPC provider...", zap.Reflect("VolumeID", volumeID))

	err = vpcs.APIRetry.FlexyRetry(vpcs.Logger, func() (error, bool) {
		_, err = vpcs.Apiclient.VolumeService().GetVolume(volumeID, vpcs.Logger)
		// Keep retry, until GetVolume returns volume not found
		if err != nil {
			return err, skipRetry(err.(*models.Error))
		}
		return err, false // continue retry as we are not seeing error which means volume is available
	})

	if err == nil {
		vpcs.Logger.Info("Volume got deleted.", zap.Reflect("volumeID", volumeID))
	}
	return err
}
