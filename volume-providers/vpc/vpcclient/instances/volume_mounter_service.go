/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package instances

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"net/http"
)

const (
	//VpcPathPrefix  VPC URL path prefix
	VpcPathPrefix = "v1/instances"
	//IksPathPrefix  IKS URL path prefix
	IksPathPrefix = "v2/storage/clusters/{cluster-id}/workers"
)

// VolumeAttachManager operations
//go:generate counterfeiter -o fakes/volume_attach_service.go --fake-name VolumeAttachService . VolumeAttachManager
type VolumeAttachManager interface {
	// Create the volume with authorisation by passing required information in the volume object
	AttachVolume(*models.VolumeAttachment, *zap.Logger) (*models.VolumeAttachment, error)
	// GetVolumeAttachment retrives the single VolumeAttachment based on the instance ID and attachmentID
	GetVolumeAttachment(*models.VolumeAttachment, *zap.Logger) (*models.VolumeAttachment, error)
	// ListVolumeAttachment retrives the VolumeAttachment list for given server
	ListVolumeAttachment(*models.VolumeAttachment, *zap.Logger) (*models.VolumeAttachmentList, error)
	// Delete the volume
	DetachVolume(*models.VolumeAttachment, *zap.Logger) (*http.Response, error)
}

// VolumeAttachService ...
type VolumeAttachService struct {
	client                       client.SessionClient
	pathPrefix                   string
	receiverError                error
	populatePathPrefixParameters func(request *client.Request, volumeAttachmentTemplate *models.VolumeAttachment) *client.Request
}

var _ VolumeAttachManager = &VolumeAttachService{}

// New ...
func New(clientIn client.SessionClient) VolumeAttachManager {
	err := models.Error{}
	return &VolumeAttachService{
		client:        clientIn,
		pathPrefix:    VpcPathPrefix,
		receiverError: &err,
		populatePathPrefixParameters: func(request *client.Request, volumeAttachmentTemplate *models.VolumeAttachment) *client.Request {
			request.PathParameter(instanceIDParam, *volumeAttachmentTemplate.InstanceID)
			return request
		},
	}
}

// IKSVolumeAttachService ...
type IKSVolumeAttachService struct {
	VolumeAttachService
}

var _ VolumeAttachManager = &IKSVolumeAttachService{}

// NewIKSVolumeAttachmentManager ...
func NewIKSVolumeAttachmentManager(clientIn client.SessionClient) VolumeAttachManager {
	err := models.IksError{}
	return &IKSVolumeAttachService{
		VolumeAttachService{
			client:        clientIn,
			pathPrefix:    IksPathPrefix,
			receiverError: &err,
			populatePathPrefixParameters: func(request *client.Request, volumeAttachmentTemplate *models.VolumeAttachment) *client.Request {
				request.PathParameter(instanceIDParam, *volumeAttachmentTemplate.InstanceID)
				request.PathParameter(clusterIDParam, *volumeAttachmentTemplate.ClusterID)
				return request
			},
		},
	}
}
