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
	"strconv"
)

// CreateVolume Get the volume by using ID
func (vpcs *VPCSession) CreateVolume(volumeRequest provider.Volume, ctxLogger *zap.Logger) (volumeResponse *provider.Volume, err error) {
	ctxLogger.Debug("Entry of CreateVolume method...")
	defer ctxLogger.Debug("Exit from CreateVolume method...")

	ctxLogger.Info("Basic validation for CreateVolume request... ", zap.Reflect("RequestedVolumeDetails", volumeRequest))
	err = validateVolumeRequest(volumeRequest)
	if err != nil {
		return nil, err
	}
	ctxLogger.Info("Successfully validated inputs for CreateVolume request... ")
	// Pending error handling
	// TODO: Check if the volume already exists with same name.
	// We can do this by scanning all volumes. But requesting the VPC team to get
	// an API for getting volume details with name instead of only ID.

	iops := ToInt64(*volumeRequest.Iops)

	// Build the template we'll send to RIAAS
	volumeTemplate := &models.Volume{
		Name:     *volumeRequest.Name,
		Capacity: int64(*volumeRequest.Capacity),
		Iops:     iops,
		Tags:     volumeRequest.VPCVolume.Tags,
		ResourceGroup: &models.ResourceGroup{
			ID: volumeRequest.VPCVolume.ResourceGroup.ID,
		},
		Generation: models.GenerationType(vpcs.Config.VPCBlockProviderName),
		Profile: &models.Profile{
			Name: volumeRequest.VPCVolume.Profile.Name,
		},
		Zone: &models.Zone{
			Name: volumeRequest.Az,
		},
	}

	ctxLogger.Info("Calling VPC provider for volume creation...")
	var volume *models.Volume
	err = retry(func() error {
		volume, err = vpcs.Apiclient.VolumeService().CreateVolume(volumeTemplate, ctxLogger)
		return err
	})

	if err != nil {
		ctxLogger.Debug("Failed to create volume from VPC provider", zap.Reflect("BackendError", err))
		return nil, userError.GetUserError("FailedToPlaceOrder", err)
	}

	ctxLogger.Info("Successfully created volume from VPC provider...", zap.Reflect("VolumeDetails", volume))

	// Converting volume to lib volume type
	volumeResponse = FromProviderToLibVolume(volume, vpcs.Logger)
	return volumeResponse, err
}

// validateVolumeRequest validating volume request
func validateVolumeRequest(volumeRequest provider.Volume) (err error) {
	// Volume name should not be empty
	if volumeRequest.Name == nil {
		return userError.GetUserError("InvalidVolumeName", nil, nil)
	} else if len(*volumeRequest.Name) == 0 {
		return userError.GetUserError("InvalidVolumeName", nil, *volumeRequest.Name)
	}

	// Capacity should not be empty
	if volumeRequest.Capacity == nil {
		return userError.GetUserError("VolumeCapacityInvalid", nil, nil)
	} else if *volumeRequest.Capacity <= 0 {
		return userError.GetUserError("VolumeCapacityInvalid", nil, *volumeRequest.Capacity)
	}

	// General purpose profiles does not allow IOPs setting
	if volumeRequest.VPCVolume.Profile.Name != "general-purpose" && (volumeRequest.Iops == nil || *volumeRequest.Iops <= strconv.Itoa(0)) {
		return userError.GetUserError("IopsInvalid", nil, *volumeRequest.Iops)
	}

	// General purpose profiles does not allow IOPs setting
	if *volumeRequest.Iops > strconv.Itoa(0) && volumeRequest.VPCVolume.Profile.Name == "general-purpose" {
		return userError.GetUserError("VolumeProfileIopsInvalid", nil)
	}
	return nil
}
