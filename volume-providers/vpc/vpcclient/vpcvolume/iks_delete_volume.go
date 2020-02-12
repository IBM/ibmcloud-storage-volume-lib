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

// DeleteVolume POSTs to /volumes
func (vs *IKSVolumeService) DeleteVolume(volumeID string, ctxLogger *zap.Logger) error {
	ctxLogger.Debug("Entry Backend DeleteVolume")
	defer ctxLogger.Debug("Exit Backend DeleteVolume")

	defer util.TimeTracker("DeleteVolume", time.Now())

	operation := &client.Operation{
		Name:        "DeleteVolume",
		Method:      "DELETE",
		PathPattern: volumeIDPath,
	}

	var apiErr models.Error

	request := vs.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command", zap.Reflect("URL", request.URL()), zap.Reflect("Operation", operation))

	_, err := request.PathParameter(volumeIDParam, volumeID).JSONError(&apiErr).Invoke()
	if err != nil {
		return err
	}

	return nil
}
