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
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
)

// DeleteVolume POSTs to /volumes
func (vs *VolumeService) DeleteVolume(volumeID string) error {
	operation := &client.Operation{
		Name:        "DeleteVolume",
		Method:      "DELETE",
		PathPattern: volumeIDPath,
	}

	var apiErr models.Error

	_, err := vs.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).JSONError(&apiErr).Invoke()
	if err != nil {
		return err
	}

	return nil
}