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

// Snapshot ...
type Snapshot struct {
	Href          string         `json:"href,omitempty"`
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	ResourceGroup *ResourceGroup `json:"resource_group,omitempty"`
	CRN           string         `json:"crn,omitempty"`
	Status        StatusType     `json:"status,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
}

// SnapshotList ...
type SnapshotList struct {
	Snapshots []*Snapshot `json:"snapshot,omitempty"`
}

type SnapshotManager interface {
	// Snapshot operations
	// Create the snapshot on the volume
	CreateSnapshot(volumeID string, snapshotTemplate *Snapshot) (*Snapshot, error)

	// Delete the snapshot
	DeleteSnapshot(volumeID string, snapshotID string) error

	// Get the snapshot
	GetSnapshot(snapshotID string) (*Snapshot, error)

	// List all the  snapshots for a given volume
	ListSnapshots(volumeID string) ([]*SnapshotList, error)

        // Set tag for a snapshot
        SetSnapshotTag(volumeID string, snapshotID string, tagName string) error

        // Delete tag of a snapshot
        DeleteSnapshotTag(volumeID string, snapshotID string, tagName string) error

        // List all tags of a snapshot
        ListSnapshotTags(volumeID string, snapshotID string) (*[]string, error)

        // Check if the given tag exists on a snapshot
        CheckSnapshotTag(volumeID string, snapshotID string, tagName string) error
}
