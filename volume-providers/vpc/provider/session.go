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
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/instances"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/riaas"
	"go.uber.org/zap"
)

// VPCSession implements lib.Session
type VPCSession struct {
	VPCAccountID          string
	Config                *config.VPCProviderConfig
	ContextCredentials    provider.ContextCredentials
	VolumeType            provider.VolumeType
	Provider              provider.VolumeProvider
	Apiclient             riaas.RegionalAPI
	APIClientVolAttachMgr instances.VolumeAttachManager
	APIVersion            string
	Logger                *zap.Logger
}

const (
	// VPC storage provider
	VPC = provider.VolumeProvider("VPC")
	// VolumeType ...
	VolumeType = provider.VolumeType("vpc-block")
	// SnapshotMask ...
	SnapshotMask = "id,username,capacityGb,createDate,snapshotCapacityGb,parentVolume[snapshotSizeBytes],parentVolume[snapshotCapacityGb],parentVolume[id],parentVolume[storageTierLevel],parentVolume[notes],storageType[keyName],serviceResource[datacenter[name]],billingItem[location,hourlyFlag],provisionedIops,lunId,originalVolumeName,storageTierLevel,notes"
)

var (
	// DeleteVolumeReason ...
	DeleteVolumeReason = "deleted by ibm-volume-lib on behalf of user request"
)

// Close at present does nothing
func (*VPCSession) Close() {
	// Do nothing for now
}

// GetProviderDisplayName returns the name of the VPC provider
func (vpcs *VPCSession) GetProviderDisplayName() provider.VolumeProvider {
	return VPC
}

// ProviderName ...
func (vpcs *VPCSession) ProviderName() provider.VolumeProvider {
	return VPC
}

// Type ...
func (vpcs *VPCSession) Type() provider.VolumeType {
	return VolumeType
}
