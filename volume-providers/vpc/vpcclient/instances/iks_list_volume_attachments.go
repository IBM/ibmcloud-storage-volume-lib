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

// ListVolumeAttachment retrives the list volume attachments with givne volume attachment details
func (vs *IKSVolumeAttachService) ListVolumeAttachment(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) (*models.VolumeAttachmentList, error) {
	defer util.TimeTracker("ListVolumeAttachment", time.Now())

	operation := &client.Operation{
		Name:   "ListVolumeAttachment",
		Method: "GET",
	}

	operation.PathPattern = vs.pathPrefix + "getAttachmentsList"

	var volumeAttachmentList models.VolumeAttachmentList
	apiErr := vs.receiverError

	request := vs.client.NewRequest(operation)

	ctxLogger.Info("Equivalent curl command  details and query parameters", zap.Reflect("URL", request.URL()),zap.Reflect("request", request), zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate), zap.Reflect("Operation", operation), zap.Reflect(IksClusterQuery, *volumeAttachmentTemplate.ClusterID), zap.Reflect(IksWorkerQuery, *volumeAttachmentTemplate.InstanceID))
	request = request.AddQueryValue(IksClusterQuery, *volumeAttachmentTemplate.ClusterID)
	request = request.AddQueryValue(clusterIDParam, *volumeAttachmentTemplate.InstanceID)

	_, err := request.JSONSuccess(&volumeAttachmentList).JSONError(apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while getting volume attachments list", zap.Error(err))
		return nil, err
	}
	return &volumeAttachmentList, nil
}
