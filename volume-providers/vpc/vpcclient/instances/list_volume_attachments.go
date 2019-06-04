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
func (vs *VolumeAttachService) ListVolumeAttachment(instanceID string, ctxLogger *zap.Logger) (*models.VolumeAttachmentList, error) {
	defer util.TimeTracker("GetAttachStatus", time.Now())

	operation := &client.Operation{
		Name:        "GetAttachStatus",
		Method:      "GET",
		PathPattern: instanceIDvolumeAttachmentPath,
	}

	var volumeAttachmentList models.VolumeAttachmentList
	var apiErr models.Error

	request := vs.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command  details", zap.Reflect("URL", request.URL()), zap.Reflect("instanceID", instanceID), zap.Reflect("Operation", operation))
	ctxLogger.Info("Pathparameters", zap.Reflect(instanceIDParam, instanceID))
	req := request.PathParameter(instanceIDParam, instanceID)
	_, err := req.JSONSuccess(&volumeAttachmentList).JSONError(&apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while getting volume attahment", zap.Error(err))
		return nil, err
	}
	return &volumeAttachmentList, nil
}
