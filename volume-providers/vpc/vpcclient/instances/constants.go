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

const (
	instancesPath                    = "instances"
	instanceIDParam                  = "instance-id"
	volumeIDParam                    = "volume-id"
	instanceIDPath                   = instancesPath + "/{" + instanceIDParam + "}"
	volumeAttachmentPath             = "volume_attachments"
	instanceIDvolumeAttachmentPath   = instanceIDPath + "/" + volumeAttachmentPath
	instanceIDvolumeIDAttachmentPath = instanceIDvolumeAttachmentPath + "/{" + volumeIDParam + "}"
)
