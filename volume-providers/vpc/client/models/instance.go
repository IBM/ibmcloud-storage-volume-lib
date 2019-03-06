/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package models

// StatusType ...
type StatusType string

// String ...
func (i StatusType) String() string { return string(i) }

// Status values
const (
	StatusFailed     StatusType = "failed"
	StatusPaused     StatusType = "paused"
	StatusPausing    StatusType = "pausing"
	StatusPending    StatusType = "pending"
	StatusRunning    StatusType = "running"
	StatusStarting   StatusType = "starting"
	StatusStopped    StatusType = "stopped"
	StatusStopping   StatusType = "stopping"
	StatusRestarting StatusType = "restarting"
	StatusResuming   StatusType = "resuming"
)

// Instance ...
type Instance struct {
	BootVolumeAttachment    *VolumeAttachment   `json:"boot_volume_attachment,omitempty"`
	CPU                     *CPU                `json:"cpu,omitempty"`
	CRN                     string              `json:"crn,omitempty"`
	Generation              GenerationType      `json:"generation,omitempty"`
	GPU                     *GPU                `json:"gpu,omitempty"`
	Href                    string              `json:"href,omitempty"`
	ID                      string              `json:"id,omitempty"`
	Image                   *Image              `json:"image,omitempty"`
	Keys                    *[]Key              `json:"keys,omitempty"`
	Memory                  int64               `json:"memory,omitempty"`
	Name                    string              `json:"name,omitempty"`
	NetworkInterfaces       *[]NetworkInterface `json:"network_interfaces,omitempty"`
	PrimaryNetworkInterface *NetworkInterface   `json:"primary_network_interface,omitempty"`
	Profile                 *Profile            `json:"profile,omitempty"`
	ResourceGroup           *ResourceGroup      `json:"resource_group,omitempty"`
	Status                  StatusType          `json:"status,omitempty"`
	Tags                    []string            `json:"tags,omitempty"`
	VolumeAttachments       *[]VolumeAttachment `json:"volume_attachments,omitempty"`
	VPC                     *VPC                `json:"vpc,omitempty"`
	Zone                    *Zone               `json:"zone,omitempty"`
}

// InstanceList ...
type InstanceList struct {
	Instances  []*Instance `json:"instances,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	TotalCount int         `json:"total_count,omitempty"`
}
