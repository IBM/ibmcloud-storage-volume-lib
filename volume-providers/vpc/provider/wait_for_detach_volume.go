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
	//"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	//providerError "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"go.uber.org/zap"
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

	err = vpcs.APIRetry.FlexyRetry(vpcs.Logger, func() (error, bool) {
		currentVolAttachment, err := vpcs.GetVolumeAttachment(volumeAttachmentTemplate)
		vpcs.Logger.Info("Info", zap.Reflect("VolAtt", currentVolAttachment))
		// In case of error we should not retry as there are two conditions for error
		// 1- some issues at endpoint side --> Which is already covered in vpcs.GetVolumeAttachment
		// 2- Attachment not found --> in this case we should not re-try as it has been deleted
		if err != nil {
			return err, true
		}
		return err, false
	})

	// Could be a success case
	if err != nil {
		if errMsg, ok := err.(util.Message); ok {
			if errMsg.Code == userError.VolumeAttachFindFailed {
				vpcs.Logger.Info("Volume detachment is complete")
				return nil
			}
		}
	}
	/*
		}, func(intf interface{}, err *models.Error) bool {
			// Skip API retry logic, if there is any error keep retry as per configuration
			if err != nil {
				// stop re-try, as attchment find failed, because volume is already detached
				if err.Errors[0].Code == userError.VolumeAttachFindFailed {
					return true
				}
				// keep re-try for all other errors
				return false
			}
			return false // Keep retry until timeout
		})

		// Return nil in case of successfully volume detached after re-try
		if err != nil {
			if errMsg, ok := err.(*models.Error); ok {
				if errMsg.Errors[0].Code == userError.VolumeAttachFindFailed {
					// Consider volume detachment is complete if  error code is VolumeAttachFindFailed
					vpcs.Logger.Info("Volume detachment is complete")
					return nil
				}
			}
		}*/

	userErr := userError.GetUserError(string(userError.VolumeDetachTimedOut), err, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID)
	vpcs.Logger.Info("Wait for detach timed out", zap.Error(userErr))
	return userErr
}
