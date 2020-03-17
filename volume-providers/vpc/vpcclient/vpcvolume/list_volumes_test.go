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
	"net/url"
	"testing"
)

func TestListVolumes(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	testCases := []struct {
		name string

		// Response
		status  int
		content string

		limit   int
		start   string
		filters *models.ListVolumeFilters

		// Expected return
		expectErr string
		verify    func(*testing.T)
		muxVerify func(*testing.T, *http.Request)
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
			name:   "Verify that limit is added to the query",
			limit:  12,
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"limit": []string{"12"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name:   "Verify that start is added to the query",
			start:  "x-y-z",
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"start": []string{"x-y-z"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name: "Verify that resource_group.id is added to the query",
			filters: &models.ListVolumeFilters{
				ResourceGroupID: "rgid",
			},
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"resource_group.id": []string{"rgid"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name: "Verify that tag is added to the query",
			filters: &models.ListVolumeFilters{
				Tag: "taggg",
			},
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"tag": []string{"taggg"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name: "Verify that zone.name is added to the query",
			filters: &models.ListVolumeFilters{
				ZoneName: "zone",
			},
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"zone.name": []string{"zone"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name: "Verify that volume name is added to the query",
			filters: &models.ListVolumeFilters{
				VolumeName: "testname",
			},
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"name": []string{"testname"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			mux, client, teardown := test.SetupServer(t)
			test.SetupMuxResponse(t, mux, vpcvolume.Version+"/volumes", http.MethodGet, nil, testcase.status, testcase.content, nil)

			defer teardown()

			logger.Info("Test case being executed", zap.Reflect("testcase", testcase.name))

			volumeService := vpcvolume.New(client)

			volumes, err := volumeService.ListVolumes(testcase.limit, testcase.start, testcase.filters, logger)
			logger.Info("Volumes", zap.Reflect("volumes", volumes))

			if testcase.expectErr != "" && assert.Error(t, err) {
				assert.Equal(t, testcase.expectErr, err.Error())
				assert.Nil(t, volumes)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, volumes)
			}

			if testcase.verify != nil {
				testcase.verify(t)
			}
		})
	}
}
