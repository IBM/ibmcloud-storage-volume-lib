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

// GetAttachStatus volume
func (vpcs *VPCSession) GetAttachStatus(volumeAttachRequest provider.VolumeAttachRequest) (provider.VolumeResponse, error) {
	vpcs.Logger.Debug("Entry of GetAttachStatus method...")
	defer vpcs.Logger.Debug("Exit from GetAttachStatus method...")
	var err error
	vpcs.Logger.Info("Validating basic inputs for GetAttachStatus method...", zap.Reflect("volumeAttachRequest", volumeAttachRequest))
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
	vpcs.Logger.Info("Getting Attach  status from VPC provider...")
	var volumeAttachResult *models.VolumeAttachment
	err = retry(vpcs.Logger, func() error {
		volumeAttachResult, err = vpcs.Apiclient.VolumeMountService().GetAttachStatus(&volumeAttachment, vpcs.Logger)
		return err
	})
	if err != nil {
		userErr := userError.GetUserError(string(userError.VolumeAttachFindFailed), err, volumeAttachRequest.VPCVolumeAttachment.Volume.VolumeID, volumeAttachRequest.VPCVolumeAttachment.InstanceID)
		volumeResponse.Message = userErr.Error()
		return volumeResponse, userErr
	}
	volumeResponse.Status = provider.SUCCESS
	volumeResponse.VPCVolumeAttachment = &volumeAttachResult.VolumeAttachment
	volumeResponse.VPCVolumeAttachment.Volume = FromProviderToLibVolume(volumeAttachResult.Volume, vpcs.Logger)
	vpcs.Logger.Info("Successfully fetched volume attachment from VPC provider", zap.Reflect("volumeResponse", volumeResponse))
	return volumeResponse, nil
}
