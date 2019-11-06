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

// GetVolumeAttachment retrives the volume attach status with given volume attachment details
func (vs *IKSVolumeAttachService) GetVolumeAttachment(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) (*models.VolumeAttachment, error) {
	defer util.TimeTracker("IKS GetVolumeAttachment", time.Now())

	operation := &client.Operation{
		Name:        "GetVolumeAttachment",
		Method:      "GET",
		PathPattern: vs.pathPrefix + "getAttachment",
	}

	apiErr := vs.receiverError
	var volumeAttachment models.VolumeAttachment

	operationRequest := vs.client.NewRequest(operation)
	operationRequest = operationRequest.SetQueryValue(IksClusterQueryKey, *volumeAttachmentTemplate.ClusterID)
	operationRequest = operationRequest.SetQueryValue(IksWorkerQueryKey, *volumeAttachmentTemplate.InstanceID)
	operationRequest = operationRequest.SetQueryValue(IksVolumeAttachmentIDQueryKey, volumeAttachmentTemplate.ID)

	ctxLogger.Info("Equivalent curl command and query parameters", zap.Reflect("URL", operationRequest.URL()), zap.Reflect(IksClusterQueryKey, *volumeAttachmentTemplate.ClusterID), zap.Reflect(IksWorkerQueryKey, *volumeAttachmentTemplate.InstanceID), zap.Reflect(IksVolumeAttachmentIDQueryKey, volumeAttachmentTemplate.ID))

	_, err := operationRequest.JSONSuccess(&volumeAttachment).JSONError(apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while getting volume attachment", zap.Error(err))
		return nil, err
	}
	ctxLogger.Info("Successfuly retrieved the volume attachment", zap.Reflect("volumeAttachment", volumeAttachment))
	return &volumeAttachment, err
}
