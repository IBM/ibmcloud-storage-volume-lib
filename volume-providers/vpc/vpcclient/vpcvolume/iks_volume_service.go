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
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
)

const (
	//IksV2PathPrefix ...
	IksV2PathPrefix = "v2/storage/"
)

// IKSVolumeService ...
type IKSVolumeService struct {
	VolumeService
	pathPrefix    string
	receiverError error
}

var _ VolumeManager = &IKSVolumeService{}

// NewIKSVolumeService ...
func NewIKSVolumeService(client client.SessionClient) VolumeManager {
	err := models.IksError{}
	iksVolumeService := &IKSVolumeService{
		VolumeService: VolumeService{
			client: client,
		},
		pathPrefix:    IksV2PathPrefix,
		receiverError: &err,
	}
	return iksVolumeService
}
