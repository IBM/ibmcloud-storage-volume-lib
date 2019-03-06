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

// CheckSnapshotTag checks if the given tag exists on a snapshot
func (vs *VolumeService) CheckSnapshotTag(volumeID string, snapshotID string, tagName string) error {
	operation := &client.Operation{
		Name:        "CheckSnapshotTag",
		Method:      "GET",
		PathPattern: snapshotTagNamePath,
	}

	var apiErr models.Error

	req := vs.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).PathParameter(snapshotIDParam, snapshotID).PathParameter(snapshotTagParam, tagName).JSONError(&apiErr)
	_, err := req.Invoke()
	if err != nil {
		return err
	}

	return nil
}
