/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package messages

// ReasonCode ...
type ReasonCode string

// Attach / Detach problems
const (
	//VolumeAttachFailed indicates if volume attach to instance is failed
	VolumeAttachFailed = ReasonCode("VolumeAttachFailed")
	//VolumeDetachFailed indicates if volume detach from instance is failed
	VolumeDetachFailed = ReasonCode("VolumeDetachFailed")
	//VolumeAttachFindFailed indicates if the volume attachment is not found with given request
	VolumeAttachFindFailed = ReasonCode("VolumeAttachFindFailed")
)
