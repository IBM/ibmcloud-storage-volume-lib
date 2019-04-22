/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import "time"

// VPCVolume specific	parameters
type VPCVolume struct {
	Href              string              `json:"href,omitempty"`
	ResourceGroup     *ResourceGroup      `json:"resource_group,omitempty"`
	Generation        GenerationType      `json:"generation,omitempty"`
	Profile           *Profile            `json:"profile,omitempty"`
	Tags              []string            `json:"tags,omitempty"`
	VolumeAttachments *[]VolumeAttachment `json:"volume_attachments,omitempty"`
	CRN               string              `json:"crn,omitempty"`
}

// GenerationType ...
type GenerationType string

// String ...
func (i GenerationType) String() string { return string(i) }

// ResourceGroup ...
type ResourceGroup struct {
	Href string `json:"href,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Profile ...
type Profile struct {
	Href string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
	CRN  string `json:"crn,omitempty"`
}

// VolumeAttachment ...
type VolumeAttachment struct {
	Href       string `json:"href,omitempty"`
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`

	Volume    *Volume    `json:"volume,omitempty"`
	Status    string     `json:"status,omitempty"` //attaching, attached, detaching
	Type      string     `json:"type,omitempty"`   //boot, data
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// If set to true, when deleting the instance the volume will also be deleted
	DeleteVolumeOnInstanceDelete bool `json:"delete_volume_on_instance_delete,omitempty"`
}
