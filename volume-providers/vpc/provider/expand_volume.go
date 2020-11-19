/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2020 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/metrics"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"time"
)

// ExpandVolume Get the volume by using ID
func (vpcs *VPCSession) ExpandVolume(expandVolumeRequest provider.ExpandVolumeRequest) (size int64, err error) {
	vpcs.Logger.Debug("Entry of ExpandVolume method...")
	defer vpcs.Logger.Debug("Exit from ExpandVolume method...")
	defer metrics.UpdateDurationFromStart(vpcs.Logger, "ExpandVolume", time.Now())

	// Get volume details
	existVolume, err := vpcs.GetVolume(expandVolumeRequest.VolumeID)
	if err != nil {
		return -1, err
	}
	// Return existing Capacity if its greater or equal to expandable size
	if existVolume.Capacity != nil && int64(*existVolume.Capacity) >= expandVolumeRequest.Capacity {
		return int64(*existVolume.Capacity), nil
	}
	vpcs.Logger.Info("Successfully validated inputs for ExpandVolume request... ")

	// Build the template to send to backend
	volumeTemplate := &models.Volume{
		Capacity: expandVolumeRequest.Capacity,
	}

	vpcs.Logger.Info("Calling VPC provider for volume expand...")
	var volume *models.Volume
	err = retry(vpcs.Logger, func() error {
		volume, err = vpcs.Apiclient.VolumeService().ExpandVolume(expandVolumeRequest.VolumeID, volumeTemplate, vpcs.Logger)
		return err
	})

	if err != nil {
		vpcs.Logger.Debug("Failed to expand volume from VPC provider", zap.Reflect("BackendError", err))
		return -1, userError.GetUserError("FailedToPlaceOrder", err)
	}

	vpcs.Logger.Info("Successfully accepted volume expansion request, now waiting for volume state equal to available")
	err = WaitForValidVolumeState(vpcs, volume.ID)
	if err != nil {
		return -1, userError.GetUserError("VolumeNotInValidState", err, volume.ID)
	}

	vpcs.Logger.Info("Volume got valid (available) state", zap.Reflect("VolumeDetails", volume))
	return expandVolumeRequest.Capacity, nil
}
