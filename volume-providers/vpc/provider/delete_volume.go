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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"go.uber.org/zap"
)

// DeleteVolume deletes the volume
func (vpcs *VPCSession) DeleteVolume(vol *provider.Volume) error {
	vpcs.Logger.Info("Entry DeleteVolume()", zap.Reflect("vol", vol))
	var err error
	err = retry(func() error {
		err = vpcs.Apiclient.Volume().DeleteVolume(vol.VolumeID)
		return err
	})

	if err != nil {
		vpcs.Logger.Error("Error occured while deleting the volume", zap.Error(err))
		return reasoncode.GetUserError("FailedToDeleteVolume", err)
	}
	vpcs.Logger.Info("Exit DeleteVolume()")
	return err
}
