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

// GetSnapshot GETs from /volumes
func (vs *VolumeService) GetSnapshot(volumeID string, snapshotID string) (*models.Snapshot, error) {
	operation := &client.Operation{
		Name:        "GetSnapshot",
		Method:      "GET",
		PathPattern: snapshotIDPath,
	}

	req := vs.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).PathParameter(snapshotIDParam, snapshotID)

	var snapshot models.Snapshot
	var apiErr models.Error

	_, err := req.JSONSuccess(&snapshot).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}
