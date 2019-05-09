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
	volumeServiceFakes "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/vpcvolume/fakes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestCreateVolume(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		volumeService *volumeServiceFakes.VolumeService
	)

	testCases := []struct {
		testCaseName   string
		baseVolume     *models.Volume
		providerVolume provider.Volume
		profileName    string

		setup func(providerVolume *provider.Volume)

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, volumeResponse *provider.Volume, err error)
	}{
		{
			testCaseName: "Volume capacity is nil",
			baseVolume: &models.Volume{
				ID:     "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:   "test volume name",
				Status: models.StatusType("OK"),
				Iops:   int64(1000),
				Zone:   &models.Zone{Name: "test-zone"},
			},
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("test volume name"),
				Capacity: nil,
			},
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Volume name is nil",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
			},
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Volume name is empty",
			baseVolume: &models.Volume{
				ID:     "16f293bf-test-4bff-816f-e199c0c65db5",
				Status: models.StatusType("OK"),
				Name:   "",
				Iops:   int64(1000),
				Zone:   &models.Zone{Name: "test-zone"},
			},
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String(""),
			},
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Volume capacity is zero",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("test volume name"),
				Capacity: Int(0),
			},
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Volume with general-purpose profile and invalid iops",
			profileName:  "general-purpose",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("test volume name"),
				Capacity: Int(10),
				Iops:     String("1000"),
				VPCVolume: provider.VPCVolume{
					Profile: &provider.Profile{Name: profileName},
				},
			},
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Volume with no validation issues",
			profileName:  "general-purpose",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("test volume name"),
				Capacity: Int(10),
				Iops:     String("0"),
				VPCVolume: provider.VPCVolume{
					Profile:       &provider.Profile{Name: profileName},
					ResourceGroup: &provider.ResourceGroup{ID: "default resource group id", Name: "default resource group"},
				},
			},
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
				assert.Nil(t, err)
			},
		}, {
			testCaseName: "Volume creaion failure",
			profileName:  "general-purpose",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("test volume name"),
				Capacity: Int(10),
				Iops:     String("0"),
				VPCVolume: provider.VPCVolume{
					Profile:       &provider.Profile{Name: profileName},
					ResourceGroup: &provider.ResourceGroup{ID: "default resource group id", Name: "default resource group"},
				},
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description: Volume creation failed. ",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Volume with test-purpose profile and invalid iops",
			profileName:  "test-purpose",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("test volume name"),
				Capacity: Int(10),
				VPCVolume: provider.VPCVolume{
					Profile: &provider.Profile{Name: profileName},
				},
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description: Volume creation failed. ",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Volume creaion with resource group ID and Name empty",
			profileName:  "general-purpose",
			providerVolume: provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("test volume name"),
				Capacity: Int(10),
				Iops:     String("0"),
				VPCVolume: provider.VPCVolume{
					Profile:       &provider.Profile{Name: profileName},
					ResourceGroup: &provider.ResourceGroup{},
				},
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description: Volume creation failed. ",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, volumeResponse *provider.Volume, err error) {
				assert.Nil(t, volumeResponse)
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

			if testcase.expectedErr != "" {
				volumeService.CreateVolumeReturns(testcase.baseVolume, errors.New(testcase.expectedReasonCode))
			} else {
				volumeService.CreateVolumeReturns(testcase.baseVolume, nil)
			}
			volume, err := vpcs.CreateVolume(testcase.providerVolume)
			logger.Info("Volume details", zap.Reflect("volume", volume))

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, reasoncode.ReasonCode(testcase.expectedReasonCode), util.ErrorReasonCode(err))
			}

			if testcase.verify != nil {
				testcase.verify(t, volume, err)
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
