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

// GetVolume POSTs to /volumes
func (vs *VolumeService) GetVolume(volumeID string, ctxLogger *zap.Logger) (*models.Volume, error) {
	ctxLogger.Debug("Entry Backend GetVolume")
	defer ctxLogger.Debug("Exit Backend GetVolume")

	defer util.TimeTracker("GetVolume", time.Now())

	operation := &client.Operation{
		Name:        "GetVolume",
		Method:      "GET",
		PathPattern: volumeIDPath,
	}

	var volume models.Volume
	var apiErr models.Error

	request := vs.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command", zap.Reflect("URL", request.URL()), zap.Reflect("Operation", operation))

	req := request.PathParameter(volumeIDParam, volumeID)
	_, err := req.JSONSuccess(&volume).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &volume, nil
}
