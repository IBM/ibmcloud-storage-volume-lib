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

type VolumeManager interface {
	// Provider name
	ProviderName() VolumeProvider

	// Type returns the underlying volume type
	Type() VolumeType

	// Volume operations
	// Create the volume with authorization by passing required information in the volume object
	VolumeCreate(VolumeRequest Volume) (*Volume, error)

	// Create the volume from snapshot with snapshot tags
	VolumeCreateFromSnapshot(snapshot Snapshot, tags map[string]string) (*Volume, error)

	// Delete the volume
	VolumeDelete(*Volume) error

	// Get the volume by using ID  //
	VolumeGet(id string) (*Volume, error)

	// Others
	// Get volume lists by using snapshot tags
	VolumesList(tags map[string]string) ([]*Volume, error)

	// GetVolumeByRequestID fetch the volume by request ID.
	// Request Id is the one that is returned when volume is provsioning request is
	// placed with Iaas provider.
	GetVolumeByRequestID(requestID string) (*Volume, error)
}
