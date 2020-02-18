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

// UpdateVolume POSTs to /volumes
func (vs *IKSVolumeService) UpdateVolume(volumeTemplate *models.Volume, ctxLogger *zap.Logger) error {
	ctxLogger.Debug("Entry Backend IKSVolumeService.UpdateVolume")
	defer ctxLogger.Debug("Exit Backend IKSVolumeService.UpdateVolume")

	defer util.TimeTracker("IKSVolumeService.UpdateVolume", time.Now())

	operation := &client.Operation{
		Name:        "UpdateVolume",
		Method:      "PUT",
		PathPattern: vs.pathPrefix + updateVolume,
	}
	apiErr := vs.receiverError
	request := vs.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command", zap.Reflect("URL", request.URL()), zap.Reflect("Operation", operation), zap.Reflect("volumeTemplate", volumeTemplate))

	_, err := request.JSONBody(volumeTemplate).JSONError(apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Update volume failed with error", zap.Error(err), zap.Error(apiErr))
	}
	return err
}
