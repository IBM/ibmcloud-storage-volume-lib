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

import "time"

// Volume ...
type Volume struct {
	Href              string              `json:"href,omitempty"`
	ID                string              `json:"id,omitempty"`
	Name              string              `json:"name,omitempty"`
	Capacity          int64               `json:"capacity,omitempty"`
	Iops              int64               `json:"iops,omitempty"`
	ResourceGroup     *ResourceGroup      `json:"resource_group,omitempty"`
	Tags              []string            `json:"tags,omitempty"`
	Generation        GenerationType      `json:"generation,omitempty"`
	Profile           *Profile            `json:"profile,omitempty"`
	Snapshot          *Snapshot           `json:"snapshot,omitempty"`
	CreatedAt         *time.Time          `json:"created_at,omitempty"`
	Status            StatusType          `json:"status,omitempty"`
	VolumeAttachments *[]VolumeAttachment `json:"volume_attachments,omitempty"`

	Zone *Zone  `json:"zone,omitempty"`
	CRN  string `json:"crn,omitempty"`
}

// ListVolumeFilters ...
type ListVolumeFilters struct {
	ResourceGroupID string
	Tag             string
	ZoneName        string
}

// VolumeList ...
type VolumeList struct {
	Volumes    []*Volume `json:"volumes,omitempty"`
	Limit      int       `json:"limit,omitempty"`
	TotalCount int       `json:"total_count,omitempty"`
}
