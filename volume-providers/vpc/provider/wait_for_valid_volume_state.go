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

// WaitForValidVolumeState checks the volume for valid status
func WaitForValidVolumeState(vpcs *VPCSession, volumeID string) (err error) {
	vpcs.Logger.Debug("Entry of WaitForValidVolumeState method...")
	defer vpcs.Logger.Debug("Exit from WaitForValidVolumeState method...")

	vpcs.Logger.Info("Getting volume details from VPC provider...", zap.Reflect("VolumeID", volumeID))

	var volume *models.Volume
	err = retry(vpcs.Logger, func() error {
		volume, err = vpcs.Apiclient.VolumeService().GetVolume(volumeID, vpcs.Logger)
		vpcs.Logger.Info("Getting volume details from VPC provider...", zap.Reflect("volume", volume))
		if volume != nil && volume.Status == validVolumeStatus {
			vpcs.Logger.Info("Volume got valid (available) state", zap.Reflect("VolumeDetails", volume))
			return nil
		}
		return err
	})

	if err != nil {
		vpcs.Logger.Info("Volume could not get valid (available) state", zap.Reflect("VolumeDetails", volume))
		return userError.GetUserError("VolumeNotInValidState", err, volumeID)
	}

	return nil
}
