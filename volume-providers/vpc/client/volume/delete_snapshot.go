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

// DeleteSnapshot DELETEs to /volumes
func (vs *VolumeService) DeleteSnapshot(volumeID string, snapshotID string) error {
	operation := &client.Operation{
		Name:        "DeleteSnapshot",
		Method:      "DELETE",
		PathPattern: snapshotIDPath,
	}

	var apiErr models.Error

	_, err := vs.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).PathParameter(snapshotIDParam, snapshotID).JSONError(&apiErr).Invoke()
	if err != nil {
		return err
	}

	return nil
}
