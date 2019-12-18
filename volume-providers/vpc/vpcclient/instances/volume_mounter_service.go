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

// VolumeAttachManager operations
//go:generate counterfeiter -o fakes/volume_attach_service.go --fake-name VolumeAttachService . VolumeAttachManager
type VolumeAttachManager interface {
	// Create the volume with authorisation by passing required information in the volume object
	AttachVolume(*models.VolumeAttachment, *zap.Logger) (*models.VolumeAttachment, error)
	// GetVolumeAttachment retrives the single VolumeAttachment based on the instance ID and attachmentID
	GetVolumeAttachment(*models.VolumeAttachment, *zap.Logger) (*models.VolumeAttachment, error)
	// ListVolumeAttachments retrives the VolumeAttachment list for given server
	ListVolumeAttachments(*models.VolumeAttachment, *zap.Logger) (*models.VolumeAttachmentList, error)
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

// IKSVolumeAttachService ...
type IKSVolumeAttachService struct {
	client                  client.SessionClient
	pathPrefix              string
	receiverError           error
	populateQueryParameters func(request *client.Request, volumeAttachmentTemplate *models.VolumeAttachment) *client.Request
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

var _ VolumeAttachManager = &IKSVolumeAttachService{}

// NewIKSVolumeAttachmentManager ...
func NewIKSVolumeAttachmentManager(clientIn client.SessionClient) VolumeAttachManager {
	err := models.IksError{}
	return &IKSVolumeAttachService{
		client:        clientIn,
		pathPrefix:    IksPathPrefix,
		receiverError: &err,
	}
}
