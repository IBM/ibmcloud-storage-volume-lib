/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/riaas/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestGetVolume(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	testCases := []struct {
		name string

		// Response
		status  int
		content string

		// Expected return
		expectErr string
		verify    func(*testing.T, *provider.Volume, error)
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
			name:    "Verify that the volume is parsed correctly",
			status:  http.StatusOK,
			content: "{\"id\":\"vol1\",\"name\":\"vol1\",\"capacity\":10,\"iops\":3000,\"status\":\"pending\",\"zone\":{\"name\":\"test-1\",\"href\":\"https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/test-1\"},\"crn\":\"crn:v1:bluemix:public:is:test-1:a/rg1::volume:vol1\"}",
			verify: func(t *testing.T, volume *provider.Volume, err error) {
				if assert.NotNil(t, volume) {
					assert.Equal(t, "vol1", volume.VolumeID)
				}
			},
		}, {
			name:    "False positive: What if the volume ID is not matched",
			status:  http.StatusOK,
			content: "{\"id\":\"wrong-vol\",\"name\":\"wrong-vol\",\"capacity\":10,\"iops\":3000,\"status\":\"pending\",\"zone\":{\"name\":\"test-1\",\"href\":\"https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/test-1\"},\"crn\":\"crn:v1:bluemix:public:is:test-1:a/rg1::volume:wrong-vol\"}",
			verify: func(t *testing.T, volume *provider.Volume, err error) {
				if assert.NotNil(t, volume) {
					assert.NotEqual(t, "vol1", volume.VolumeID)
				}
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			mux, client, teardown := test.SetupServer(t)
			emptyString := ""
			test.SetupMuxResponse(t, mux, "/volumes/volume-id", http.MethodGet, &emptyString, testcase.status, testcase.content, nil)

			assert.NotNil(t, client)
			defer teardown()

			vpcs, err := GetTestOpenSession(t, client, logger)
			assert.NotNil(t, vpcs)
			assert.Nil(t, err)

			volume, err := vpcs.GetVolume("volume-id")
			logger.Info("Volume details", zap.Reflect("volume", volume))

			if testcase.expectErr != "" && assert.Error(t, err) {
				assert.Equal(t, testcase.expectErr, err.Error())
				assert.Nil(t, volume)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, volume)
			}

			if testcase.verify != nil {
				testcase.verify(t, volume, err)
			}

		})
	}
}
