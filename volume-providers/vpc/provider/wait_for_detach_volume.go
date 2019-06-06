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
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"go.uber.org/zap"
	"time"
)

// WaitForDetachVolume waits for volume to be detached from node. e.g waits till no volume attachment is found
func (vpcs *VPCSession) WaitForDetachVolume(volumeAttachmentTemplate provider.VolumeAttachmentRequest) error {
	vpcs.Logger.Debug("Entry of WaitForDetachVolume method...")
	defer vpcs.Logger.Debug("Exit from WaitForDetachVolume method...")
	var err error
	vpcs.Logger.Info("Validating basic inputs for WaitForDetachVolume method...", zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate))
	err = vpcs.validateAttachVolumeRequest(volumeAttachmentTemplate)
	if err != nil {
		return err
	}
	maxTimeout, maxRetryAttempt, retryGapDuration := vpcs.Config.GetTimeOutParameters()
	retryCount := 0
	vpcs.Logger.Info("Waiting for volume to be detached", zap.Int("maxTimeout", maxTimeout))
	for retryCount < maxRetryAttempt {
		currentVolAttachment, errAPI := vpcs.GetVolumeAttachment(volumeAttachmentTemplate)
		if errAPI != nil {
			if errMsg, ok := errAPI.(util.Message); ok {
				if errMsg.Code == userError.VolumeAttachFindFailed {
					// Consider volume detachment is complete if  error code is VolumeAttachFindFailed
					vpcs.Logger.Info("Volume detachment is complete", zap.Int("retry attempt", retryCount), zap.Int("max retry attepmts", maxRetryAttempt))
					return nil
				}
				// do not retry if there is another error
				vpcs.Logger.Error("Error occured while finding volume attachment", zap.Error(errAPI))
				userErr := userError.GetUserError(string(userError.VolumeDetachFailed), errAPI, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID)
				return userErr
			}
		}
		// retry when volume attachment is still there
		retryCount = retryCount + 1
		vpcs.Logger.Info("Volume is still detaching. Retry..", zap.Int("retry attempt", retryCount), zap.Int("max retry attepmts", maxRetryAttempt), zap.Reflect("currentVolAttachment", currentVolAttachment))
		time.Sleep(retryGapDuration)
	}
	userErr := userError.GetUserError(string(userError.VolumeDetachTimedOut), err, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID, vpcs.Config.Timeout)
	vpcs.Logger.Info("Wait for detach timed out", zap.Error(userErr))
	return userErr
}
