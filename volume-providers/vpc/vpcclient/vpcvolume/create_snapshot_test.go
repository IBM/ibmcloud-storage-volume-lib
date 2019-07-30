/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpcvolume_test

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/riaas/test"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/vpcvolume"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"testing"
)

func TestCreateSnapshot(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	testCases := []struct {
		name string

		// backend url
		url string

		// Response
		status  int
		content string

		// Expected return
		expectErr string
		verify    func(*testing.T, *models.Snapshot, error)
	}{
		{
			name:   "Verify that the correct endpoint is invoked",
			status: http.StatusNoContent,
			url:    "/volumes/volume-id/snapshots",
		}, {
			name:      "Verify that a 404 is returned to the caller",
			status:    http.StatusNotFound,
			url:       vpcvolume.Version + "/volumes/volume-id/snapshots",
			content:   "{\"errors\":[{\"message\":\"testerr\"}]}",
			expectErr: "Trace Code:, testerr Please check ",
		}, {
			name:    "Verify that the snapshot is parsed correctly",
			status:  http.StatusOK,
			url:     vpcvolume.Version + "/volumes/volume-id/snapshots",
			content: "{\"id\":\"snapshot1\",\"status\":\"pending\"}",
			verify: func(t *testing.T, snapshot *models.Snapshot, err error) {
				assert.NotNil(t, snapshot)
			},
		}, {
			name:    "Verify that the snapshot is parsed correctly",
			status:  http.StatusOK,
			content: "{\"id\":\"snapshot1\",\"status\":\"pending\"}",
			url:     vpcvolume.Version + "/volumes/volume-id/snapshots",
			verify: func(t *testing.T, snapshot *models.Snapshot, err error) {
				assert.NotNil(t, snapshot)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			template := &models.Snapshot{
				Name: "snapshot-name",
				ID:   "snapshot-id",
			}
			mux, client, teardown := test.SetupServer(t)
			requestBody := `{
        			"id":"snapshot-id",
  			        "name":"snapshot-name"
      			}`
			requestBody = strings.Join(strings.Fields(requestBody), "") + "\n"
			test.SetupMuxResponse(t, mux, testcase.url, http.MethodPost, &requestBody, testcase.status, testcase.content, nil)

			defer teardown()

			logger.Info("Test case being executed", zap.Reflect("testcase", testcase.name))

			snapshotService := vpcvolume.NewSnapshotManager(client)

			snapshot, err := snapshotService.CreateSnapshot("volume-id", template, logger)
			logger.Info("Snapshot", zap.Reflect("snapshot", snapshot))

			// vpc snapshot functionality is not yet ready. It would return error for now
			if testcase.verify != nil {
				testcase.verify(t, snapshot, err)
			}
		})
	}
}
