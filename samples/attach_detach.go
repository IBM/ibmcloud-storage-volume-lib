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

//VolumeAttachmentManager ...
type VolumeAttachmentManager struct {
	Session   provider.Session
	Logger    *zap.Logger
	RequestID string
}

// NewVolumeAttachmentManager ...
func NewVolumeAttachmentManager(session provider.Session, logger *zap.Logger, requestID string) *VolumeAttachmentManager {
	return &VolumeAttachmentManager{
		Session:   session,
		Logger:    logger,
		RequestID: requestID,
	}
}

var instanceID string

//var volumeID string
//var clusterID string
var volumeAttachmentReq provider.VolumeAttachmentRequest

//AttachVolume ...
func (vam *VolumeAttachmentManager) AttachVolume() {
	vam.setupVolumeAttachmentRequest()
	response, err := vam.Session.AttachVolume(volumeAttachmentReq)
	if err != nil {
		updateRequestID(err, vam.RequestID)
		vam.Logger.Error("Failed to attach the volume", zap.Error(err))
		return
	}
	volumeAttachmentReq.VPCVolumeAttachment = &provider.VolumeAttachment{
		ID: response.VPCVolumeAttachment.ID,
	}
	response, err = vam.Session.WaitForAttachVolume(volumeAttachmentReq)
	if err != nil {
		updateRequestID(err, vam.RequestID)
		vam.Logger.Error("Failed to complete volume attach", zap.Error(err))
	}
	fmt.Println("Volume attachment", response, err)
}

//DetachVolume ...
func (vam *VolumeAttachmentManager) DetachVolume() {
	vam.setupVolumeAttachmentRequest()
	response, err := vam.Session.DetachVolume(volumeAttachmentReq)
	if err != nil {
		updateRequestID(err, vam.RequestID)
		vam.Logger.Error("Failed to detach the volume", zap.Error(err))
		return
	}
	err = vam.Session.WaitForDetachVolume(volumeAttachmentReq)
	if err != nil {
		updateRequestID(err, vam.RequestID)
		vam.Logger.Error("Failed to complete volume detach", zap.Error(err))
	}
	fmt.Println("Volume attachment", response, err)
}

func (vam *VolumeAttachmentManager) setupVolumeAttachmentRequest() {
	fmt.Printf("Enter the volume id: ")
	_, _ = fmt.Scanf("%s", &volumeID)
	fmt.Printf("Enter the instance id: ")
	_, _ = fmt.Scanf("%s", &instanceID)
	fmt.Printf("Enter the cluster id: ")
	_, _ = fmt.Scanf("%s", &clusterID)
	volumeAttachmentReq = provider.VolumeAttachmentRequest{
		VolumeID:   *volumeID,
		InstanceID: instanceID,
		VPCVolumeAttachment: &provider.VolumeAttachment{
			DeleteVolumeOnInstanceDelete: false,
		},
		IKSVolumeAttachment: &provider.IKSVolumeAttachment{
			ClusterID: clusterID,
		},
	}

}
