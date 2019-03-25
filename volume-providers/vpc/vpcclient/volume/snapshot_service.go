/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package volume

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
)

// Snapshot operations
type SnapshotManager interface {
	// Create the snapshot on the volume
	CreateSnapshot(volumeID string, snapshotTemplate *models.Snapshot) (*models.Snapshot, error)

	// Delete the snapshot
	DeleteSnapshot(volumeID string, snapshotID string) error

	// Get the snapshot
	GetSnapshot(volumeID string, snapshotID string) (*models.Snapshot, error)

	// List all the  snapshots for a given volume
	ListSnapshots(volumeID string) (*models.SnapshotList, error)

	// Set tag for a snapshot
	SetSnapshotTag(volumeID string, snapshotID string, tagName string) error

	// Delete tag of a snapshot
	DeleteSnapshotTag(volumeID string, snapshotID string, tagName string) error

	// List all tags of a snapshot
	ListSnapshotTags(volumeID string, snapshotID string) (*[]string, error)

	// Check if the given tag exists on a snapshot
	CheckSnapshotTag(volumeID string, snapshotID string, tagName string) error
}

// SnapshotService ...
type SnapshotService struct {
	client client.ClientSession
}

var _ SnapshotManager = &SnapshotService{}

// New ...
func NewSnapshotManager(client client.ClientSession) SnapshotManager {
	return &SnapshotService{
		client: client,
	}
}
