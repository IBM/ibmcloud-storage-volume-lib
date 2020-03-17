/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpcvolume

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// ListVolumes GETs /volumes
func (vs *VolumeService) ListVolumes(limit int, start string, filters *models.ListVolumeFilters, ctxLogger *zap.Logger) (*models.VolumeList, error) {
	ctxLogger.Debug("Entry Backend ListVolumes")
	defer ctxLogger.Debug("Exit Backend ListVolumes")

	defer util.TimeTracker("ListVolumes", time.Now())

	operation := &client.Operation{
		Name:        "ListVolumes",
		Method:      "GET",
		PathPattern: volumesPath,
	}

	var volumes models.VolumeList
	var apiErr models.Error

	request := vs.client.NewRequest(operation)
	ctxLogger.Info("Equivalent curl command", zap.Reflect("URL", request.URL()), zap.Reflect("Operation", operation))

	req := request.JSONSuccess(&volumes).JSONError(&apiErr)

	if limit > 0 {
		req.AddQueryValue("limit", strconv.Itoa(limit))
	}

	if start != "" {
		req.AddQueryValue("start", start)
	}

	if filters != nil {
		if filters.ResourceGroupID != "" {
			req.AddQueryValue("resource_group.id", filters.ResourceGroupID)
		}
		if filters.Tag != "" {
			req.AddQueryValue("tag", filters.Tag)
		}
		if filters.ZoneName != "" {
			req.AddQueryValue("zone.name", filters.ZoneName)
		}
		if filters.VolumeName != "" {
			req.AddQueryValue("name", filters.VolumeName)
		}
	}

	_, err := req.Invoke()
	if err != nil {
		return nil, err
	}

	return &volumes, nil
}
