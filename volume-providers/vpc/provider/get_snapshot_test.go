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
	"errors"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	serviceFakes "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/vpcvolume/fakes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestGetSnapshotWithVolumeID(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		snapshotService *serviceFakes.SnapshotService
		volumeService   *serviceFakes.VolumeService
	)

	testCases := []struct {
		testCaseName string
		volumeID     string
		snapshotID   string
		baseSnapshot *models.Snapshot
		baseVolume   *models.Volume
		setup        func()

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, snapshotResponse *provider.Snapshot, err error)
	}{
		{
			testCaseName: "Not supported yet",
			volumeID:     "v6f293bf-test-4bff-816f-e199c0c65db5",
			snapshotID:   "16f293bf-test-4bff-816f-e199c0c65db5",
			baseVolume: &models.Volume{
				ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     "test-volume-name",
				Status:   models.StatusType("OK"),
				Capacity: int64(10),
				Iops:     int64(1000),
				Zone:     &models.Zone{Name: "test-zone"},
			},
			baseSnapshot: &models.Snapshot{
				ID:     "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:   "test-snapshot-name",
				Status: models.StatusType("OK"),
			},
			verify: func(t *testing.T, snapshotResponse *provider.Snapshot, err error) {
				assert.NotNil(t, snapshotResponse)
				assert.Nil(t, err)
			},
		}, {
			testCaseName:       "Wrong volume ID",
			volumeID:           "Wrong volume ID",
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Wrong volume ID' volume ID is not valid. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_id_invalid, BackendError:, RC:400}",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, snapshotResponse *provider.Snapshot, err error) {
				assert.Nil(t, snapshotResponse)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.testCaseName, func(t *testing.T) {
			vpcs, uc, sc, err := GetTestOpenSession(t, logger)
			assert.NotNil(t, vpcs)
			assert.NotNil(t, uc)
			assert.NotNil(t, sc)
			assert.Nil(t, err)

			snapshotService = &serviceFakes.SnapshotService{}
			assert.NotNil(t, snapshotService)
			uc.SnapshotServiceReturns(snapshotService)

			volumeService = &serviceFakes.VolumeService{}
			assert.NotNil(t, volumeService)
			uc.VolumeServiceReturns(volumeService)

			if testcase.expectedErr != "" {
				snapshotService.GetSnapshotReturns(testcase.baseSnapshot, errors.New(testcase.expectedReasonCode))
				volumeService.GetVolumeReturns(testcase.baseVolume, errors.New(testcase.expectedReasonCode))
			} else {
				snapshotService.GetSnapshotReturns(testcase.baseSnapshot, nil)
				volumeService.GetVolumeReturns(testcase.baseVolume, nil)
			}
			snapshot, err := vpcs.GetSnapshotWithVolumeID(testcase.snapshotID, testcase.volumeID)
			logger.Info("Snapshot details", zap.Reflect("snapshot", snapshot))

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, reasoncode.ReasonCode(testcase.expectedReasonCode), util.ErrorReasonCode(err))
			}

			if testcase.verify != nil {
				testcase.verify(t, snapshot, err)
			}

		})
	}
}

func TestGetSnapshot(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		snapshotService *serviceFakes.SnapshotService
	)

	testCases := []struct {
		testCaseName string

		snapshotID   string
		baseSnapshot *models.Snapshot
		setup        func()

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, snapshotResponse *provider.Snapshot, err error)
	}{
		{
			testCaseName: "Not supported yet",
			snapshotID:   "16f293bf-test-4bff-816f-e199c0c65db5",
			verify: func(t *testing.T, snapshotResponse *provider.Snapshot, err error) {
				assert.Nil(t, snapshotResponse)
				assert.Nil(t, err)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.testCaseName, func(t *testing.T) {
			vpcs, uc, sc, err := GetTestOpenSession(t, logger)
			assert.NotNil(t, vpcs)
			assert.NotNil(t, uc)
			assert.NotNil(t, sc)
			assert.Nil(t, err)

			snapshotService = &serviceFakes.SnapshotService{}
			assert.NotNil(t, snapshotService)
			uc.SnapshotServiceReturns(snapshotService)

			if testcase.expectedErr != "" {
				snapshotService.GetSnapshotReturns(testcase.baseSnapshot, errors.New(testcase.expectedReasonCode))
			} else {
				snapshotService.GetSnapshotReturns(testcase.baseSnapshot, nil)
			}
			snapshot, err := vpcs.GetSnapshot(testcase.snapshotID)
			logger.Info("Snapshot details", zap.Reflect("snapshot", snapshot))

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, reasoncode.ReasonCode(testcase.expectedReasonCode), util.ErrorReasonCode(err))
			}

			if testcase.verify != nil {
				testcase.verify(t, snapshot, err)
			}

		})
	}
}
