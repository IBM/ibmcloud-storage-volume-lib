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

// GetSnapshot GETs from /volumes
func (ss *SnapshotService) GetSnapshot(volumeID string, snapshotID string, ctxLogger *zap.Logger) (*models.Snapshot, error) {
	ctxLogger.Info("Entry Backend GetSnapshot")
	defer ctxLogger.Info("Exit Backend GetSnapshot")

	defer util.TimeTracker("GetShanpshot", time.Now())

	operation := &client.Operation{
		Name:        "GetSnapshot",
		Method:      "GET",
		PathPattern: snapshotIDPath,
	}

	var snapshot models.Snapshot
	var apiErr models.Error

	request := ss.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command", zap.Reflect("URL", request.URL()))

	req := request.PathParameter(volumeIDParam, volumeID).PathParameter(snapshotIDParam, snapshotID)

	_, err := req.JSONSuccess(&snapshot).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}
