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

	userErr := userError.GetUserError(string(userError.VolumeDetachTimedOut), err, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID)
	vpcs.Logger.Info("Wait for detach timed out", zap.Error(userErr))
	return userErr
}
