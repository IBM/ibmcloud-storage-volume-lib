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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"

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
	var err error
	vpcs.Logger.Info("Validating basic inputs for Attach method...", zap.Reflect("volumeAttachRequest", volumeAttachmentRequest))
	err = vpcs.validateAttachVolumeRequest(volumeAttachmentRequest)
	if err != nil {
		return nil, err
	}
	var volumeAttachResult *models.VolumeAttachment
	// First , check if volume is already attached or attaching to given instance
	vpcs.Logger.Info("Checking if volume is already attached ")
	currentVolAttachment, err := vpcs.GetVolumeAttachment(volumeAttachmentRequest)
	if err == nil && currentVolAttachment != nil && currentVolAttachment.Status != StatusDetaching {
		vpcs.Logger.Info("volume is already attached", zap.Reflect("currentVolAttachment", currentVolAttachment))
		return currentVolAttachment, nil
	}
	//Try attaching volume if it's not already attached or there is error in getting current volume attachment
	vpcs.Logger.Info("Attaching volume from VPC provider...")
	volumeAttachment := models.NewVolumeAttachment(volumeAttachmentRequest)

	err = vpcs.APIRetry.FlexyRetry(vpcs.Logger, func() (interface{}, error) {
		volumeAttachResult, err = vpcs.APIClientVolAttachMgr.AttachVolume(&volumeAttachment, vpcs.Logger)
		return volumeAttachResult, err
	}, func(intf interface{}, err *models.Error) bool {
		// Skip retry as per common errors
		if err != nil {
			return skipRetry(err)
		}
		// stop retry, as there is no error
		return true
	})

	if err != nil {
		userErr := userError.GetUserError(string(userError.VolumeAttachFailed), err, volumeAttachmentRequest.VolumeID, volumeAttachmentRequest.InstanceID)
		return nil, userErr
	}
	varp := volumeAttachResult.ToVolumeAttachmentResponse()
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
