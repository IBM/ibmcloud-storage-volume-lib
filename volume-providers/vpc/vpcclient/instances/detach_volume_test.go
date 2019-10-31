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

func TestDetachVolume(t *testing.T) {
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
		verify    func(*testing.T, *http.Response, error)
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
			name:   "Verify that the volume detachment is done correctly",
			status: http.StatusOK,
			verify: func(t *testing.T, httpResponse *http.Response, err error) {
				if assert.Nil(t, err) {
					assert.Nil(t, httpResponse)
				}
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {

			template := &models.VolumeAttachment{
				ID:         "volumeattachmentid",
				Name:       "volumeattachment",
				InstanceID: &instanceID,
				Volume: &models.Volume{
					ID:       "volume-id",
					Name:     "volume-name",
					Capacity: 10,
					ResourceGroup: &models.ResourceGroup{
						ID: "rg1",
					},
					Zone: &models.Zone{Name: "test-1"},
				},
			}

			mux, client, teardown := test.SetupServer(t)
			test.SetupMuxResponse(t, mux, "/v1/instances/testinstance/volume_attachments/volumeattachmentid", http.MethodDelete, nil, testcase.status, testcase.content, nil)

			defer teardown()

			logger.Info("Test case being executed", zap.Reflect("testcase", testcase.name))

			volumeAttachService := instances.New(client)

			response, err := volumeAttachService.DetachVolume(template, logger)

			if testcase.expectErr != "" && assert.Error(t, err) {
				assert.Equal(t, testcase.expectErr, err.Error())
				assert.NotNil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
			}
		})
	}
}
