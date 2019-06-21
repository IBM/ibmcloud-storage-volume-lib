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
	"net/http"
	"time"
)

// DetachVolume retrives the volume attach status with givne volume attachment details
func (vs *VolumeAttachService) DetachVolume(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) (*http.Response, error) {
	defer util.TimeTracker("DetachVolume", time.Now())

	operation := &client.Operation{
		Name:        "DetachVolume",
		Method:      "DELETE",
		PathPattern: vs.pathPrefix + instanceIDattachmentIDPath,
	}

	var apiErr models.Error

	request := vs.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command  details", zap.Reflect("URL", request.URL()), zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate), zap.Reflect("Operation", operation))
	ctxLogger.Info("Pathparameters", zap.Reflect(instanceIDParam, volumeAttachmentTemplate.InstanceID), zap.Reflect(attachmentIDParam, volumeAttachmentTemplate.ID))
	req := request.PathParameter(instanceIDParam, *volumeAttachmentTemplate.InstanceID)
	req = request.PathParameter(attachmentIDParam, volumeAttachmentTemplate.ID)
	if volumeAttachmentTemplate.ClusterID != nil {
		// IKS case - requires ClusterID in  the request
		req = req.AddQueryValue("clusterID", *volumeAttachmentTemplate.ClusterID)
	}
	resp, err := req.JSONError(&apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while deleting volume attahment", zap.Error(err))
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			// volume Attachment is deleted. So do not want to retry
			return resp, apiErr
		}
	}
	ctxLogger.Info("Successfuly deleted the volume attachment")
	return resp, err
}
