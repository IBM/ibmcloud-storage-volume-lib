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
	"time"
)

// GetVolume POSTs to /volumes
func (vs *VolumeService) GetVolume(volumeID string) (*models.Volume, error) {
	defer util.TimeTracker("GetVolume", time.Now())

	operation := &client.Operation{
		Name:        "GetVolume",
		Method:      "GET",
		PathPattern: volumeIDPath,
	}

	req := vs.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID)
	var volume models.Volume
	var apiErr models.Error

	_, err := req.JSONSuccess(&volume).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &volume, nil
}
