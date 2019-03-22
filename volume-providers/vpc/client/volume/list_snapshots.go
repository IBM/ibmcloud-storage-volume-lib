/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package volume

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/client/models"
)

// ListSnapshots GETs /volumes/snapshots
func (ss *SnapshotService) ListSnapshots(volumeID string) (*models.SnapshotList, error) {
	operation := &client.Operation{
		Name:        "ListSnapshots",
		Method:      "GET",
		PathPattern: snapshotsPath,
	}

	var snapshots models.SnapshotList
	var apiErr models.Error

	_, err := ss.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).JSONSuccess(&snapshots).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &snapshots, nil
}
