/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpcvolume

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
)

// ListVolumeTags GETs /volumes/tags
func (vs *VolumeService) ListVolumeTags(volumeID string) (*[]string, error) {
	operation := &client.Operation{
		Name:        "ListVolumeTags",
		Method:      "GET",
		PathPattern: volumeTagsPath,
	}

	var tags []string
	var apiErr models.Error

	req := vs.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).JSONSuccess(&tags).JSONError(&apiErr)
	_, err := req.Invoke()
	if err != nil {
		return nil, err
	}

	return &tags, nil
}
