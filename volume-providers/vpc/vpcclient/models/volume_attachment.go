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
	"strings"
	"time"
)

// Device ...
type Device struct {
	ID string `json:"id"`
}

// VolumeAttachment for riaas client
type VolumeAttachment struct {
	ID   string `json:"id"`
	Href string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
	// Status of volume attachment named - attaching , attached, detaching
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"` //boot, data
	// InstanceID this volume is attached to
	InstanceID *string    `json:"instance_id,omitempty"`
	ClusterID  *string    `json:"clusterID,omitempty"`
	Device     *Device    `json:"device,omitempty"`
	Volume     *Volume    `json:"volume,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	// If set to true, when deleting the instance the volume will also be deleted
	DeleteVolumeOnInstanceDelete bool `json:"delete_volume_on_instance_delete,omitempty"`
}

// VolumeAttachmentList ...
type VolumeAttachmentList struct {
	VolumeAttachments []VolumeAttachment `json:"volume_attachments,omitempty"`
}

// NewVolumeAttachment creates VolumeAttachment from VolumeAttachmentRequest
func NewVolumeAttachment(volumeAttachmentRequest provider.VolumeAttachmentRequest) VolumeAttachment {
	va := VolumeAttachment{
		InstanceID: &volumeAttachmentRequest.InstanceID,
		Volume: &Volume{
			ID: volumeAttachmentRequest.VolumeID,
		},
	}
	if volumeAttachmentRequest.VPCVolumeAttachment != nil {
		va.ID = volumeAttachmentRequest.VPCVolumeAttachment.ID
		va.Href = volumeAttachmentRequest.VPCVolumeAttachment.Href
		va.Name = volumeAttachmentRequest.VPCVolumeAttachment.Name
		va.DeleteVolumeOnInstanceDelete = volumeAttachmentRequest.VPCVolumeAttachment.DeleteVolumeOnInstanceDelete

	}
	if volumeAttachmentRequest.IKSVolumeAttachment != nil {
		va.ClusterID = volumeAttachmentRequest.IKSVolumeAttachment.ClusterID
	}
	return va
}

//ToVolumeAttachmentResponse converts VolumeAttachment VolumeAttachmentResponse
func (va *VolumeAttachment) ToVolumeAttachmentResponse() *provider.VolumeAttachmentResponse {
	varp := &provider.VolumeAttachmentResponse{
		VolumeAttachmentRequest: provider.VolumeAttachmentRequest{
			VolumeID: va.Volume.ID,
			VPCVolumeAttachment: &provider.VolumeAttachment{
				DeleteVolumeOnInstanceDelete: va.DeleteVolumeOnInstanceDelete,
				ID:                           va.ID,
				Href:                         va.Href,
				Name:                         va.Name,
				Type:                         va.Type,
			},
		},
		Status:    va.Status,
		CreatedAt: va.CreatedAt,
	}
	if va.InstanceID != nil {
		varp.InstanceID = *va.InstanceID
	}

	//Set DevicePath
	if va.Device != nil {
		devicepath := va.Device.ID
		generation := "gc" //default
		if va.Volume != nil && va.Volume.Generation != "" {
			generation = va.Volume.Generation.String()
		}

		//prepend "/dev/" for generation=1 (gc)
		if generation == "gc" && !strings.HasPrefix(devicepath, "/dev/") {
			devicepath = "/dev/" + va.Device.ID
		}
		varp.VolumeAttachmentRequest.VPCVolumeAttachment.Device = &provider.Device{ID: devicepath}
	}
	return varp
}
