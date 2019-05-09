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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	snapshotServiceFakes "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/vpcvolume/fakes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestListSnapshots(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		snapshotService *snapshotServiceFakes.SnapshotService
	)

	testCases := []struct {
		testCaseName string
		volumeID     string
		baseVolume   *models.Volume

		setup func()

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, volumes []*provider.Snapshot, err error)
	}{
		{
			testCaseName: "Not supported yet",
			verify: func(t *testing.T, snapshots []*provider.Snapshot, err error) {
				assert.Nil(t, snapshots)
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

			snapshotService = &snapshotServiceFakes.SnapshotService{}
			assert.NotNil(t, snapshotService)
			uc.SnapshotServiceReturns(snapshotService)

			snapshotService.ListSnapshotsReturns(nil, nil)

			snapshots, err := vpcs.ListSnapshots()
			logger.Info("Snapshots details", zap.Reflect("Snapshots", snapshots))

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, reasoncode.ReasonCode(testcase.expectedReasonCode), util.ErrorReasonCode(err))
			}

			if testcase.verify != nil {
				testcase.verify(t, snapshots, err)
			}

		})
	}
}
