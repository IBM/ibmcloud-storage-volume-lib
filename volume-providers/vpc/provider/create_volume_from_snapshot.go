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
	"go.uber.org/zap"
)

// CreateVolumeFromSnapshot creates the volume by using ID
func (vpcs *VPCSession) CreateVolumeFromSnapshot(snapshot provider.Snapshot, tags map[string]string) (*provider.Volume, error) {
	vpcs.Logger.Info("Entry CreateVolumeFromSnapshot", zap.Reflect("Snapshot", snapshot))
	defer vpcs.Logger.Info("Exit CreateVolumeFromSnapshot", zap.Reflect("Snapshot", snapshot))

	return nil, nil
}
