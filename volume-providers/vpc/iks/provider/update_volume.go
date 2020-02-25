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
	vpc_provider "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"time"
)

const (
	customProfile = "custom"
	minSize       = 10
)

// UpdateVolume updates the volume with given information
func (vpcIks *IksVpcSession) UpdateVolume(volumeRequest provider.Volume) (err error) {
	vpcIks.Logger.Debug("Entry of UpdateVolume method...")
	defer vpcIks.Logger.Debug("Exit from UpdateVolume method...")
	defer metrics.UpdateDurationFromStart(vpcIks.Logger, "UpdateVolume", time.Now())

	vpcIks.Logger.Info("Basic validation for UpdateVolume request... ", zap.Reflect("RequestedVolumeDetails", volumeRequest))

	// Build the template to send to backend
	volumeTemplate := models.NewVolume(volumeRequest)
	err = validateVolumeRequest(volumeRequest)
	if err != nil {
		return err
	}
	vpcIks.Logger.Info("Successfully validated inputs for UpdateVolume request... ")

	vpcIks.Logger.Info("Calling  provider for volume update...")
	err = vpcIks.APIRetry.FlexyRetry(vpcIks.Logger, func() (error, bool) {
		err = vpcIks.IksSession.Apiclient.VolumeService().UpdateVolume(&volumeTemplate, vpcIks.Logger)
		return err, err == nil || vpc_provider.SkipRetryForIKS(err)
	})

	if err != nil {
		vpcIks.Logger.Debug("Failed to update volume", zap.Reflect("BackendError", err))
		return userError.GetUserError("UpdateFailed", err)
	}

	return err
}

// validateVolumeRequest validating volume request
func validateVolumeRequest(volumeRequest provider.Volume) error {

	// Volume name should not be empty
	if len(volumeRequest.VolumeID) == 0 {
		return userError.GetUserError("ErrorRequiredFieldMissing", nil, "VolumeID")
	}
	// Provider name should not be empty
	if len(volumeRequest.Provider) == 0 {
		return userError.GetUserError("ErrorRequiredFieldMissing", nil, "Provider")
	}
	// VolumeType  should not be empty
	if len(volumeRequest.VolumeType) == 0 {
		return userError.GetUserError("ErrorRequiredFieldMissing", nil, "VolumeType")
	}

	return nil
}
