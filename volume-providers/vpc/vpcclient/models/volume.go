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
	"strconv"
	"strings"
	"time"
)

const (
	//ClusterIDTagName ...
	ClusterIDTagName = "clusterid"
	//VolumeStatus ...
	VolumeStatus = "status"
)

// Volume ...
type Volume struct {
	Href                string               `json:"href,omitempty"`
	ID                  string               `json:"id,omitempty"`
	Name                string               `json:"name,omitempty"`
	Capacity            int64                `json:"capacity,omitempty"`
	Iops                int64                `json:"iops,omitempty"`
	VolumeEncryptionKey *VolumeEncryptionKey `json:"encryption_key,omitempty"`
	ResourceGroup       *ResourceGroup       `json:"resource_group,omitempty"`
	Tags                []string             `json:"tags,omitempty"`
	Profile             *Profile             `json:"profile,omitempty"`
	Snapshot            *Snapshot            `json:"snapshot,omitempty"`
	CreatedAt           *time.Time           `json:"created_at,omitempty"`
	Status              StatusType           `json:"status,omitempty"`
	VolumeAttachments   *[]VolumeAttachment  `json:"volume_attachments,omitempty"`

	Zone       *Zone  `json:"zone,omitempty"`
	CRN        string `json:"crn,omitempty"`
	Cluster    string `json:"cluster,omitempty"`
	Provider   string `json:"provider,omitempty"`
	VolumeType string `json:"volume_type,omitempty"`
}

// ListVolumeFilters ...
type ListVolumeFilters struct {
	ResourceGroupID string
	Tag             string
	ZoneName        string
	VolumeName      string
}

// VolumeList ...
type VolumeList struct {
	Volumes    []*Volume `json:"volumes,omitempty"`
	Limit      int       `json:"limit,omitempty"`
	TotalCount int       `json:"total_count,omitempty"`
}

//NewVolume created model volume from provider volume
func NewVolume(volumeRequest provider.Volume) Volume {
	// Build the template to send to backend

	volume := Volume{
		ID:       volumeRequest.VolumeID,
		CRN:      volumeRequest.CRN,
		Capacity: int64(*volumeRequest.Capacity),
		Tags:     volumeRequest.VPCVolume.Tags,
		Profile: &Profile{
			Name: volumeRequest.VPCVolume.Profile.Name,
		},
		Zone: &Zone{
			Name: volumeRequest.Az,
		},
		Provider:   string(volumeRequest.Provider),
		VolumeType: string(volumeRequest.VolumeType),
	}
	if volumeRequest.Name != nil {
		volume.Name = *volumeRequest.Name
	}
	if volumeRequest.VPCVolume.ResourceGroup != nil {
		volume.ResourceGroup = &ResourceGroup{
			ID:   volumeRequest.VPCVolume.ResourceGroup.ID,
			Name: volumeRequest.VPCVolume.ResourceGroup.Name,
		}
	}

	if volumeRequest.Iops != nil {
		value, err := strconv.ParseInt(*volumeRequest.Iops, 10, 64)
		if err != nil {
			volume.Iops = 0
		}
		volume.Iops = value
	}
	if volumeRequest.VPCVolume.VolumeEncryptionKey != nil && len(volumeRequest.VPCVolume.VolumeEncryptionKey.CRN) > 0 {
		encryptionKeyCRN := volumeRequest.VPCVolume.VolumeEncryptionKey.CRN
		volume.VolumeEncryptionKey = &VolumeEncryptionKey{CRN: encryptionKeyCRN}
	}

	//volume.initCluster()
	volume.Cluster = volumeRequest.Attributes[ClusterIDTagName]
	volume.Status = StatusType(volumeRequest.Attributes[VolumeStatus])
	return volume
}

func (vol Volume) initCluster() {
	for _, tag := range vol.Tags {
		if strings.Contains(tag, ClusterIDTagName) {
			clusterID := tag[len(ClusterIDTagName)+1:]
			vol.Cluster = clusterID
			break
		}
	}
}
