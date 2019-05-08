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

	"errors"
	"go.uber.org/zap"
)

// GetVolumeAttachment  get the volume attachment based on the request
func (vpcs *VPCSession) GetVolumeAttachment(volumeAttachmentRequest provider.VolumeAttachmentRequest) (*provider.VolumeAttachmentResponse, error) {
	vpcs.Logger.Debug("Entry of GetVolumeAttachment method...")
	defer vpcs.Logger.Debug("Exit from GetVolumeAttachment method...")
	var err error
	vpcs.Logger.Info("Validating basic inputs for GetVolumeAttachment method...", zap.Reflect("volumeAttachRequest", volumeAttachmentRequest))
	err = vpcs.validateAttachVolumeRequest(volumeAttachmentRequest)
	if err != nil {
		return nil, err
	}
	volumeAttachment := models.NewVolumeAttachment(volumeAttachmentRequest)
	vpcs.Logger.Info("Getting VolumeAttachmentList from VPC provider...")
	var volumeAttachmentList *models.VolumeAttachmentList
	err = retry(vpcs.Logger, func() error {
		volumeAttachmentList, err = vpcs.Apiclient.VolumeAttachService().ListVolumeAttachment(&volumeAttachment, vpcs.Logger)
		return err
	})
	if err != nil {
		// API call is failed
		userErr := userError.GetUserError(string(userError.VolumeAttachFindFailed), err, volumeAttachmentRequest.VolumeID, volumeAttachmentRequest.InstanceID)
		return nil, userErr
	}
	// Iterate over the volume attachment list for given instance
	for _, volumeAttachmentItem := range volumeAttachmentList.VolumeAttachments {
		// Check if volume ID is matching with requested volume ID
		if volumeAttachmentItem.Volume.ID == volumeAttachmentRequest.VolumeID {
			vpcs.Logger.Info("Successfully found volume attachment", zap.Reflect("volumeAttachment", volumeAttachmentItem))
			volumeResponse := volumeAttachmentItem.ToVolumeAttachmentResponse()
			vpcs.Logger.Info("Successfully fetched volume attachment from VPC provider", zap.Reflect("volumeResponse", volumeResponse))
			return volumeResponse, nil
		}
	}
	// No volume attahment found in the  list. So return error
	userErr := userError.GetUserError(string(userError.VolumeAttachFindFailed), errors.New("No VolumeAttachment Found"), volumeAttachmentRequest.VolumeID, volumeAttachmentRequest.InstanceID)
	vpcs.Logger.Error("Volume attachment not found", zap.Error(err))
	return nil, userErr

}
