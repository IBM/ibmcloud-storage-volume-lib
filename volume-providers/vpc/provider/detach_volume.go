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
	"net/http"
)

// Detach volume based on give volume attachment request
func (vpcs *VPCSession) Detach(volumeAttachmentTemplate provider.VolumeAttachRequest) (provider.VolumeResponse, error) {
	vpcs.Logger.Debug("Entry of Attach method...")
	defer vpcs.Logger.Debug("Exit from Attach method...")
	var err error
	vpcs.Logger.Info("Validating basic inputs for detach method...", zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate))
	err = validateAttachVolumeRequest(volumeAttachmentTemplate)
	volumeResponse := provider.VolumeResponse{}
	if err != nil {
		volumeResponse.Status = provider.FAILURE
		volumeResponse.Message = err.Error()
		return volumeResponse, err
	}

	var response *http.Response
	// First , check if volume is already attached to given instance
	vpcs.Logger.Info("Checking if volume is already attached ")
	currentVolAttachment, _ := vpcs.GetAttachStatus(volumeAttachmentTemplate)
	if currentVolAttachment.Status == provider.SUCCESS && currentVolAttachment.VPCVolumeAttachment != nil {
		vpcs.Logger.Info("Found volume attachment", zap.Reflect("currentVolAttachment", currentVolAttachment))
		if currentVolAttachment.VPCVolumeAttachment.Status != StatusDetaching {
			//Try detaching volume if it's not already in detaching
			volumeAttachment := models.VolumeAttachment{
				VolumeAttachment: *currentVolAttachment.VPCVolumeAttachment,

				Volume: &models.Volume{
					ID: currentVolAttachment.VPCVolumeAttachment.Volume.VolumeID,
				},
			}
			volumeAttachment.VolumeAttachment.InstanceID = volumeAttachmentTemplate.VPCVolumeAttachment.InstanceID
			vpcs.Logger.Info("Detaching volume from VPC provider...")
			err = retry(vpcs.Logger, func() error {
				response, err = vpcs.Apiclient.VolumeMountService().DetachVolume(&volumeAttachment, vpcs.Logger)
				return err
			})
			if err != nil {
				userErr := userError.GetUserError(string(userError.VolumeDetachFailed), err, volumeAttachmentTemplate.VPCVolumeAttachment.Volume.VolumeID, volumeAttachmentTemplate.VPCVolumeAttachment.InstanceID, volumeAttachmentTemplate.VPCVolumeAttachment.ID)
				volumeResponse.Message = userErr.Error()
				return volumeResponse, userErr
			}
			volumeResponse.Status = provider.SUCCESS
			vpcs.Logger.Info("Successfully detached volume from VPC provider", zap.Reflect("volumeResponse", volumeResponse), zap.Reflect("resp", response))
		} else {
			vpcs.Logger.Info("Volume is already getting detached")
		}
	} else {
		vpcs.Logger.Info("No volume attachment found for", zap.Reflect("currentVolAttachment", currentVolAttachment))
		volumeResponse.Status = provider.SUCCESS // consider volume detach success if its not already attached
	}
	return volumeResponse, nil
}
