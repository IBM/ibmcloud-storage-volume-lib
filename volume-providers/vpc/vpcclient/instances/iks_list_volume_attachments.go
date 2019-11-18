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

// ListVolumeAttachments retrives the list volume attachments with givne volume attachment details
func (vs *IKSVolumeAttachService) ListVolumeAttachments(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) (*models.VolumeAttachmentList, error) {
	defer util.TimeTracker("IKS ListVolumeAttachments", time.Now())

	operation := &client.Operation{
		Name:        "ListVolumeAttachment",
		Method:      "GET",
		PathPattern: vs.pathPrefix + "getAttachmentsList",
	}

	var volumeAttachmentList models.VolumeAttachmentList

	apiErr := vs.receiverError
	vs.client = vs.client.WithQueryValue(IksClusterQueryKey, *volumeAttachmentTemplate.ClusterID)
	vs.client = vs.client.WithQueryValue(IksWorkerQueryKey, *volumeAttachmentTemplate.InstanceID)

	operationRequest := vs.client.NewRequest(operation)

	ctxLogger.Info("Equivalent curl command and query parameters", zap.Reflect("URL", operationRequest.URL()), zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate), zap.Reflect("Operation", operation), zap.Reflect(IksClusterQueryKey, *volumeAttachmentTemplate.ClusterID), zap.Reflect(IksWorkerQueryKey, *volumeAttachmentTemplate.InstanceID))

	_, err := operationRequest.JSONSuccess(&volumeAttachmentList).JSONError(apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while getting volume attachments list", zap.Error(err))
		return nil, err
	}
	ctxLogger.Info("Successfuly retrieved the volume attachments")
	return &volumeAttachmentList, nil
}
