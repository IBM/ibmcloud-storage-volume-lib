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
	"go.uber.org/zap"
	"strconv"
)

// GetVolume Get the volume by using ID
func (vpcs *VPCSession) GetVolume(id string) (*provider.Volume, error) {
	vpcs.Logger.Info("In provider GetVolume method")

	volume, _ := vpcs.Apiclient.Volume().GetVolume(id)
	vpcs.Logger.Info("Volume details", zap.Reflect("Volume", volume))

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
	return respVolume, nil
}
