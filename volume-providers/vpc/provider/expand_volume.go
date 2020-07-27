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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/metrics"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"time"
)

// ExpandVolume Get the volume by using ID
func (vpcs *VPCSession) ExpandVolume(expandVolumeRequest provider.ExpandVolumeRequest) (volumeResponse *provider.Volume, err error) {
	vpcs.Logger.Debug("Entry of ExpandVolume method...")
	defer vpcs.Logger.Debug("Exit from ExpandVolume method...")
	defer metrics.UpdateDurationFromStart(vpcs.Logger, "ExpandVolume", time.Now())

	vpcs.Logger.Info("Basic validation for ExpandVolume request... ", zap.Reflect("RequestedVolumeDetails", expandVolumeRequest))
	status, err := validateExpandVolumeRequest(expandVolumeRequest)
	if !status && err != nil {
		return nil, err
	}
	vpcs.Logger.Info("Successfully validated inputs for ExpandVolume request... ")

	// Build the template to send to backend
	volumeTemplate := &models.Volume{
		Capacity: int64(*expandVolumeRequest.Capacity),
	}

	vpcs.Logger.Info("Calling VPC provider for volume expand...")
	var volume *models.Volume
	err = retry(vpcs.Logger, func() error {
		volume, err = vpcs.Apiclient.VolumeService().ExpandVolume(volumeTemplate, vpcs.Logger)
		return err
	})

	if err != nil {
		vpcs.Logger.Debug("Failed to expand volume from VPC provider", zap.Reflect("BackendError", err))
		return nil, userError.GetUserError("FailedToPlaceOrder", err)
	}

	vpcs.Logger.Info("Successfully expanded volume from VPC provider...", zap.Reflect("VolumeDetails", volume))

	vpcs.Logger.Info("Waiting for volume to be in valid (available) state", zap.Reflect("VolumeDetails", volume))
	err = WaitForValidVolumeState(vpcs, volume.ID)
	if err != nil {
		return nil, userError.GetUserError("VolumeNotInValidState", err, volume.ID)
	}
	vpcs.Logger.Info("Volume got valid (available) state", zap.Reflect("VolumeDetails", volume))

	// Converting volume to lib volume type
	volumeResponse = FromProviderToLibVolume(volume, vpcs.Logger)
	vpcs.Logger.Info("VolumeResponse", zap.Reflect("volumeResponse", volumeResponse))
	return volumeResponse, err
}

// validateExpandVolumeRequest validating volume request
func validateExpandVolumeRequest(volumeRequest provider.ExpandVolumeRequest) (bool, error) {
	// Volume name should not be empty
	if volumeRequest.Name == nil {
		return false, userError.GetUserError("InvalidVolumeName", nil, nil)
	} else if len(*volumeRequest.Name) == 0 {
		return false, userError.GetUserError("InvalidVolumeName", nil, *volumeRequest.Name)
	}

	// Capacity should not be empty
	if volumeRequest.Capacity == nil {
		return false, userError.GetUserError("VolumeCapacityInvalid", nil, nil)
	} else if *volumeRequest.Capacity < minSize {
		return false, userError.GetUserError("VolumeCapacityInvalid", nil, *volumeRequest.Capacity)
	}
	return true, nil
}
