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
	//Its non bloking call and does not wait to complete the attachment
	AttachVolume(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error)
	//Detach detaches the volume/ fileset from the server
	//Its non bloking call and does not wait to complete the detachment
	DetachVolume(detachRequest VolumeAttachmentRequest) (*http.Response, error)

	//WaitForAttachVolume waits for the volume to be attached to the host
	//Return error if wait is timed out OR there is other error
	WaitForAttachVolume(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error)

	//WaitForDetachVolume waits for the volume to be detached from the host
	//Return error if wait is timed out OR there is other error
	WaitForDetachVolume(detachRequest VolumeAttachmentRequest) error

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
	// Only for SL provider
	SoftlayerOptions map[string]string `json:"softlayerOptions,omitempty"`
	// Only for VPC provider
	VPCVolumeAttachment *VolumeAttachment `json:"vpcVolumeAttachment"`
	// Only IKS provider
	IKSVolumeAttachment *IKSVolumeAttachment `json:"iksVolumeAttachment"`
}
