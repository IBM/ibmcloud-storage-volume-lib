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

const (
	customProfile = "custom"
	minSize       = 10
)

// CreateVolume Get the volume by using ID
func (vpcs *VPCSession) CreateVolume(volumeRequest provider.Volume) (volumeResponse *provider.Volume, err error) {
	vpcs.Logger.Debug("Entry of CreateVolume method...")
	defer vpcs.Logger.Debug("Exit from CreateVolume method...")

	vpcs.Logger.Info("Basic validation for CreateVolume request... ", zap.Reflect("RequestedVolumeDetails", volumeRequest))
	resourceGroup, iops, err := validateVolumeRequest(volumeRequest)
	if err != nil {
		return nil, err
	}
	vpcs.Logger.Info("Successfully validated inputs for CreateVolume request... ")

	// Build the template to send to backend
	volumeTemplate := &models.Volume{
		Name:          *volumeRequest.Name,
		Capacity:      int64(*volumeRequest.Capacity),
		Iops:          iops,
		Tags:          volumeRequest.VPCVolume.Tags,
		ResourceGroup: &resourceGroup,
		Generation:    models.GenerationType(vpcs.Config.VPCBlockProviderName),
		Profile: &models.Profile{
			Name: volumeRequest.VPCVolume.Profile.Name,
		},
		Zone: &models.Zone{
			Name: volumeRequest.Az,
		},
	}

	var encryptionKeyCRN string
	if volumeRequest.VPCVolume.VolumeEncryptionKey != nil && len(volumeRequest.VPCVolume.VolumeEncryptionKey.CRN) > 0 {
		encryptionKeyCRN = volumeRequest.VPCVolume.VolumeEncryptionKey.CRN
		volumeTemplate.VolumeEncryptionKey = &models.VolumeEncryptionKey{CRN: encryptionKeyCRN}
	}

	vpcs.Logger.Info("Calling VPC provider for volume creation...")
	var volume *models.Volume
	err = retry(vpcs.Logger, func() error {
		volume, err = vpcs.Apiclient.VolumeService().CreateVolume(volumeTemplate, vpcs.Logger)
		return err
	})

	if err != nil {
		vpcs.Logger.Debug("Failed to create volume from VPC provider", zap.Reflect("BackendError", err))
		return nil, userError.GetUserError("FailedToPlaceOrder", err)
	}

	vpcs.Logger.Info("Successfully created volume from VPC provider...", zap.Reflect("VolumeDetails", volume))

	// Converting volume to lib volume type
	volumeResponse = FromProviderToLibVolume(volume, vpcs.Logger)
	return volumeResponse, err
}

// validateVolumeRequest validating volume request
func validateVolumeRequest(volumeRequest provider.Volume) (models.ResourceGroup, int64, error) {
	resourceGroup := models.ResourceGroup{}
	var iops int64
	iops = 0
	// Volume name should not be empty
	if volumeRequest.Name == nil {
		return resourceGroup, iops, userError.GetUserError("InvalidVolumeName", nil, nil)
	} else if len(*volumeRequest.Name) == 0 {
		return resourceGroup, iops, userError.GetUserError("InvalidVolumeName", nil, *volumeRequest.Name)
	}

	// Capacity should not be empty
	if volumeRequest.Capacity == nil {
		return resourceGroup, iops, userError.GetUserError("VolumeCapacityInvalid", nil, nil)
	} else if *volumeRequest.Capacity < minSize {
		return resourceGroup, iops, userError.GetUserError("VolumeCapacityInvalid", nil, *volumeRequest.Capacity)
	}

	// Read user provided error, no harm to pass the 0 values to RIaaS in case of tiered profiles
	if volumeRequest.Iops != nil {
		iops = ToInt64(*volumeRequest.Iops)
	}
	if volumeRequest.VPCVolume.Profile.Name != customProfile && iops > 0 {
		return resourceGroup, iops, userError.GetUserError("VolumeProfileIopsInvalid", nil)
	}

	// validate and add resource group ID or Name whichever is provided by user
	if volumeRequest.VPCVolume.ResourceGroup == nil {
		return resourceGroup, iops, userError.GetUserError("EmptyResourceGroup", nil)
	}

	// validate and add resource group ID or Name whichever is provided by user
	if len(volumeRequest.VPCVolume.ResourceGroup.ID) == 0 && len(volumeRequest.VPCVolume.ResourceGroup.Name) == 0 {
		return resourceGroup, iops, userError.GetUserError("EmptyResourceGroupIDandName", nil)
	}

	if len(volumeRequest.VPCVolume.ResourceGroup.ID) > 0 {
		resourceGroup.ID = volumeRequest.VPCVolume.ResourceGroup.ID
	}
	if len(volumeRequest.VPCVolume.ResourceGroup.Name) > 0 {
		// get the resource group ID from resource group name as Name is not supported by RIaaS
		resourceGroup.Name = volumeRequest.VPCVolume.ResourceGroup.Name
	}
	return resourceGroup, iops, nil
}
