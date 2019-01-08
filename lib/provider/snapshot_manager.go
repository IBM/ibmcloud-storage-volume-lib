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

type SnapshotManager interface {
	// Create snapshot space
	SnapshotOrder(VolumeRequest Volume) error

	// Snapshot operations
	// Create the snapshot on the volume
	SnapshotCreate(volume *Volume, tags map[string]string) (*Snapshot, error)

	// Delete the snapshot
	SnapshotDelete(*Snapshot) error

	// Get the snapshot
	SnapshotGet(snapshotID string) (*Snapshot, error)

	// Snapshot list by using tags
	SnapshotsList() ([]*Snapshot, error)

	//List all the  snapshots for a given volume
	ListAllSnapshots(volumeID string) ([]*Snapshot, error)
}
