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
	providerutils "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/util"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"time"
)

// CreateVolume POSTs to /volumes
func (vs *VolumeService) CreateVolume(volumeTemplate *models.Volume) (*models.Volume, error) {
	defer providerutils.TimeTracker("CreateVolume", time.Now())

	operation := &client.Operation{
		Name:        "CreateVolume",
		Method:      "POST",
		PathPattern: volumesPath,
	}

	var volume models.Volume
	var apiErr models.Error

	_, err := vs.client.NewRequest(operation).JSONBody(volumeTemplate).JSONSuccess(&volume).JSONError(&apiErr).Invoke()
	if err != nil {
		return nil, err
	}

	return &volume, nil
}
