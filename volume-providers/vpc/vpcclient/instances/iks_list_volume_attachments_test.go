/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package instances_test

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/instances"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/riaas/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIKSListVolumeAttachment(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	instanceID := "testinstance"
	clusterID := "testcluster"
	// IKS tests
	mux, client, teardown := test.SetupServer(t)
	content := "{\"volume_attachments\":[{\"id\":\"volumeattachmentid1\", \"name\":\"volume attachment\", \"device\": {\"id\":\"xvdc\"}, \"volume\": {\"id\":\"volume-id1\",\"name\":\"volume-name\",\"capacity\":10,\"iops\":3000,\"status\":\"pending\"}, \"instance_id\":\"testinstance\"}]}"

	test.SetupMuxResponse(t, mux, "/v2/storage/getAttachmentslist", http.MethodGet, nil, http.StatusOK, content, nil)
	volumeAttachService := instances.NewIKSVolumeAttachmentManager(client)

	template := &models.VolumeAttachment{
		ID:         "volumeattachmentid",
		Name:       "volume attachment",
		ClusterID:  &clusterID,
		InstanceID: &instanceID,
		Volume: &models.Volume{
			ID:       "volume-id",
			Name:     "volume-name",
			Capacity: 10,
			ResourceGroup: &models.ResourceGroup{
				ID: "rg1",
			},
			Generation: models.GenerationType("gc"),
			Zone:       &models.Zone{Name: "test-1"},
		},
	}
	defer teardown()

	volumeAttachmentsList, err := volumeAttachService.ListVolumeAttachment(template, logger)

	assert.NoError(t, err)
	assert.NotNil(t, volumeAttachmentsList)
}
