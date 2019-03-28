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
	"go.uber.org/zap"
)

// DeleteSnapshot delete snapshot
func (vpcs *VPCSession) DeleteSnapshot(snapshot *provider.Snapshot) error {
	vpcs.Logger.Info("Entry DeleteSnapshot()", zap.Reflect("snapshot", snapshot))
	var err error
	_, err = vpcs.GetSnapshot(snapshot.SnapshotID)
	if err != nil {
		return reasoncode.GetUserError("StorageFindFailedWithSnapshotId", err, snapshot.SnapshotID, "Not a valid snapshot ID")
	}

	err = retry(func() error {
		err = vpcs.Apiclient.SnapshotService().DeleteSnapshot("", snapshot.SnapshotID)
		return err
	})

	if err != nil {
		vpcs.Logger.Error("Error occured while deleting the snapshot", zap.Error(err))
		return reasoncode.GetUserError("FailedToDeleteSnapshot", err)
	}
	vpcs.Logger.Info("Exit DeleteSnapshot()", zap.Reflect("snapshot", snapshot))
	return nil
}
