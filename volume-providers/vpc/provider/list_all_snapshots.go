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

// ListAllSnapshots list all snapshots
func (vpcs *VPCSession) ListAllSnapshots(volumeID string) ([]*provider.Snapshot, error) {
	vpcs.Logger.Info("Entry ListAllSnapshots", zap.Reflect("VolumeID", volumeID))
	defer vpcs.Logger.Info("Exit ListAllSnapshots", zap.Reflect("VolumeID", volumeID))

	return nil, nil
}
