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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
)

// GetVolume Get the volume by using ID
func (vpcs *VPCSession) CreateSnapshot(volume *provider.Volume, tags map[string]string) (*provider.Snapshot, error) {
	// Step 1: validate input
	return nil, nil
}