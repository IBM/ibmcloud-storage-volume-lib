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

// CreateVolume Get the volume by using ID
func (vpcs *VPCSession) CreateVolume(volumeRequest provider.Volume) (*provider.Volume, error) {
	vpcs.Logger.Info("Entry CreateVolume", zap.Reflect("volumeRequest", volumeRequest))
	defer vpcs.Logger.Info("Exit CreateVolume", zap.Reflect("volumeRequest", volumeRequest))

	vpcs.Logger.Info("Creating volume as per order request .... ", zap.Reflect("volumeRequest", volumeRequest))

	var err error
	var volume *models.Volume

	// Volume name should not be empty
	if len(*volumeRequest.Name) == 0 {
		return nil, reasoncode.GetUserError("InvalidVolumeName", nil)
	}

	// Capacity should not be empty
	if volumeRequest.Capacity == nil || *volumeRequest.Capacity == 0 {
		return nil, reasoncode.GetUserError("VolumeCapacityInvalid", nil)
	}

	// General purpose profiles does not allow IOPs setting
	if volumeRequest.VPCVolume.Profile.Name != "general-purpose" && (volumeRequest.Iops == nil || *volumeRequest.Iops <= strconv.Itoa(0)) {
		return nil, reasoncode.GetUserError("IopsInvalid", nil)
	}

	// General purpose profiles does not allow IOPs setting
	if *volumeRequest.Iops > strconv.Itoa(0) && volumeRequest.VPCVolume.Profile.Name == "general-purpose" {
		return nil, reasoncode.GetUserError("VolumeProfileIopsInvalid", nil)
	}

	// Pending error handling
	// TODO: Check if the volume already exists with same name.
	// We can do this by scanning all volumes. But requesting the VPC team to get
	// an API for getting volume details with name instead of only ID.

	iops := ToInt64(*volumeRequest.Iops)

	// Build the template we'll send to RIAAS
	volumeTemplate := &models.Volume{
		Name:     *volumeRequest.Name,
		Capacity: int64(*volumeRequest.Capacity),
		Iops:     iops,
		Tags:     volumeRequest.VPCVolume.Tags,
		ResourceGroup: &models.ResourceGroup{
			ID: volumeRequest.VPCVolume.ResourceGroup.ID,
		},
		Generation: "gt",
		Profile: &models.Profile{
			Name: volumeRequest.VPCVolume.Profile.Name,
		},
		Zone: &models.Zone{
			Name: volumeRequest.Az,
		},
	}

	vpcs.Logger.Info("Volume request details", zap.Reflect("Volume request template", volumeTemplate))
	err = retry(func() error {
		volume, err = vpcs.Apiclient.VolumeService().CreateVolume(volumeTemplate)
		return err
	})

	if err != nil {
		return nil, reasoncode.GetUserError("FailedToPlaceOrder", err)
	}

	vpcs.Logger.Info("Created volume details", zap.Reflect("Volume", volume))

	var volumeResponse *provider.Volume
	volumeResponse, err = vpcs.GetVolume(volume.ID)

	return volumeResponse, err
}
