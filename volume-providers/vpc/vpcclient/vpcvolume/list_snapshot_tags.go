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
	"time"
)

// ListSnapshotTags GETs /volumes/snapshots/tags
func (ss *SnapshotService) ListSnapshotTags(volumeID string, snapshotID string) (*[]string, error) {
	defer util.TimeTracker("ListSnapshotTags", time.Now())

	operation := &client.Operation{
		Name:        "ListSnapshotTags",
		Method:      "GET",
		PathPattern: snapshotTagsPath,
	}

	var tags []string
	var apiErr models.Error

	req := ss.client.NewRequest(operation).PathParameter(volumeIDParam, volumeID).PathParameter(snapshotIDParam, snapshotID).JSONSuccess(&tags).JSONError(&apiErr)
	_, err := req.Invoke()
	if err != nil {
		return nil, err
	}

	return &tags, nil
}
