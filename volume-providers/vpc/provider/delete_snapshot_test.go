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

func TestDeleteSnapshot(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		snapshotService *serviceFakes.SnapshotService
		volumeService   *serviceFakes.VolumeService
	)

	testCases := []struct {
		testCaseName     string
		baseVolume       *models.Volume
		baseSnapshot     *models.Snapshot
		providerVolume   *provider.Volume
		providerSnapshot *provider.Snapshot

		tags  map[string]string
		setup func()

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, err error)
	}{
		{
			testCaseName: "Not supported yet",
			providerSnapshot: &provider.Snapshot{
				SnapshotID: "s6f293bf-test-4bff-816f-e199c0c65db5",
				Volume: provider.Volume{
					VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
					Name:     String("Test volume"),
					Capacity: Int(10),
					Iops:     String("1000"),
					VPCVolume: provider.VPCVolume{
						Profile: &provider.Profile{Name: "general-purpose"},
					},
				},
			},
			verify: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		}, {
			testCaseName: "Not a valid snapshot",
			providerSnapshot: &provider.Snapshot{
				SnapshotID: "s6f293bf-test-4bff-816f-e199c0c65db5",
				Volume: provider.Volume{
					VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
					Name:     String("Test volume"),
					Capacity: Int(10),
					Iops:     String("1000"),
					VPCVolume: provider.VPCVolume{
						Profile: &provider.Profile{Name: "general-purpose"},
					},
				},
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Not a valid snapshot ID",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, err error) {
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
				snapshotService.DeleteSnapshotReturns(errors.New(testcase.expectedReasonCode))
				snapshotService.GetSnapshotReturns(testcase.baseSnapshot, errors.New(testcase.expectedReasonCode))
			} else {
				snapshotService.DeleteSnapshotReturns(nil)
				snapshotService.GetSnapshotReturns(testcase.baseSnapshot, nil)
			}
			err = vpcs.DeleteSnapshot(testcase.providerSnapshot)

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
