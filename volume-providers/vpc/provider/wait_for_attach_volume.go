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
	//"time"
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
	fobj := NewFlexyRetryDefault()
	err = fobj.FlexyRetry(vpcs.Logger, func() (interface{}, error) {
		currentVolAttachment, errAPI := vpcs.GetVolumeAttachment(volumeAttachmentTemplate)
		return currentVolAttachment, errAPI
	}, func(intf interface{}, err *models.Error) bool {
		if err != nil {
			return skipRetry(err)
		}

		if intf.(*provider.VolumeAttachmentResponse).Status == StatusAttached {
			return false
		}
		return true
	})
	/*maxTimeout, maxRetryAttempt, retryGapDuration := vpcs.Config.GetTimeOutParameters()
	retryCount := 0
	vpcs.Logger.Info("Waiting for volume to be attached", zap.Int("maxTimeout", maxTimeout))
	for retryCount < maxRetryAttempt {
		currentVolAttachment, errAPI := vpcs.GetVolumeAttachment(volumeAttachmentTemplate)
		if errAPI == nil && currentVolAttachment.Status == StatusAttached {
			// volume is attached return no error
			vpcs.Logger.Info("Volume attachment is complete", zap.Int("retry attempt", retryCount), zap.Int("max retry attepmts", maxRetryAttempt), zap.Reflect("currentVolAttachment", currentVolAttachment))
			return currentVolAttachment, nil
		} else if errAPI != nil {
			// do not retry if there is error
			vpcs.Logger.Error("Error occured while finding volume attachment", zap.Error(errAPI))
			userErr := userError.GetUserError(string(userError.VolumeAttachFailed), errAPI, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID)
			return nil, userErr
		}
		// retry if attach status is not "attached"
		retryCount = retryCount + 1
		vpcs.Logger.Info("Volume is still attaching. Retry..", zap.Int("retry attempt", retryCount), zap.Int("max retry attepmts", maxRetryAttempt), zap.Reflect("currentVolAttachment", currentVolAttachment))
		time.Sleep(retryGapDuration)
	}*/
	userErr := userError.GetUserError(string(userError.VolumeAttachTimedOut), nil, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID, vpcs.Config.Timeout)
	vpcs.Logger.Info("Wait for attach timed out", zap.Error(userErr))

	return nil, userErr
}
