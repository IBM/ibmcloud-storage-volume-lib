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

// SetVolumeTag sets tag for a volume
func (vs *VolumeService) SetVolumeTag(volumeID string, tagName string) error {
	defer util.TimeTracker("SetVolumeTag", time.Now())

	operation := &client.Operation{
		Name:        "SetVolumeTag",
		Method:      "PUT",
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
