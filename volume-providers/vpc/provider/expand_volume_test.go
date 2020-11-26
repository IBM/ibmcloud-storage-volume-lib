/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2020 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"errors"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	volumeServiceFakes "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/vpcvolume/fakes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestExpandVolume(t *testing.T) {
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		volumeService *volumeServiceFakes.VolumeService
	)

	testCases := []struct {
		testCaseName       string
		volumeID           string
		baseVolume         *models.Volume
		newSize            int64
		expectedErr        string
		expectedSize       int64
		expectedReasonCode string
	}{
		{
			testCaseName: "OK",
			volumeID:     "16f293bf-test-4bff-816f-e199c0c65db5",
			baseVolume: &models.Volume{
				ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
				Status:   models.StatusType("available"),
				Capacity: int64(10),
				Iops:     int64(1000),
				Zone:     &models.Zone{Name: "test-zone"},
			},
			newSize:      20,
			expectedSize: 20,
		},
		{
			testCaseName: "same size-success",
			volumeID:     "16f293bf-test-4bff-816f-e199c0c65db5",
			baseVolume: &models.Volume{
				ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
				Status:   models.StatusType("available"),
				Capacity: int64(10),
				Iops:     int64(1000),
				Zone:     &models.Zone{Name: "test-zone"},
			},
			newSize:      10,
			expectedSize: 10,
		},
		{
			testCaseName:       "volume not found",
			volumeID:           "16f293bf-test-4bff-816f-e199c0c65db5",
			baseVolume:         nil,
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Wrong volume ID' volume ID is not valid. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_id_invalid, BackendError:, RC:400}",
			expectedReasonCode: "ErrorUnclassified",
			newSize:            10,
			expectedSize:       -1,
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.testCaseName, func(t *testing.T) {
			vpcs, uc, sc, err := GetTestOpenSession(t, logger)
			assert.NotNil(t, vpcs)
			assert.NotNil(t, uc)
			assert.NotNil(t, sc)
			assert.Nil(t, err)

			volumeService = &volumeServiceFakes.VolumeService{}
			assert.NotNil(t, volumeService)
			uc.VolumeServiceReturns(volumeService)

			if testcase.expectedErr != "" {
				volumeService.GetVolumeReturns(testcase.baseVolume, errors.New(testcase.expectedReasonCode))
				volumeService.ExpandVolumeReturns(testcase.baseVolume, errors.New(testcase.expectedReasonCode))
			} else {
				volumeService.GetVolumeReturns(testcase.baseVolume, nil)
				volumeService.ExpandVolumeReturns(testcase.baseVolume, nil)
			}
			requestExp := provider.ExpandVolumeRequest{VolumeID: testcase.volumeID, Capacity: testcase.newSize}
			size, err := vpcs.ExpandVolume(requestExp)

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, reasoncode.ReasonCode(testcase.expectedReasonCode), util.ErrorReasonCode(err))
			}
			assert.Equal(t, size, testcase.expectedSize)
		})
	}
}
