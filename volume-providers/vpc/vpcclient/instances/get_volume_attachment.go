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
func (vs *VolumeAttachService) GetVolumeAttachment(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) (*models.VolumeAttachment, error) {
	defer util.TimeTracker("DetachVolume", time.Now())

	operation := &client.Operation{
		Name:        "GetVolumeAttachment",
		Method:      "GET",
		PathPattern: vs.pathPrefix + instanceIDattachmentIDPath,
	}

	var apiErr models.Error
	var volumeAttachment models.VolumeAttachment
	request := vs.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command  details", zap.Reflect("URL", request.URL()), zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate), zap.Reflect("Operation", operation))
	ctxLogger.Info("Pathparameters", zap.Reflect(instanceIDParam, volumeAttachmentTemplate.InstanceID), zap.Reflect(attachmentIDParam, volumeAttachmentTemplate.ID))
	req := request.PathParameter(instanceIDParam, *volumeAttachmentTemplate.InstanceID)
	req = request.PathParameter(attachmentIDParam, volumeAttachmentTemplate.ID)
	if volumeAttachmentTemplate.ClusterID != nil {
		// IKS case - requires ClusterID in  the request
		req = req.AddQueryValue("clusterID", *volumeAttachmentTemplate.ClusterID)
	}
	_, err := req.JSONSuccess(&volumeAttachment).JSONError(&apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while getting volume attahment", zap.Error(err))
		return nil, err
	}
	ctxLogger.Info("Successfuly retrieved the volume attachment", zap.Reflect("volumeAttachment", volumeAttachment))
	return &volumeAttachment, err
}
