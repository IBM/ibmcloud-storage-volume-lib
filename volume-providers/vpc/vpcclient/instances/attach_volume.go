/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package instances

import (
	"fmt"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"time"
)

// AttachVolume attached volume to instances with givne volume attachment details
func (vs *VolumeMountService) AttachVolume(volumeAttachmentTemplate *models.VolumeAttachment) (*models.VolumeAttachment, error) {
	defer util.TimeTracker("AttachVolume", time.Now())

	operation := &client.Operation{
		Name:        "AttachVolume",
		Method:      "POST",
		PathPattern: instanceIDvolumeAttachmentPath,
	}

	var volumeAttachment models.VolumeAttachment
	var apiErr models.Error
	request := vs.client.NewRequest(operation).PathParameter(instanceIDvolumeAttachmentPath, volumeAttachmentTemplate.InstanceID).JSONBody(volumeAttachmentTemplate).JSONSuccess(&volumeAttachment).JSONError(&apiErr)
	fmt.Println("Volume attachment request", request)
	_, err := request.Invoke()
	if err != nil {
		return nil, err
	}

	return &volumeAttachment, nil
}

//DetachVolume detach volume with given volume AttachmentID
func (vs *VolumeMountService) DetachVolume(volumeAttachmentID string) error {
	return nil
}
