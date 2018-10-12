/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package softlayer_block

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/common"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
)

const (
	// Softlayer storage provider
	SoftLayer            = provider.VolumeProvider("SOFTLAYER-BLOCK")
	SoftLayerEndurance   = provider.VolumeProviderType("SOFTLAYER_ENDURANCE")
	SoftLayerPerformance = provider.VolumeProviderType("SOFTLAYER_PERFORMANCE")

	VolumeTypeBlock = provider.VolumeType("VOLUME-Block")
)

// SLBlockSession implements lib.Session
type SLBlockSession struct {
	common.SLSession
}

// Close at present does nothing
func (*SLBlockSession) Close() {
	// Do nothing for now
}

// GetProviderDisplayName returns the name of the SoftLayer provider
// DEPRECATED
func (sls *SLBlockSession) GetProviderDisplayName() provider.VolumeProvider {
	return SoftLayer
}

func (sls *SLBlockSession) ProviderName() provider.VolumeProvider {
	return SoftLayer
}
