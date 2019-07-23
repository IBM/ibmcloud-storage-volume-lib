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

	err = vpcs.APIRetry.FlexyRetry(vpcs.Logger, func() (interface{}, error) {
		currentVolAttachment, err := vpcs.GetVolumeAttachment(volumeAttachmentTemplate)
		return currentVolAttachment, err
	}, func(intf interface{}, err *models.Error) bool {
		// Skip API retry logic, if there is any error keep retry as per configuration
		if err != nil {
			// stop re-try, as attchment find failed, because volume is already detached
			if err.Errors[0].Code == userError.VolumeAttachFindFailed {
				return false
			}
			// keep re-try for all other errors
			return true
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
	}

	userErr := userError.GetUserError(string(userError.VolumeDetachTimedOut), err, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID, vpcs.Config.Timeout)
	vpcs.Logger.Info("Wait for detach timed out", zap.Error(userErr))
	return userErr
}
