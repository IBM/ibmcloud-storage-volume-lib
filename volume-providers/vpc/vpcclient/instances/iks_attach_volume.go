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
		Name:   "AttachVolume",
		Method: "POST",
	}
	operation.PathPattern = vs.pathPrefix + "createAttachment"

	var volumeAttachment models.VolumeAttachment
	apiErr := vs.receiverError

	request := vs.client.NewRequest(operation)

	request = request.SetQueryValue(IksClusterQuery, *volumeAttachmentTemplate.ClusterID)
	request = request.SetQueryValue(IksWorkerQuery, *volumeAttachmentTemplate.InstanceID)
	vol := *volumeAttachmentTemplate.Volume
	request = request.SetQueryValue(IksVolumeQuery, vol.ID)
	ctxLogger.Info("Equivalent curl command  details and query parameters", zap.Reflect("URL", request.URL()), zap.Reflect("Payload", volumeAttachmentTemplate), zap.Reflect("Operation", operation), zap.Reflect(IksClusterQuery, volumeAttachmentTemplate.InstanceID), zap.Reflect(IksWorkerQuery, volumeAttachmentTemplate.InstanceID), zap.Reflect(IksVolumeQuery, vol.ID))
	_, err := request.JSONBody(volumeAttachmentTemplate).JSONSuccess(&volumeAttachment).JSONError(apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	ctxLogger.Info("Successfuly attached the volume")

	return &volumeAttachment, nil
}
