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
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestListVolumeAttachment(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	instanceID := "testinstance"

	testCases := []struct {
		name string

		// Response
		status  int
		content string

		// Expected return
		expectErr string
		verify    func(*testing.T, *models.VolumeAttachmentList, error)
	}{
		{
			name:   "Verify that the correct endpoint is invoked",
			status: http.StatusNoContent,
		}, {
			name:      "Verify that a 404 is returned to the caller",
			status:    http.StatusNotFound,
			content:   "{\"errors\":[{\"message\":\"testerr\"}]}",
			expectErr: "Trace Code:, testerr Please check ",
		}, {
			name:    "Verify that the volume attachment is done correctly",
			status:  http.StatusOK,
			content: "{\"id\":\"volume attachment id\", \"name\":\"volume attachment\",\"volume\": {\"id\":\"volume-id\",\"name\":\"volume-name\",\"capacity\":10,\"iops\":3000,\"status\":\"pending\"}}",
			verify: func(t *testing.T, volumeAttachmentList *models.VolumeAttachmentList, err error) {
				assert.NotNil(t, volumeAttachmentList)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			mux, client, teardown := test.SetupServer(t)
			test.SetupMuxResponse(t, mux, "/instances/testinstance/volume_attachments", http.MethodGet, nil, testcase.status, testcase.content, nil)

			defer teardown()

			logger.Info("Test case being executed", zap.Reflect("testcase", testcase.name))

			volumeAttachService := instances.New(client)

			volumeAttachments, err := volumeAttachService.ListVolumeAttachment(instanceID, logger)

			if testcase.expectErr != "" && assert.Error(t, err) {
				assert.Equal(t, testcase.expectErr, err.Error())
				assert.Nil(t, volumeAttachments)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, volumeAttachments)
			}
		})
	}
}
