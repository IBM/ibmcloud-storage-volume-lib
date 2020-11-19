/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpcvolume

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"time"
)

// ExpandVolume PATCH to /volumes
func (vs *VolumeService) ExpandVolume(volumeID string, volumeTemplate *models.Volume, ctxLogger *zap.Logger) (*models.Volume, error) {
	ctxLogger.Debug("Entry Backend ExpandVolume")
	defer ctxLogger.Debug("Exit Backend ExpandVolume")

	defer util.TimeTracker("ExpandVolume", time.Now())

	operation := &client.Operation{
		Name:        "ExpandVolume",
		Method:      "PATCH",
		PathPattern: volumeIDPath,
	}

	var volume models.Volume
	var apiErr models.Error

	request := vs.client.NewRequest(operation)
	req := request.PathParameter(volumeIDParam, volumeID)
	ctxLogger.Info("Equivalent curl command and payload details", zap.Reflect("URL", req.URL()), zap.Reflect("Payload", volumeTemplate), zap.Reflect("Operation", operation))
	_, err := req.JSONBody(volumeTemplate).JSONSuccess(&volume).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &volume, nil
}
