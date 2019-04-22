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

// CreateSnapshot POSTs to /volumes
func (ss *SnapshotService) CreateSnapshot(volumeID string, snapshotTemplate *models.Snapshot, ctxLogger *zap.Logger) (*models.Snapshot, error) {
	ctxLogger.Debug("Entry Backend CreateSpanShot")
	defer ctxLogger.Debug("Exit Backend CreateSnapshot")

	defer util.TimeTracker("CreateSnapshot", time.Now())

	operation := &client.Operation{
		Name:        "CreateSnapshot",
		Method:      "POST",
		PathPattern: snapshotsPath,
	}

	var snapshot models.Snapshot
	var apiErr models.Error

	request := ss.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command and payload details", zap.Reflect("URL", request.URL()), zap.Reflect("Payload", snapshotTemplate), zap.Reflect("Operation", operation))

	_, err := request.PathParameter(volumeIDParam, volumeID).JSONBody(snapshotTemplate).JSONSuccess(&snapshot).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}
