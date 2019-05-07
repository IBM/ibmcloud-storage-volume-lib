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

// SnapshotManager ...
type SnapshotManager interface {
	// Create snapshot space
	OrderSnapshot(VolumeRequest Volume) error

	// Snapshot operations
	// Create the snapshot on the volume
	CreateSnapshot(volume *Volume, tags map[string]string) (*Snapshot, error)

	// Delete the snapshot
	DeleteSnapshot(*Snapshot) error

	// Get the snapshot
	GetSnapshot(snapshotID string) (*Snapshot, error)

	// Get the snapshot with volume ID
	GetSnapshotWithVolumeID(volumeID string, snapshotID string) (*Snapshot, error)

	// Snapshot list by using tags
	ListSnapshots() ([]*Snapshot, error)

	//List all the  snapshots for a given volume
	ListAllSnapshots(volumeID string) ([]*Snapshot, error)
}
