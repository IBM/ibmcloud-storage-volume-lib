/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package models

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
)

// VolumeAttachment ...
type VolumeAttachment struct {
	provider.VolumeAttachment
	Volume *Volume `json:"volume,omitempty"`
}

// VolumeAttachmentList ...
type VolumeAttachmentList struct {
	VolumeAttachments []VolumeAttachment `json:"volume_attachments,omitempty"`
}
