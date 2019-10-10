/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	vpcprovider "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/provider"
	"net/http"
)

// IksVpcSession implements lib.Session for VPC IKS dual session
type IksVpcSession struct {
	vpcprovider.VPCSession                         // Holds VPC/Riaas session by default
	IksSession             *vpcprovider.VPCSession // Holds IKS session
}

var _ provider.Session = &IksVpcSession{}

const (
	// Provider storage provider
	Provider = provider.VolumeProvider("IKS-VPC-Block")
	// VolumeType ...
	VolumeType = provider.VolumeType("VPC-Block")
)

// Close at present does nothing
func (vpcIks *IksVpcSession) Close() {
	// Do nothing for now
}

// GetProviderDisplayName returns the name of the VPC provider
func (vpcIks *IksVpcSession) GetProviderDisplayName() provider.VolumeProvider {
	return Provider
}

// ProviderName ...
func (vpcIks *IksVpcSession) ProviderName() provider.VolumeProvider {
	return Provider
}

// Type ...
func (vpcIks *IksVpcSession) Type() provider.VolumeType {
	return VolumeType
}

// AttachVolume attach volume based on given volume attachment request
func (vpcIks *IksVpcSession) AttachVolume(volumeAttachmentRequest provider.VolumeAttachmentRequest) (*provider.VolumeAttachmentResponse, error) {
	vpcIks.Logger.Debug("Entry of IksVpcSession.AttachVolume method...")
	defer vpcIks.Logger.Debug("Exit from IksVpcSession.AttachVolume method...")
	return vpcIks.IksSession.AttachVolume(volumeAttachmentRequest)
}

// DetachVolume attach volume based on given volume attachment request
func (vpcIks *IksVpcSession) DetachVolume(volumeAttachmentRequest provider.VolumeAttachmentRequest) (*http.Response, error) {
	vpcIks.IksSession.Logger.Debug("Entry of IksVpcSession.DetachVolume method...")
	defer vpcIks.Logger.Debug("Exit from IksVpcSession.DetachVolume method...")
	return vpcIks.IksSession.DetachVolume(volumeAttachmentRequest)
}

// GetVolumeAttachment attach volume based on given volume attachment request
func (vpcIks *IksVpcSession) GetVolumeAttachment(volumeAttachmentRequest provider.VolumeAttachmentRequest) (*provider.VolumeAttachmentResponse, error) {
	vpcIks.Logger.Debug("Entry of IksVpcSession.GetVolumeAttachment method...")
	defer vpcIks.Logger.Debug("Exit from IksVpcSession.GetVolumeAttachment method...")
	return vpcIks.IksSession.GetVolumeAttachment(volumeAttachmentRequest)
}

// WaitForAttachVolume attach volume based on given volume attachment request
func (vpcIks *IksVpcSession) WaitForAttachVolume(volumeAttachmentRequest provider.VolumeAttachmentRequest) (*provider.VolumeAttachmentResponse, error) {
	vpcIks.Logger.Debug("Entry of IksVpcSession.WaitForAttachVolume method...")
	defer vpcIks.Logger.Debug("Exit from IksVpcSession.WaitForAttachVolume method...")
	return vpcIks.IksSession.WaitForAttachVolume(volumeAttachmentRequest)
}

// WaitForDetachVolume attach volume based on given volume attachment request
func (vpcIks *IksVpcSession) WaitForDetachVolume(volumeAttachmentRequest provider.VolumeAttachmentRequest) error {
	vpcIks.Logger.Debug("Entry of IksVpcSession.WaitForDetachVolume method...")
	defer vpcIks.Logger.Debug("Exit from IksVpcSession.WaitForDetachVolume method...")
	return vpcIks.IksSession.WaitForDetachVolume(volumeAttachmentRequest)
}
