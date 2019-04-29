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
	//"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"

	"go.uber.org/zap"
)

//VpcVolumeAttachment ...
const VpcVolumeAttachment = "vpcVolumeAttachment"

// Attach volume baed on given volume attachment request
func (vpcs *VPCSession) Attach(volumeAttachRequest provider.VolumeAttachRequest) (provider.VolumeResponse, error) {
	vpcs.Logger.Debug("Entry of Attach method...")
	defer vpcs.Logger.Debug("Exit from Attach method...")
	var err error
	vpcs.Logger.Info("Validating basic inputs for Attach method...", zap.Reflect("volumeAttachRequest", volumeAttachRequest))
	err = validateAttachVolumeRequest(volumeAttachRequest)
	volumeResponse := provider.VolumeResponse{}
	if err != nil {
		volumeResponse.Status = provider.FAILURE
		volumeResponse.Message = err.Error()
		return volumeResponse, err
	}
	volumeAttachment := models.VolumeAttachment{
		VolumeAttachment: *volumeAttachRequest.VPCVolumeAttachment,
		Volume: &models.Volume{
			ID: volumeAttachRequest.VPCVolumeAttachment.Volume.VolumeID,
		},
	}
	var volumeAttachResult *models.VolumeAttachment
	// First , check if volume is already attached to given instance
	vpcs.Logger.Info("Checking if volume is already attached ")
	currentVolAttachment, _ := vpcs.GetAttachStatus(volumeAttachRequest)
	if currentVolAttachment.Status == provider.SUCCESS && currentVolAttachment.VPCVolumeAttachment != nil {
		vpcs.Logger.Info("volume is already attached", zap.Reflect("currentVolAttachment", currentVolAttachment))
		return currentVolAttachment, nil
	}
	//Try attaching volume if it's not already attached or there is error in getting current volume attachment
	vpcs.Logger.Info("Attaching volume from VPC provider...")
	err = retry(func() error {
		volumeAttachResult, err = vpcs.Apiclient.VolumeMountService().AttachVolume(&volumeAttachment, vpcs.Logger)
		return err
	})
	if err != nil {
		userErr := userError.GetUserError(string(userError.VolumeAttachFailed), err, volumeAttachRequest.VPCVolumeAttachment.Volume.VolumeID, volumeAttachRequest.VPCVolumeAttachment.InstanceID)
		volumeResponse.Message = userErr.Error()
		return volumeResponse, userErr
	}
	volumeResponse.Status = provider.SUCCESS
	volumeResponse.VPCVolumeAttachment = &volumeAttachResult.VolumeAttachment
	vpcs.Logger.Info("Successfully attached volume from VPC provider", zap.Reflect("volumeResponse", volumeResponse))
	return volumeResponse, nil
}

// validateVolume validating volume ID
func validateAttachVolumeRequest(volumeAttachRequest provider.VolumeAttachRequest) error {
	var err error
	vpcVolumeAttachmentMissing := volumeAttachRequest.VPCVolumeAttachment == nil
	volumeMissing := vpcVolumeAttachmentMissing || volumeAttachRequest.VPCVolumeAttachment.Volume == nil
	// Check for InstanceID - required validation
	if vpcVolumeAttachmentMissing || len(volumeAttachRequest.VPCVolumeAttachment.InstanceID) == 0 {
		err = userError.GetUserError(string(reasoncode.ErrorRequiredFieldMissing), nil, "VolumeAttachRequest.VPCVolumeAttachment.InstanceID")
		return err
	}
	// Check for VolumeID - required validation
	if volumeMissing || len(volumeAttachRequest.VPCVolumeAttachment.Volume.VolumeID) == 0 {
		err = userError.GetUserError(string(reasoncode.ErrorRequiredFieldMissing), nil, "VolumeAttachRequest.VPCVolumeAttachment.Volume.VolumeID")
		return err
	}
	return nil
}
