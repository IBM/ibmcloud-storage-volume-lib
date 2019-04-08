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
	"strconv"
)

// GetVolume gets the volume by using ID
func (vpcs *VPCSession) GetVolume(id string) (*provider.Volume, error) {
	vpcs.Logger.Info("Entry GetVolume", zap.Reflect("VolumeID", id))
	defer vpcs.Logger.Info("Exit GetVolume", zap.Reflect("VolumeID", id))

	var err error
	var volume *models.Volume

	err = retry(func() error {
		volume, err = vpcs.Apiclient.VolumeService().GetVolume(id)
		return err
	})

	if err != nil {
		return nil, reasoncode.GetUserError("StorageFindFailedWithVolumeId", err, id, "Not a valid volume ID")
	}

	vpcs.Logger.Info("Successfully retrieved the volume details", zap.Reflect("Volume", volume))

	volumeCap := int(volume.Capacity)
	iops := strconv.Itoa(int(volume.Iops))
	respVolume := &provider.Volume{
		VolumeID:     volume.ID,
		Provider:     VPC,
		Capacity:     &volumeCap,
		Iops:         &iops,
		VolumeType:   VolumeType,
		CreationTime: *volume.CreatedAt,
		Region:       volume.Zone.Name,
	}

	vpcs.Logger.Info("Successfully retrieved the provider volume details", zap.Reflect("Provider volume", respVolume))
	return respVolume, err
}
