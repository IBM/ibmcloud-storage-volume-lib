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

// CreateSnapshot POSTs to /volumes
func (vs *VolumeService) CreateSnapshot(volumeID string, snapshotTemplate *models.Snapshot) (*models.Snapshot, error) {
	operation := &client.Operation{
		Name:        "CreateSnapshot",
		Method:      "POST",
		PathPattern: snapshotsPath,
	}

	var snapshot models.Snapshot
	var apiErr models.Error

	_, err := vs.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).JSONBody(snapshotTemplate).JSONSuccess(&snapshot).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}
