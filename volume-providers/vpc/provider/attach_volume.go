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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"time"

	"go.uber.org/zap"
)

//VpcVolumeAttachment ...
const (
	VpcVolumeAttachment = "vpcVolumeAttachment"
	StatusAttached      = "attached"
	StatusAttaching     = "attaching"
	StatusDetaching     = "detaching"
)

// AttachVolume attach volume based on given volume attachment request
func (vpcs *VPCSession) AttachVolume(volumeAttachmentRequest provider.VolumeAttachmentRequest) (*provider.VolumeAttachmentResponse, error) {
	vpcs.Logger.Debug("Entry of AttachVolume method...")
	defer vpcs.Logger.Debug("Exit from AttachVolume method...")
	defer metrics.UpdateDurationFromStart(vpcs.Logger, "AttachVolume", time.Now())
	var err error
	vpcs.Logger.Info("Validating basic inputs for Attach method...", zap.Reflect("volumeAttachRequest", volumeAttachmentRequest))
	err = vpcs.validateAttachVolumeRequest(volumeAttachmentRequest)
	if err != nil {
		return nil, err
	}
	var volumeAttachResult *models.VolumeAttachment
	var varp *provider.VolumeAttachmentResponse
	volumeAttachment := models.NewVolumeAttachment(volumeAttachmentRequest)

	err = vpcs.APIRetry.FlexyRetry(vpcs.Logger, func() (error, bool) {
		// First , check if volume is already attached or attaching to given instance
		vpcs.Logger.Info("Checking if volume is already attached by other thread")
		currentVolAttachment, err := vpcs.GetVolumeAttachment(volumeAttachmentRequest)
		if err == nil && currentVolAttachment != nil && currentVolAttachment.Status != StatusDetaching {
			vpcs.Logger.Info("Volume is already attached", zap.Reflect("currentVolAttachment", currentVolAttachment))
			varp = currentVolAttachment
			return nil, true // stop retry volume already attached
		}
		//Try attaching volume if it's not already attached or there is error in getting current volume attachment
		vpcs.Logger.Info("Attaching volume from VPC provider...", zap.Bool("IKSEnabled?", vpcs.Config.IsIKS))
		volumeAttachResult, err = vpcs.APIClientVolAttachMgr.AttachVolume(&volumeAttachment, vpcs.Logger)
		// Keep retry, until we get the proper volumeAttachResult object
		if err != nil {
			return err, skipRetryForObiviousErrors(err, vpcs.Config.IsIKS)
		}
		varp = volumeAttachResult.ToVolumeAttachmentResponse(vpcs.Config.VPCBlockProviderType)
		return err, true // stop retry as no error
	})

	if err != nil {
		userErr := userError.GetUserError(string(userError.VolumeAttachFailed), err, volumeAttachmentRequest.VolumeID, volumeAttachmentRequest.InstanceID)
		return nil, userErr
	}
	vpcs.Logger.Info("Successfully attached volume from VPC provider", zap.Reflect("volumeResponse", varp))
	return varp, nil
}

// validateVolume validating volume ID
func (vpcs *VPCSession) validateAttachVolumeRequest(volumeAttachRequest provider.VolumeAttachmentRequest) error {
	var err error
	// Check for InstanceID - required validation
	if len(volumeAttachRequest.InstanceID) == 0 {
		err = userError.GetUserError(string(reasoncode.ErrorRequiredFieldMissing), nil, "InstanceID")
		vpcs.Logger.Error("volumeAttachRequest.InstanceID is required", zap.Error(err))
		return err
	}
	// Check for VolumeID - required validation
	if len(volumeAttachRequest.VolumeID) == 0 {
		err = userError.GetUserError(string(reasoncode.ErrorRequiredFieldMissing), nil, "VolumeID")
		vpcs.Logger.Error("volumeAttachRequest.VolumeID is required", zap.Error(err))
		return err
	}
	return nil
}
