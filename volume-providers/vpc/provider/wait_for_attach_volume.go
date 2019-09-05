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
	"go.uber.org/zap"
	"time"
)

// WaitForAttachVolume waits for volume to be attached to node. e.g waits till status becomes attached
func (vpcs *VPCSession) WaitForAttachVolume(volumeAttachmentTemplate provider.VolumeAttachmentRequest) (*provider.VolumeAttachmentResponse, error) {
	vpcs.Logger.Debug("Entry of WaitForAttachVolume method...")
	defer vpcs.Logger.Debug("Exit from WaitForAttachVolume method...")
	defer metrics.UpdateDurationFromStart(vpcs.Logger, "WaitForAttachVolume", time.Now())

	vpcs.Logger.Info("Validating basic inputs for WaitForAttachVolume method...", zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate))
	err := vpcs.validateAttachVolumeRequest(volumeAttachmentTemplate)
	if err != nil {
		return nil, err
	}

	var currentVolAttachment *provider.VolumeAttachmentResponse
	err = vpcs.APIRetry.FlexyRetryWithConstGap(vpcs.Logger, func() (error, bool) {
		currentVolAttachment, err = vpcs.GetVolumeAttachment(volumeAttachmentTemplate)
		if err != nil {
			// Need to stop retry as there is an error while getting attachment
			// considering that vpcs.GetVolumeAttachment already re-tried
			return err, true
		}
		// Stop retry in case of volume is attached
		return err, currentVolAttachment != nil && currentVolAttachment.Status == StatusAttached
	})
	// Success case, checks are required in case of timeout happened and volume is still not attached state
	if err == nil && (currentVolAttachment != nil && currentVolAttachment.Status == StatusAttached) {
		return currentVolAttachment, nil
	}

	userErr := userError.GetUserError(string(userError.VolumeAttachTimedOut), nil, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID)
	vpcs.Logger.Info("Wait for attach timed out", zap.Error(userErr))

	return nil, userErr
}
