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
	//"errors"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	volumeServiceFakes "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/vpcvolume/fakes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestUpdateVolume(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		volumeService *volumeServiceFakes.VolumeService
	)

	testCases := []struct {
		testCaseName   string
		providerVolume provider.Volume
		profileName    string

		setup func(providerVolume *provider.Volume)

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, err error)
	}{
		{
			testCaseName: "Volume Update Success",
			providerVolume: provider.Volume{
				VolumeID:   "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:       String("test volume name"),
				Capacity:   nil,
				Provider:   provider.VolumeProvider("vpc-classic"),
				VolumeType: provider.VolumeType("block"),
			},
			verify: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		}, {
			testCaseName: "VolumeID Empty",
			providerVolume: provider.Volume{
				Name:     String("test volume name"),
				Capacity: Int(0),
			},
			expectedErr:        "{Code:ErrorRequiredFieldMissing, Type:InvalidRequest, Description:[VolumeID] is required to complete the operation., BackendError:, RC:400}",
			expectedReasonCode: "ErrorRequiredFieldMissing",
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Volume Provider Empty",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("test volume name"),
				Capacity: Int(0),
			},
			expectedErr:        "{Code:ErrorRequiredFieldMissing, Type:InvalidRequest, Description:[Provider] is required to complete the operation., BackendError:, RC:400}",
			expectedReasonCode: "ErrorRequiredFieldMissing",
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

			volumeService = &volumeServiceFakes.VolumeService{}
			assert.NotNil(t, volumeService)
			uc.VolumeServiceReturns(volumeService)
			/*
				if testcase.expectedErr != "" {
					volumeService.UpdateVolumeReturns(errors.New(testcase.expectedReasonCode))
				} else {
					volumeService.UpdateVolumeReturns(nil)
				}*/
			err = vpcs.UpdateVolume(testcase.providerVolume)

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, testcase.expectedErr, err.Error())
			}

			if testcase.verify != nil {
				testcase.verify(t, err)
			}

		})
	}
}

// String returns a pointer to the string value provided
func String(v string) *string {
	return &v
}

// Int returns a pointer to the int value provided
func Int(v int) *int {
	return &v
}
