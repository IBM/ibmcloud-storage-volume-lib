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

// CheckVolumeTag checks if the given tag exists on a volume
func (vs *VolumeService) CheckVolumeTag(volumeID string, tagName string) error {
	operation := &client.Operation{
		Name:        "CheckVolumeTag",
		Method:      "GET",
		PathPattern: volumeTagNamePath,
	}

	var apiErr models.Error

	req := vs.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).PathParameter(volumeTagParam, tagName).JSONError(&apiErr)
	_, err := req.Invoke()
	if err != nil {
		return err
	}

	return nil
}
