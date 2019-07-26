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

	//"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
)

// WaitForAttachVolume waits for volume to be attached to node. e.g waits till status becomes attached
func (vpcs *VPCSession) WaitForAttachVolume(volumeAttachmentTemplate provider.VolumeAttachmentRequest) (*provider.VolumeAttachmentResponse, error) {
	vpcs.Logger.Debug("Entry of WaitForAttachVolume method...")
	defer vpcs.Logger.Debug("Exit from WaitForAttachVolume method...")
	vpcs.Logger.Info("Validating basic inputs for WaitForAttachVolume method...", zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate))
	err := vpcs.validateAttachVolumeRequest(volumeAttachmentTemplate)
	if err != nil {
		return nil, err
	}

	var currentVolAttachment *provider.VolumeAttachmentResponse
	err = vpcs.APIRetry.FlexyRetry(vpcs.Logger, func() (error, bool) {
		currentVolAttachment, err = vpcs.GetVolumeAttachment(volumeAttachmentTemplate)
		if err != nil {
			// Need to stop retry as there is an error while getting attachment
			// considering that vpcs.GetVolumeAttachment already re-tried
			return err, true
		}
		// Stop retry in case of volume is attached
		return err, currentVolAttachment != nil && currentVolAttachment.Status == StatusAttached
	})
	// Success case
	if err == nil && (currentVolAttachment != nil && currentVolAttachment.Status == StatusAttached) {
		return currentVolAttachment, nil
	}
	/*
		}, func(intf interface{}, err *models.Error) bool {
			// Skip API retry logic, if there is any error keep retry as per configuration
			if err != nil {
				return skipRetry(err)
			}
			// return true in case of volume in attached status else false for retry
			return intf.(*provider.VolumeAttachmentResponse).Status == StatusAttached
		})*/
	userErr := userError.GetUserError(string(userError.VolumeAttachTimedOut), nil, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID)
	vpcs.Logger.Info("Wait for attach timed out", zap.Error(userErr))

	return nil, userErr
}
