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
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
)

const (
	validVolumeStatus = "available"
)

// WaitForVolumeAvailableState checks the volume for valid status
func (vpcs *VPCSession) WaitForVolumeAvailableState(volumeID string) (err error) {
	vpcs.Logger.Debug("Entry of WaitForVolumeAvailableState method...")
	defer vpcs.Logger.Debug("Exit from WaitForVolumeAvailableState method...")

	vpcs.Logger.Info("Basic validation for volume ID...", zap.Reflect("VolumeID", volumeID))
	// validating volume ID
	err = validateVolumeID(volumeID)
	if err != nil {
		return err
	}

	vpcs.Logger.Info("Getting volume details from VPC provider...", zap.Reflect("VolumeID", volumeID))

	var volume *models.Volume
	err = retry(vpcs.Logger, func() error {
		volume, err = vpcs.Apiclient.VolumeService().GetVolume(volumeID, vpcs.Logger)
		vpcs.Logger.Info("Getting volume details from VPC provider...", zap.Reflect("volume", volume))
		if volume.Status == validVolumeStatus {
			vpcs.Logger.Info("Volume got available state", zap.Reflect("VolumeDetails", volume))
			return nil
		}
		return err
	})

	if err != nil {
		vpcs.Logger.Info("Volume could not get available state", zap.Reflect("VolumeDetails", volume))
		return userError.GetUserError("StorageFindFailedWithVolumeId", err, volumeID)
	}

	return nil
}
