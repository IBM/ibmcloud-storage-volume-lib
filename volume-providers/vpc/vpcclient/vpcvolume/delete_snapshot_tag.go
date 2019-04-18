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

// DeleteSnapshotTag deletes tag of a snapshot
func (ss *SnapshotService) DeleteSnapshotTag(volumeID string, snapshotID string, tagName string, ctxLogger *zap.Logger) error {
	ctxLogger.Info("Entry Backend DeleteSnapshotTag")
	defer ctxLogger.Info("Exit Backend DeleteSnapshotTag")

	defer util.TimeTracker("DeleteSnapshotTag", time.Now())

	operation := &client.Operation{
		Name:        "DeleteSnapshotTag",
		Method:      "DELETE",
		PathPattern: snapshotTagNamePath,
	}

	var apiErr models.Error

	request := ss.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command", zap.Reflect("URL", request.URL()))

	req := request.PathParameter(volumeIDParam, volumeID).PathParameter(snapshotIDParam, snapshotID).PathParameter(snapshotTagParam, tagName).JSONError(&apiErr)
	_, err := req.Invoke()
	if err != nil {
		return err
	}

	return nil
}
