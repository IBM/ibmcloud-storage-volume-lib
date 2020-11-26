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

// VolumeManager ...
type VolumeManager interface {
	// Provider name
	ProviderName() VolumeProvider

	// Type returns the underlying volume type
	Type() VolumeType

	// Volume operations
	// Create the volume with authorization by passing required information in the volume object
	CreateVolume(VolumeRequest Volume) (*Volume, error)

	// Create the volume from snapshot with snapshot tags
	CreateVolumeFromSnapshot(snapshot Snapshot, tags map[string]string) (*Volume, error)

	// UpdateVolume the volume
	UpdateVolume(Volume) error
	// Delete the volume
	DeleteVolume(*Volume) error

	// Get the volume by using ID  //
	GetVolume(id string) (*Volume, error)

	// Get the volume by using Name,
	// actually some of providers(like VPC) has the capability to provide volume
	// details by usig user provided volume name
	GetVolumeByName(name string) (*Volume, error)

	// Get volume lists by using filters
	ListVolumes(limit int, start string, tags map[string]string) (*VolumeList, error)

	// GetVolumeByRequestID fetch the volume by request ID.
	// Request Id is the one that is returned when volume is provsioning request is
	// placed with Iaas provider.
	GetVolumeByRequestID(requestID string) (*Volume, error)

	//AuthorizeVolume allows aceess to volume  based on given authorization
	AuthorizeVolume(volumeAuthorization VolumeAuthorization) error

	// Volume operations
	// Expand the volume with authorization by passing required information in the volume object
	ExpandVolume(expandVolumeRequest ExpandVolumeRequest) (int64, error)
}
