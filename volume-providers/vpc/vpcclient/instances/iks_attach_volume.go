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
func (vs *IKSVolumeAttachService) AttachVolume(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) (*models.VolumeAttachment, error) {
	defer util.TimeTracker("IKS AttachVolume", time.Now())

	operation := &client.Operation{
		Name:        "AttachVolume",
		Method:      "POST",
		PathPattern: vs.pathPrefix + "createAttachment",
	}

	var volumeAttachment models.VolumeAttachment
	apiErr := vs.receiverError

	operationRequest := vs.client.NewRequest(operation)

	operationRequest = operationRequest.SetQueryValue(IksClusterQueryKey, *volumeAttachmentTemplate.ClusterID)
	operationRequest = operationRequest.SetQueryValue(IksWorkerQueryKey, *volumeAttachmentTemplate.InstanceID)
	vol := *volumeAttachmentTemplate.Volume
	operationRequest = operationRequest.SetQueryValue(IksVolumeQueryKey, vol.ID)

	ctxLogger.Info("Equivalent curl command and query parameters", zap.Reflect("URL", operationRequest.URL()), zap.Reflect("Payload", volumeAttachmentTemplate), zap.Reflect("Operation", operation), zap.Reflect(IksClusterQueryKey, volumeAttachmentTemplate.ClusterID), zap.Reflect(IksWorkerQueryKey, volumeAttachmentTemplate.InstanceID), zap.Reflect(IksVolumeQueryKey, vol.ID))

	_, err := operationRequest.JSONBody(volumeAttachmentTemplate).JSONSuccess(&volumeAttachment).JSONError(apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	ctxLogger.Info("Successfuly attached the volume")
	return &volumeAttachment, nil
}
