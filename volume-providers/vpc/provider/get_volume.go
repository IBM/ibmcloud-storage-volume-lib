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

// GetVolume gets the volume by using ID
func (vpcs *VPCSession) GetVolume(id string) (respVolume *provider.Volume, err error) {
	vpcs.Logger.Debug("Entry of GetVolume method...")
	defer vpcs.Logger.Debug("Exit from GetVolume method...")

	vpcs.Logger.Info("Basic validation for volume ID...", zap.Reflect("VolumeID", id))
	// validating volume ID
	err = validateVolumeID(id)
	if err != nil {
		return nil, err
	}

	vpcs.Logger.Info("Getting volume details from VPC provider...", zap.Reflect("VolumeID", id))

	var volume *models.Volume
	err = retry(func() error {
		volume, err = vpcs.Apiclient.VolumeService().GetVolume(id, vpcs.Logger)
		return err
	})

	if err != nil {
		return nil, userError.GetUserError("StorageFindFailedWithVolumeId", err, id)
	}

	vpcs.Logger.Info("Successfully retrieved volume details from VPC provider", zap.Reflect("VolumeDetails", volume))

	// Converting volume to lib volume type
	respVolume = FromProviderToLibVolume(volume, vpcs.Logger)
	return respVolume, err
}

// validateVolumeID validating basic volume ID
func validateVolumeID(volumeID string) (err error) {
	if IsValidVolumeIDFormat(volumeID) {
		return nil
	}
	err = userError.GetUserError("InvalidVolumeID", nil, volumeID)
	return
}
