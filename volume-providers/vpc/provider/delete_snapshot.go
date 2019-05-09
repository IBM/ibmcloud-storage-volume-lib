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
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"go.uber.org/zap"
)

// DeleteSnapshot delete snapshot
func (vpcs *VPCSession) DeleteSnapshot(snapshot *provider.Snapshot) error {
	vpcs.Logger.Info("Entry DeleteSnapshot", zap.Reflect("snapshot", snapshot))
	defer vpcs.Logger.Info("Exit DeleteSnapshot", zap.Reflect("snapshot", snapshot))

	var err error
	_, err = vpcs.GetSnapshot(snapshot.SnapshotID)
	if err != nil {
		return userError.GetUserError("StorageFindFailedWithSnapshotId", err, snapshot.SnapshotID, "Not a valid snapshot ID")
	}

	err = retry(vpcs.Logger, func() error {
		err = vpcs.Apiclient.SnapshotService().DeleteSnapshot(snapshot.Volume.VolumeID, snapshot.SnapshotID, vpcs.Logger)
		return err
	})

	if err != nil {
		return userError.GetUserError("FailedToDeleteSnapshot", err)
	}

	vpcs.Logger.Info("Successfully deleted the snapshot with backend (vpcclient) call)")
	return nil
}
