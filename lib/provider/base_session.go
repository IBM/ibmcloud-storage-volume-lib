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

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
)

// BaseSession ... Empty implementation for all the methods
type BaseSession struct {
	Logger *zap.Logger
}

// ProviderName ...
func (b *BaseSession) ProviderName() VolumeProvider {
	b.Logger.Warn("Unimplemented -ProviderName()")
	return ""
}

// Type returns the underlying volume type
func (b *BaseSession) Type() VolumeType {
	b.Logger.Warn("Unimplemented -Type()")
	return ""
}

// CreateVolume operations
// Create the volume with authorization by passing required information in the volume object
func (b *BaseSession) CreateVolume(VolumeRequest Volume) (*Volume, error) {
	b.Logger.Warn("Unimplemented -CreateVolume()")
	return nil, errors.New("Unimplemented -CreateVolume()")
}

// CreateVolumeFromSnapshot Create the volume from snapshot with snapshot tags
func (b *BaseSession) CreateVolumeFromSnapshot(snapshot Snapshot, tags map[string]string) (*Volume, error) {
	b.Logger.Warn("Unimplemented -CreateVolumeFromSnapshot()")
	return nil, errors.New("Unimplemented -CreateVolumeFromSnapshot()")
}

// DeleteVolume Delete the volume
func (b *BaseSession) DeleteVolume(*Volume) error {
	b.Logger.Warn("Unimplemented - DeleteVolume()")
	return errors.New("Unimplemented - DeleteVolume()")
}

// GetVolume Get the volume by using ID  //
func (b *BaseSession) GetVolume(id string) (*Volume, error) {
	b.Logger.Warn("Unimplemented -GetVolume()")
	return nil, errors.New("Unimplemented -GetVolume()")
}

// GetVolumeByName Get the volume by using Name,
// actually some of providers(like VPC) has the capability to provide volume
// details by usig user provided volume name
func (b *BaseSession) GetVolumeByName(name string) (*Volume, error) {
	b.Logger.Warn("Unimplemented -GetVolumeByName()")
	return nil, errors.New("Unimplemented -GetVolumeByName()")
}

// ListVolumes Get volume lists by using snapshot tags
func (b *BaseSession) ListVolumes(tags map[string]string) ([]*Volume, error) {
	b.Logger.Warn("Unimplemented -ListVolumes()")
	return nil, errors.New("Unimplemented -ListVolumes()")
}

// GetVolumeByRequestID fetch the volume by request ID.
// Request Id is the one that is returned when volume is provsioning request is
// placed with Iaas provider.
func (b *BaseSession) GetVolumeByRequestID(requestID string) (*Volume, error) {
	b.Logger.Warn("Unimplemented -GetVolumeByRequestID()")
	return nil, errors.New("Unimplemented -GetVolumeByRequestID()")
}

//AuthorizeVolume allows aceess to volume  based on given authorization
func (b *BaseSession) AuthorizeVolume(volumeAuthorization VolumeAuthorization) error {
	b.Logger.Warn("Unimplemented - AuthorizeVolume()")
	return errors.New("Unimplemented - AuthorizeVolume()")
}

// OrderSnapshot Creates snapshot space
func (b *BaseSession) OrderSnapshot(VolumeRequest Volume) error {

	b.Logger.Warn("Unimplemented - OrderSnapshot()")
	return errors.New("Unimplemented - OrderSnapshot()")
}

// CreateSnapshot Create the snapshot on the volume
func (b *BaseSession) CreateSnapshot(volume *Volume, tags map[string]string) (*Snapshot, error) {

	b.Logger.Warn("Unimplemented - CreateSnapshot()")
	return nil, errors.New("Unimplemented - CreateSnapshot()")

}

// DeleteSnapshot Delete the snapshot
func (b *BaseSession) DeleteSnapshot(*Snapshot) error {

	b.Logger.Warn("Unimplemented - DeleteSnapshot()")
	return errors.New("Unimplemented - DeleteSnapshot()")
}

// GetSnapshot Get the snapshot
func (b *BaseSession) GetSnapshot(snapshotID string) (*Snapshot, error) {

	b.Logger.Warn("Unimplemented - GetSnapshot()")
	return nil, errors.New("Unimplemented - GetSnapshot()")

}

// GetSnapshotWithVolumeID Get the snapshot with volume ID
func (b *BaseSession) GetSnapshotWithVolumeID(volumeID string, snapshotID string) (*Snapshot, error) {

	b.Logger.Warn("Unimplemented - GetSnapshotWithVolumeID()")
	return nil, errors.New("Unimplemented - GetSnapshotWithVolumeID()")

}

// ListSnapshots lists by using tags
func (b *BaseSession) ListSnapshots() ([]*Snapshot, error) {

	b.Logger.Warn("Unimplemented - ListSnapshots()")
	return nil, errors.New("Unimplemented - ListSnapshots()")
}

//ListAllSnapshots List all the  snapshots for a given volume
func (b *BaseSession) ListAllSnapshots(volumeID string) ([]*Snapshot, error) {

	b.Logger.Warn("Unimplemented - ListAllSnapshots()")
	return nil, errors.New("Unimplemented - ListAllSnapshots()")

}

//AttachVolume method attaches a volume/ fileset to a server
func (b *BaseSession) AttachVolume(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error) {

	b.Logger.Warn("Unimplemented - AttachVolume()")
	return nil, errors.New("Unimplemented - AttachVolume()")

}

//DetachVolume detaches the volume/ fileset from the server
func (b *BaseSession) DetachVolume(detachRequest VolumeAttachmentRequest) (*http.Response, error) {

	b.Logger.Warn("Unimplemented - DetachVolume()")
	return nil, errors.New("Unimplemented - DetachVolume()")

}

//GetVolumeAttachment retirves the current status of given volume attach request
func (b *BaseSession) GetVolumeAttachment(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error) {

	b.Logger.Warn("Unimplemented - GetVolumeAttachment()")
	return nil, errors.New("Unimplemented - GetVolumeAttachment()")

}
