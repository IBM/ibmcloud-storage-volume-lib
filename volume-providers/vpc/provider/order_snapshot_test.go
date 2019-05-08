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
	"time"
)

func TestOrderSnapshot(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		snapshotService *serviceFakes.SnapshotService
		volumeService   *serviceFakes.VolumeService
	)

	testCases := []struct {
		name           string
		baseSnapshot   *models.Snapshot
		providerVolume provider.Volume
		baseVolume     *models.Volume
		tags           map[string]string
		setup          func()

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, err error)
	}{
		{
			name: "Not supported yet",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
			},
			baseVolume: &models.Volume{
				ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     "test-volume-name",
				Status:   models.StatusType("OK"),
				Capacity: int64(10),
				Iops:     int64(1000),
			},
			baseSnapshot: &models.Snapshot{
				ID:     "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:   "test-snapshot-name",
				Status: models.StatusType("OK"),
			},
			verify: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		}, {
			name: "Not supported yet",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
			},
			baseVolume: &models.Volume{
				ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     "test-volume-name",
				Status:   models.StatusType("OK"),
				Capacity: int64(10),
				Iops:     int64(1000),
			},
			expectedErr:        "{Code:StorageFindFailedWithSnapshotId, Type:InvalidRequest, Description:'Not a valid snapshot ID",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
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
				snapshotService.CreateSnapshotReturns(testcase.baseSnapshot, errors.New(testcase.expectedReasonCode))
				volumeService.GetVolumeReturns(testcase.baseVolume, errors.New(testcase.expectedReasonCode))
			} else {
				snapshotService.CreateSnapshotReturns(testcase.baseSnapshot, nil)
				volumeService.GetVolumeReturns(testcase.baseVolume, nil)
			}
			err = vpcs.OrderSnapshot(testcase.providerVolume)

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, reasoncode.ReasonCode(testcase.expectedReasonCode), util.ErrorReasonCode(err))
			}

			if testcase.verify != nil {
				testcase.verify(t, err)
			}

		})
	}
}

func TestOrderSnapshotTwo(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	timeNow := time.Now()

	var (
		snapshotService *serviceFakes.SnapshotService
		volumeService   *serviceFakes.VolumeService
		baseSnapshot    *models.Snapshot
		providerVolume  provider.Volume
		baseVolume      *models.Volume
	)

	providerVolume = provider.Volume{
		VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
	}
	baseVolume = &models.Volume{
		ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
		Name:     "test volume name",
		Status:   models.StatusType("OK"),
		Capacity: int64(10),
		Iops:     int64(1000),
		Zone:     &models.Zone{Name: "test-zone"},
	}
	baseSnapshot = &models.Snapshot{
		ID:        "16f293bf-test-4bff-816f-e199c0c65db5",
		Name:      "test-snapshot-name",
		Status:    models.StatusType("OK"),
		CreatedAt: &timeNow,
	}
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

	snapshotService.CreateSnapshotReturns(baseSnapshot, errors.New("ErrorUnclassified"))
	volumeService.GetVolumeReturns(baseVolume, nil)
	err = vpcs.OrderSnapshot(providerVolume)
	assert.NotNil(t, err)

	snapshotService.CreateSnapshotReturns(baseSnapshot, nil)
	volumeService.GetVolumeReturns(baseVolume, nil)
	err = vpcs.OrderSnapshot(providerVolume)
	assert.Nil(t, err)
}
