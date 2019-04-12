/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpcvolume

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
)

// VolumeManager operations
type VolumeManager interface {
	// Create the volume with authorisation by passing required information in the volume object
	CreateVolume(*models.Volume, *zap.Logger) (*models.Volume, error)

	// Delete the volume
	DeleteVolume(volumeID string) error

	// Get the volume by using ID
	GetVolume(volumeID string) (*models.Volume, error)

	// Others
	// Get volume lists by using snapshot tags
	ListVolumes(limit int, filters *models.ListVolumeFilters) (*models.VolumeList, error)

	// Set tag for a volume
	SetVolumeTag(volumeID string, tagName string) error

	// Delete tag of a volume
	DeleteVolumeTag(volumeID string, tagName string) error

	// List all tags of a volume
	ListVolumeTags(volumeID string) (*[]string, error)

	// Check if the given tag exists on a volume
	CheckVolumeTag(volumeID string, tagName string) error
}

// VolumeService ...
type VolumeService struct {
	client client.SessionClient
}

var _ VolumeManager = &VolumeService{}

// New ...
func New(client client.SessionClient) VolumeManager {
	return &VolumeService{
		client: client,
	}
}
