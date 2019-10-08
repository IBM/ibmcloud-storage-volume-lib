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
	defer util.TimeTracker("GetVolumeAttachment", time.Now())

	operation := &client.Operation{
		Name:        "GetVolumeAttachment",
		Method:      "GET",
		PathPattern: vs.pathPrefix + instanceIDattachmentIDPath,
	}

	operation.PathPattern = vs.pathPrefix + "getAttachment"

	apiErr := vs.receiverError
	var volumeAttachment models.VolumeAttachment
	request := vs.client.NewRequest(operation)

	ctxLogger.Info("Equivalent curl command  details and query parameters", zap.Reflect(IksClusterQuery, *volumeAttachmentTemplate.ClusterID), zap.Reflect(clusterIDParam, *volumeAttachmentTemplate.InstanceID), zap.Reflect(IksVolumeAttachmentIDQuery, volumeAttachmentTemplate.ID))
	request = request.AddQueryValue(IksClusterQuery, *volumeAttachmentTemplate.ClusterID)
	request = request.AddQueryValue(clusterIDParam, *volumeAttachmentTemplate.InstanceID)
	request = request.AddQueryValue(IksVolumeAttachmentIDQuery, volumeAttachmentTemplate.ID)

	_, err := request.JSONSuccess(&volumeAttachment).JSONError(apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while getting volume attachment", zap.Error(err))
		return nil, err
	}
	ctxLogger.Info("Successfuly retrieved the volume attachment", zap.Reflect("volumeAttachment", volumeAttachment))
	return &volumeAttachment, err
}
