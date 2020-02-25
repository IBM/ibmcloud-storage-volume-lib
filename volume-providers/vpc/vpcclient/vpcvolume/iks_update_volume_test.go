/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2020 All Rights Reserved.
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
	"testing"
)

func TestUpdateVolume(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	volumeTemplate := models.Volume{
		ID:         "volume-id",
		VolumeType: "block",
		Provider:   "vpc-classic",
		Cluster:    "cluster-id",
		CRN:        "crn:v1:staging:public:is:us-south-1:a/account-id::volume:volume-id",
		Tags:       []string{"tag1:val1", "tag2:val2"},
		Capacity:   2,
		Iops:       300,
	}

	testCases := []struct {
		name string

		// Response
		status        int
		volumeRequest models.Volume

		// Expected return
		expectErr string
		verify    func(*testing.T, *models.Volume, error)
	}{
		{
			name:   "Verify that the correct endpoint is invoked",
			status: http.StatusNoContent,
		}, {
			name:          "Verify that the volume is updated succesfully",
			status:        http.StatusOK,
			volumeRequest: volumeTemplate,
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			mux, client, teardown := test.SetupServer(t)
			test.SetupMuxResponse(t, mux, "/v2/storage/updateVolume", http.MethodPost, nil, http.StatusOK, "", nil)

			defer teardown()

			logger.Info("Test case being executed", zap.Reflect("testcase", testcase.name))

			volumeService := vpcvolume.NewIKSVolumeService(client)

			err := volumeService.UpdateVolume(&testcase.volumeRequest, logger)

			if testcase.expectErr != "" && assert.Error(t, err) {
				assert.Equal(t, testcase.expectErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
