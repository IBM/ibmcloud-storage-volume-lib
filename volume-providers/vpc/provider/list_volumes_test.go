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
	"strings"
	"testing"
)

func TestListVolumes(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		volumeService *volumeServiceFakes.VolumeService
	)

	testCases := []struct {
		testCaseName string
		volumeList   *models.VolumeList

		limit int
		start string
		tags  map[string]string

		setup func()

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, next_token string, volumes *provider.VolumeList, err error)
	}{
		{
			testCaseName: "Filter by zone",
			volumeList: &models.VolumeList{
				First: &models.HReference{Href: "https://eu-gb.iaas.cloud.ibm.com/v1/volumes?start=16f293bf-test-4bff-816f-e199c0c65db5\u0026limit=50\u0026zone.name=test-zone-1"},
				Next:  nil,
				Limit: 50,
				Volumes: []*models.Volume{
					{
						ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
						Name:     "test-volume-name1",
						Status:   models.StatusType("OK"),
						Capacity: int64(10),
						Iops:     int64(1000),
						Zone:     &models.Zone{Name: "test-zone-1"},
					}, {
						ID:       "23b154fr-test-4bff-816f-f213s1y34gj8",
						Name:     "test-volume-name2",
						Status:   models.StatusType("OK"),
						Capacity: int64(10),
						Iops:     int64(1000),
						Zone:     &models.Zone{Name: "test-zone-1"},
					},
				},
			},
			tags: map[string]string{
				"zone.name": "test-zone-1",
			},
			verify: func(t *testing.T, next_token string, volumes *provider.VolumeList, err error) {
				assert.NotNil(t, volumes.Volumes)
				assert.Equal(t, next_token, volumes.Next)
				assert.Nil(t, err)
			},
		}, {
			testCaseName: "Filter by zone, 1 entry per page",
			volumeList: &models.VolumeList{
				First: &models.HReference{Href: "https://eu-gb.iaas.cloud.ibm.com/v1/volumes?start=16f293bf-test-4bff-816f-e199c0c65db5\u0026limit=1\u0026zone.name=test-zone-1"},
				Next:  &models.HReference{Href: "https://eu-gb.iaas.cloud.ibm.com/v1/volumes?start=23b154fr-test-4bff-816f-f213s1y34gj8\u0026limit=1\u0026zone.name=test-zone-1"},
				Limit: 1,
				Volumes: []*models.Volume{
					{
						ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
						Name:     "test-volume-name1",
						Status:   models.StatusType("OK"),
						Capacity: int64(10),
						Iops:     int64(1000),
						Zone:     &models.Zone{Name: "test-zone-1"},
					}, {
						ID:       "23b154fr-test-4bff-816f-f213s1y34gj8",
						Name:     "test-volume-name2",
						Status:   models.StatusType("OK"),
						Capacity: int64(10),
						Iops:     int64(1000),
						Zone:     &models.Zone{Name: "test-zone-1"},
					},
				},
			},
			tags: map[string]string{
				"zone.name": "test-zone-1",
			},
			limit: 1,
			verify: func(t *testing.T, next_token string, volumes *provider.VolumeList, err error) {
				assert.NotNil(t, volumes.Volumes)
				assert.Equal(t, next_token, volumes.Next)
				assert.Nil(t, err)
			},
		}, {
			testCaseName: "Filter by zone: no volume found",    // Filter by zone where no volume is present
			volumeList: &models.VolumeList{
				First: nil,
				Next:  nil,
				Limit: 50,
				Volumes: []*models.Volume{},
			},
			tags: map[string]string{
				"zone.name": "test-zone",
			},
			verify: func(t *testing.T, next_token string, volumes *provider.VolumeList, err error) {
				assert.Nil(t, volumes.Volumes)
				assert.Equal(t, next_token, volumes.Next)
				assert.Nil(t, err)
			},
		}, {
			testCaseName: "Filter by name",
			volumeList: &models.VolumeList{
				First: &models.HReference{Href: "https://eu-gb.iaas.cloud.ibm.com/v1/volumes?start=16f293bf-test-4bff-816f-e199c0c65db5\u0026limit=50"},
				Next:  nil,
				Limit: 50,
				Volumes: []*models.Volume{
					{
						ID:       "16f293bf-test-4bff-816f-e199c0c65db5",
						Name:     "test-volume-name1",
						Status:   models.StatusType("OK"),
						Capacity: int64(10),
						Iops:     int64(1000),
						Zone:     &models.Zone{Name: "test-zone"},
					},
				},
			},
			tags: map[string]string{
				"name": "test-volume-name1",
			},
			verify: func(t *testing.T, next_token string, volumes *provider.VolumeList, err error) {
				assert.NotNil(t, volumes.Volumes)
				assert.Equal(t, next_token, volumes.Next)
				assert.Nil(t, err)
			},
		}, {
			testCaseName: "volume not found",
			tags: map[string]string{
				"name": "test-volume-name1",
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:RetrivalFailed, Description: Unable to fetch list of volumes. ",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, next_token string, volumes *provider.VolumeList, err error) {
				assert.Nil(t, volumes)
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
				volumeService.ListVolumesReturns(testcase.volumeList, errors.New(testcase.expectedReasonCode))
			} else {
				volumeService.ListVolumesReturns(testcase.volumeList, nil)
			}
			volumes, err := vpcs.ListVolumes(testcase.limit, testcase.start, testcase.tags)
			logger.Info("VolumesList details", zap.Reflect("VolumesList", volumes))

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, reasoncode.ReasonCode(testcase.expectedReasonCode), util.ErrorReasonCode(err))
			}

			if testcase.verify != nil {
				var next string
				if testcase.volumeList != nil {
					if testcase.volumeList.Next != nil {
						// "Next":{"href":"https://eu-gb.iaas.cloud.ibm.com/v1/volumes?start=3e898aa7-ac71-4323-952d-a8d741c65a68\u0026limit=1\u0026zone.name=eu-gb-1"}
						next = strings.Split(strings.Split(testcase.volumeList.Next.Href, "start=")[1], "\u0026")[0]
					}
				}
				testcase.verify(t, next, volumes, err)
			}

		})
	}
}
