/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package main

import (
	"fmt"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"go.uber.org/zap"
)

var volumeReq provider.Volume

//VolumeManager ...
type VolumeManager struct {
	Session   provider.Session
	Logger    *zap.Logger
	RequestID string
}

// NewVolumeManager ...
func NewVolumeManager(session provider.Session, logger *zap.Logger, requestID string) *VolumeManager {
	return &VolumeManager{
		Session:   session,
		Logger:    logger,
		RequestID: requestID,
	}
}

//UpdateVolume ...
func (vam *VolumeManager) UpdateVolume() {
	vam.setupVolumeRequest()
	err := vam.Session.UpdateVolume(volumeReq)
	if err != nil {
		updateRequestID(err, vam.RequestID)
		vam.Logger.Error("Failed to attach the volume", zap.Error(err))
		return
	}
	fmt.Println("Volume update", err)
}

func (vam *VolumeManager) setupVolumeRequest() {
	fmt.Printf("Enter the volume id: ")
	_, _ = fmt.Scanf("%s", &volumeID)
	fmt.Printf("Enter the provider: ")
	_, _ = fmt.Scanf("%s", &instanceID)
	fmt.Printf("Enter the cluster id: ")
	_, _ = fmt.Scanf("%s", &clusterID)
	capacity := 30
	iops := "10"
	volumeReq = provider.Volume{
		VolumeID: volumeID,
		Capacity: &capacity,
		Iops:     &iops,
		Provider: "ibm.io.vpc.block",
	}
	volumeReq.Tags = []string{"clusterid:" + clusterID, "reclaimpolicy:Delete"}

}
