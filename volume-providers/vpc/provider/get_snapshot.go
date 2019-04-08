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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
)

// GetSnapshot get snapshot
func (vpcs *VPCSession) GetSnapshot(snapshotID string) (*provider.Snapshot, error) {
	vpcs.Logger.Info("Entry GetSnapshot", zap.Reflect("SnapshotID", snapshotID))
	defer vpcs.Logger.Info("Exit GetSnapshot", zap.Reflect("SnapshotID", snapshotID))

	var err error
	var snapshot *models.Snapshot

	err = retry(func() error {
		snapshot, err = vpcs.Apiclient.SnapshotService().GetSnapshot("", snapshotID)
		return err
	})

	if err != nil {
		return nil, reasoncode.GetUserError("FailedToDeleteSnapshot", err)
	}

	vpcs.Logger.Info("Successfully retrieved the snapshot details", zap.Reflect("Snapshot", snapshot))

	volume, err := vpcs.GetVolume("")
	if err != nil {
		return nil, reasoncode.GetUserError("StorageFindFailedWithVolumeId", err, volume.VolumeID, "Not a valid volume ID")
	}

	respSnapshot := &provider.Snapshot{
		SnapshotID: snapshot.ID,
		Volume:     *volume,
	}

	vpcs.Logger.Info("Successfully retrieved the snapshot details", zap.Reflect("Provider snapshot", respSnapshot))
	return respSnapshot, nil
}
