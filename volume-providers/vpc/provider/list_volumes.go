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
	"fmt"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/metrics"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"strings"
	"time"
)

// ListVolumes list all volumes
func (vpcs *VPCSession) ListVolumes(limit int, start string, tags map[string]string) (*provider.VolumeList, error) {
	vpcs.Logger.Info("Entry ListVolumes", zap.Reflect("start", start), zap.Reflect("filters", tags))
	defer vpcs.Logger.Info("Exit ListVolumes", zap.Reflect("start", start), zap.Reflect("filters", tags))
	defer metrics.UpdateDurationFromStart(vpcs.Logger, "ListVolumes", time.Now())

	if limit < 0 {
		return nil, userError.GetUserError("InvalidListVolumesLimit", fmt.Errorf(
			"listVolumes got invalid entries request %v, supports values between 0-100", limit))
	}

	if limit > 100 {
		vpcs.Logger.Warn(fmt.Sprintf("listVolumes requested max entries of %v, supports values <=100 so defaulting value back to 100", limit))
		limit = 100
	}

	filters := &models.ListVolumeFilters{
		// Tag:          tags["tag"],
		ResourceGroupID: tags["resource_group.id"],
		ZoneName:        tags["zone.name"],
		VolumeName:      tags["name"],
	}

	vpcs.Logger.Info("Getting volumes list from VPC provider...", zap.Reflect("start", start), zap.Reflect("filters", filters))

	var volumes *models.VolumeList
	var err error
	err = retry(vpcs.Logger, func() error {
		volumes, err = vpcs.Apiclient.VolumeService().ListVolumes(limit, start, filters, vpcs.Logger)
		return err
	})

	if err != nil {
		if strings.Contains(err.Error(), "start parameter is not found")  {
			return nil, userError.GetUserError("StartVolumeIDNotFound", err, start)
		}
		return nil, userError.GetUserError("ListVolumesFailed", err)
	}

	vpcs.Logger.Info("Successfully retrieved volumes list from VPC backend", zap.Reflect("VolumesList", volumes))

	var respVolumesList = &provider.VolumeList{}
	if volumes != nil {
		if volumes.Next != nil {
			var next string
			// "Next":{"href":"https://eu-gb.iaas.cloud.ibm.com/v1/volumes?start=3e898aa7-ac71-4323-952d-a8d741c65a68\u0026limit=1\u0026zone.name=eu-gb-1"}
			if strings.Contains(volumes.Next.Href, "start=") {
				next = strings.Split(strings.Split(volumes.Next.Href, "start=")[1], "\u0026")[0]
			} else {
				vpcs.Logger.Warn("Volumes.Next.Href is not in expected format", zap.Reflect("volumes.Next.Href", volumes.Next.Href))
			}
			respVolumesList.Next = next
		}

		volumeslist := volumes.Volumes
		if volumeslist != nil && len(volumeslist) > 0 {
			for _, volItem := range volumeslist {
				volumeResponse := FromProviderToLibVolume(volItem, vpcs.Logger)
				respVolumesList.Volumes = append(respVolumesList.Volumes, volumeResponse)
			}
		}
	}
	return respVolumesList, err
}
