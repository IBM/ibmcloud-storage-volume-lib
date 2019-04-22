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

// ListSnapshots GETs /volumes/snapshots
func (ss *SnapshotService) ListSnapshots(volumeID string, ctxLogger *zap.Logger) (*models.SnapshotList, error) {
	ctxLogger.Debug("Entry Backend ListSnapshots")
	defer ctxLogger.Debug("Exit Backend ListSnapshots")

	defer util.TimeTracker("ListSnapshots", time.Now())

	operation := &client.Operation{
		Name:        "ListSnapshots",
		Method:      "GET",
		PathPattern: snapshotsPath,
	}

	var snapshots models.SnapshotList
	var apiErr models.Error

	request := ss.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command", zap.Reflect("URL", request.URL()), zap.Reflect("Operation", operation))

	_, err := request.PathParameter(volumeIDParam, volumeID).JSONSuccess(&snapshots).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &snapshots, nil
}
