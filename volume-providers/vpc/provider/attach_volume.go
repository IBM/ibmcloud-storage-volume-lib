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

const VpcVolumeAttachment = "vpcVolumeAttachment"

// Attach volume
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
	}
	vpcs.Logger.Info("Attaching volume from VPC provider...")
	var volumeAttachResult *models.VolumeAttachment
	err = retry(func() error {
		volumeAttachResult, err = vpcs.Apiclient.VolumeMountService().AttachVolume(&volumeAttachment)
		return err
	})

	if err != nil {
		userErr := userError.GetUserError("FailedToAtachVolume", err, volumeAttachRequest.VPCVolumeAttachment.Volume.VolumeID)
		volumeResponse.Message = userErr.Error()
		return volumeResponse, userErr

	}

	volumeResponse.Status = provider.SUCCESS
	vpcs.Logger.Info("Successfully attached volume from VPC provider", zap.Reflect("volumeAttachResult", volumeAttachResult))
	return volumeResponse, nil
}

// validateVolume validating volume ID
func validateAttachVolumeRequest(volumeAttachRequest provider.VolumeAttachRequest) error {
	var err error

	if volumeAttachRequest.VPCVolumeAttachment == nil {
		err = userError.GetUserError("EmptyVPCVolumeAttachment", nil, nil)
		return err
	}
	if len(volumeAttachRequest.VPCVolumeAttachment.Name) < 1 {
		err = userError.GetUserError("EmptyVPCVolumeAttachmentParameter", nil, "Name")
		return err
	}
	if volumeAttachRequest.VPCVolumeAttachment.Volume == nil {
		err = userError.GetUserError("EmptyVPCVolumeAttachmentParameter", nil, "Volume")
		return err
	}
	return nil
}
