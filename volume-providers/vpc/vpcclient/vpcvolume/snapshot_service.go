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

// SnapshotManager operations
type SnapshotManager interface {
	// Create the snapshot on the volume
	CreateSnapshot(volumeID string, snapshotTemplate *models.Snapshot, ctxLogger *zap.Logger) (*models.Snapshot, error)

	// Delete the snapshot
	DeleteSnapshot(volumeID string, snapshotID string, ctxLogger *zap.Logger) error

	// Get the snapshot
	GetSnapshot(volumeID string, snapshotID string, ctxLogger *zap.Logger) (*models.Snapshot, error)

	// List all the  snapshots for a given volume
	ListSnapshots(volumeID string, ctxLogger *zap.Logger) (*models.SnapshotList, error)

	// Set tag for a snapshot
	SetSnapshotTag(volumeID string, snapshotID string, tagName string, ctxLogger *zap.Logger) error

	// Delete tag of a snapshot
	DeleteSnapshotTag(volumeID string, snapshotID string, tagName string, ctxLogger *zap.Logger) error

	// List all tags of a snapshot
	ListSnapshotTags(volumeID string, snapshotID string, ctxLogger *zap.Logger) (*[]string, error)

	// Check if the given tag exists on a snapshot
	CheckSnapshotTag(volumeID string, snapshotID string, tagName string, ctxLogger *zap.Logger) error
}

// SnapshotService ...
type SnapshotService struct {
	client client.SessionClient
}

var _ SnapshotManager = &SnapshotService{}

// NewSnapshotManager ...
func NewSnapshotManager(client client.SessionClient) SnapshotManager {
	return &SnapshotService{
		client: client,
	}
}
