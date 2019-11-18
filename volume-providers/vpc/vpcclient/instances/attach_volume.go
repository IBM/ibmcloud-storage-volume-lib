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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"time"
)

// AttachVolume attached volume to instances with givne volume attachment details
func (vs *VolumeAttachService) AttachVolume(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) (*models.VolumeAttachment, error) {
	defer util.TimeTracker("AttachVolume", time.Now())

	operation := &client.Operation{
		Name:        "AttachVolume",
		Method:      "POST",
		PathPattern: vs.pathPrefix + instanceIDvolumeAttachmentPath,
	}

	var volumeAttachment models.VolumeAttachment
	apiErr := vs.receiverError

	operationRequest := vs.client.NewRequest(operation)

	ctxLogger.Info("Equivalent curl command and payload details", zap.Reflect("URL", operationRequest.URL()), zap.Reflect("Payload", volumeAttachmentTemplate), zap.Reflect("Operation", operation), zap.Reflect("PathParameters", volumeAttachmentTemplate.InstanceID))
	_, err := vs.populatePathPrefixParameters(operationRequest, volumeAttachmentTemplate).JSONBody(volumeAttachmentTemplate).JSONSuccess(&volumeAttachment).JSONError(apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	ctxLogger.Info("Successfuly attached the volume")
	return &volumeAttachment, nil
}
