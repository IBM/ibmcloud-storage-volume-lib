/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"net/http"
	"time"
)

const (
	// SUCCESS ...
	SUCCESS = "Success"
	// FAILURE ...
	FAILURE = "Failure"
	// NOTSUPPORTED ...
	NOTSUPPORTED = "Not supported"
)

// VolumeAttachManager ...
type VolumeAttachManager interface {
	//Attach method attaches a volume/ fileset to a server
	AttachVolume(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error)
	//Detach detaches the volume/ fileset from the server
	DetachVolume(detachRequest VolumeAttachmentRequest) (*http.Response, error)

	/*
	     Below method will be uncommented when there is support

	   	//Wait for the volume to be attached on the node
	   	WaitForAttach(devicePath string, opts map[string]string) VolumeResponse

	   	//Wait for the volume to be detached on the node
	   	WaitForDetach(devicePath string) VolumeResponse
	*/

	//GetAttachAttachment retirves the current status of given volume attach request
	GetVolumeAttachment(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error)
}

// VolumeAttachmentResponse used for both attach and detach operation
type VolumeAttachmentResponse struct {
	VolumeAttachmentRequest
	//Status status of the volume attachment success, failed, attached, attaching, detaching
	Status    string     `json:"status,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// VolumeAttachmentRequest  used for both attach and detach operation
type VolumeAttachmentRequest struct {
	VolumeID   string `json:"volumeID"`
	InstanceID string `json:"instanceID"`
	// Only for VPC provider
	VPCVolumeAttachment *VolumeAttachment `json:"vpcVolumeAttachment"`
}
