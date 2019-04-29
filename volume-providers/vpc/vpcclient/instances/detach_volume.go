/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package instances

import (
	"fmt"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"time"
)

// DetachVolume retrives the volume attach status with givne volume attachment details
func (vs *VolumeMountService) DetachVolume(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) error {
	defer util.TimeTracker("AttachVolume", time.Now())

	operation := &client.Operation{
		Name:        "GetAttachStatus",
		Method:      "GET",
		PathPattern: instanceIDvolumeAttachmentPath,
	}

	var volumeAttachmentList models.VolumeAttachmentList
	var apiErr models.Error

	request := vs.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command  details", zap.Reflect("URL", request.URL()), zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate), zap.Reflect("Operation", operation))
	ctxLogger.Info("Pathparameters", zap.Reflect(instanceIDParam, volumeAttachmentTemplate.InstanceID), zap.Reflect(volumeIDParam, volumeAttachmentTemplate.Volume.ID))
	req := request.PathParameter(instanceIDParam, volumeAttachmentTemplate.InstanceID)
	_, err := req.JSONSuccess(&volumeAttachmentList).JSONError(&apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while getting volume attahment", zap.Error(err))
		return err
	}

	for _, volumeAttachment := range volumeAttachmentList.VolumeAttachments {
		if volumeAttachment.Volume.ID == volumeAttachmentTemplate.Volume.ID {
			ctxLogger.Info("Successfully fetched volume attachment", zap.Reflect("volumeAttachment", volumeAttachment))
			return nil
		}
	}
	// Volume is not attached to instance
	// form model error so that retry won't  happen
	apiErr = models.Error{
		Errors: []models.ErrorItem{
			models.ErrorItem{
				Code:    models.ErrorCodeNotFound,
				Message: fmt.Sprintf("volume [%s]  is not attached to instance [%s]", volumeAttachmentTemplate.Volume.ID, volumeAttachmentTemplate.InstanceID),
			},
		},
	}

	return apiErr
}
