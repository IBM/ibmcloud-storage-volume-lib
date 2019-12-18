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
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"

	"go.uber.org/zap"
	"net/http"
	"time"
)

// DetachVolume detach volume based on given volume attachment request
func (vpcs *VPCSession) DetachVolume(volumeAttachmentTemplate provider.VolumeAttachmentRequest) (*http.Response, error) {
	vpcs.Logger.Debug("Entry of DetachVolume method...")
	defer vpcs.Logger.Debug("Exit from DetachVolume method...")
	defer metrics.UpdateDurationFromStart(vpcs.Logger, "DetachVolume", time.Now())
	var err error
	vpcs.Logger.Info("Validating basic inputs for detach method...", zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate))
	err = vpcs.validateAttachVolumeRequest(volumeAttachmentTemplate)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	var volumeAttachment models.VolumeAttachment

	err = vpcs.APIRetry.FlexyRetry(vpcs.Logger, func() (error, bool) {
		// First , check if volume is already attached to given instance
		vpcs.Logger.Info("Checking if volume is already attached ")
		currentVolAttachment, err := vpcs.GetVolumeAttachment(volumeAttachmentTemplate)
		if err == nil && currentVolAttachment.Status != StatusDetaching {
			// If no error and current volume is not already in detaching state ( i.e in attached or attaching state) attemp to detach
			vpcs.Logger.Info("Found volume attachment", zap.Reflect("currentVolAttachment", currentVolAttachment))
			volumeAttachment := models.NewVolumeAttachment(volumeAttachmentTemplate)
			volumeAttachment.ID = currentVolAttachment.VPCVolumeAttachment.ID
			vpcs.Logger.Info("Detaching volume from VPC provider...")
			response, err = vpcs.APIClientVolAttachMgr.DetachVolume(&volumeAttachment, vpcs.Logger)
			return err, err == nil // Retry in case of all errors
		}
		vpcs.Logger.Info("No volume attachment found for", zap.Reflect("currentVolAttachment", currentVolAttachment), zap.Error(err))
		// consider volume detach success if its  already  in Detaching or VolumeAttachment is not found
		response = &http.Response{
			StatusCode: http.StatusOK,
		}
		return nil, true // skip retry if volume is not found OR alreadd in detaching state
	})
	if err != nil {
		userErr := userError.GetUserError(string(userError.VolumeDetachFailed), err, volumeAttachmentTemplate.VolumeID, volumeAttachmentTemplate.InstanceID, volumeAttachment.ID)
		vpcs.Logger.Error("Volume detach failed with error", zap.Error(err))
		return response, userErr
	}
	vpcs.Logger.Info("Successfully detached volume from VPC provider", zap.Reflect("resp", response))
	return response, nil
}
