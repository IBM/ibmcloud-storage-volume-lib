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

func TestDeleteVolume(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		volumeService *serviceFakes.VolumeService
	)

	testCases := []struct {
		testCaseName   string
		baseVolume     *models.Volume
		providerVolume *provider.Volume

		tags  map[string]string
		setup func()

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, err error)
	}{
		{
			testCaseName: "Not supported yet",
			providerVolume: &provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("Test volume"),
				Capacity: Int(10),
				Iops:     String("1000"),
				VPCVolume: provider.VPCVolume{
					Profile: &provider.Profile{Name: "general-purpose"},
				},
			},
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		}, {
			testCaseName:       "False positive: No volume being sent",
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Not a valid volume ID",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Incorrect volume ID",
			providerVolume: &provider.Volume{
				VolumeID: "wrong volume ID",
				Name:     String("Test volume"),
				Capacity: Int(10),
				Iops:     String("1000"),
				VPCVolume: provider.VPCVolume{
					Profile:       &provider.Profile{Name: "general-purpose"},
					ResourceGroup: &provider.ResourceGroup{ID: "default resource group id", Name: "default resource group"},
				},
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Not a valid volume ID",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Incorrect volume ID",
			providerVolume: &provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("Test volume"),
				Capacity: Int(10),
				Iops:     String("1000"),
				VPCVolume: provider.VPCVolume{
					Profile:       &provider.Profile{Name: "general-purpose"},
					ResourceGroup: &provider.ResourceGroup{ID: "default resource group id", Name: "default resource group"},
				},
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Not a valid volume ID",
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

			volumeService = &serviceFakes.VolumeService{}
			assert.NotNil(t, volumeService)
			uc.VolumeServiceReturns(volumeService)

			if testcase.expectedErr != "" {
				volumeService.DeleteVolumeReturns(errors.New(testcase.expectedReasonCode))
				volumeService.GetVolumeReturns(testcase.baseVolume, errors.New(testcase.expectedReasonCode))
			} else {
				volumeService.DeleteVolumeReturns(nil)
				volumeService.GetVolumeReturns(testcase.baseVolume, nil)
			}

			err = vpcs.DeleteVolume(testcase.providerVolume)

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

func TestDeleteVolumeTwo(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		volumeService *serviceFakes.VolumeService
	)

	var baseVolume *models.Volume
	var providerVolume *provider.Volume

	baseVolume = &models.Volume{
		ID:     "16f293bf-test-4bff-816f-e199c0c65db5",
		Name:   "test volume name",
		Status: models.StatusType("OK"),
		Iops:   int64(1000),
		Zone:   &models.Zone{Name: "test-zone"},
	}

	providerVolume = &provider.Volume{
		VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
		Name:     String("Test volume"),
		Capacity: Int(10),
		Iops:     String("1000"),
		VPCVolume: provider.VPCVolume{
			Profile:       &provider.Profile{Name: "general-purpose"},
			ResourceGroup: &provider.ResourceGroup{ID: "default resource group id", Name: "default resource group"},
		},
	}

	vpcs, uc, sc, err := GetTestOpenSession(t, logger)
	assert.NotNil(t, vpcs)
	assert.NotNil(t, uc)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	volumeService = &serviceFakes.VolumeService{}
	assert.NotNil(t, volumeService)
	uc.VolumeServiceReturns(volumeService)

	volumeService.DeleteVolumeReturns(nil)
	volumeService.GetVolumeReturns(nil, errors.New("ErrorUnclassified"))

	err = vpcs.DeleteVolume(providerVolume)
	assert.Nil(t, err)

	volumeService.DeleteVolumeReturns(errors.New("FailedToDeleteVolume"))
	volumeService.GetVolumeReturns(baseVolume, nil)

	err = vpcs.DeleteVolume(providerVolume)
	assert.NotNil(t, err)
}
