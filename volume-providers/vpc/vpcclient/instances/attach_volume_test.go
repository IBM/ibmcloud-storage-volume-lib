/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package instances_test

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/instances"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/riaas/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"testing"
)

func TestAttachVolume(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	instanceID := "testinstance"

	testCases := []struct {
		name string

		// Response
		status  int
		content string

		// Expected return
		expectErr string
		verify    func(*testing.T, *models.VolumeAttachment, error)
	}{
		{
			name:   "Verify that the correct endpoint is invoked",
			status: http.StatusNoContent,
		}, {
			name:      "Verify that a 404 is returned to the caller",
			status:    http.StatusNotFound,
			content:   "{\"errors\":[{\"message\":\"testerr\"}]}",
			expectErr: "Trace Code:, testerr Please check ",
		}, {
			name:    "Verify that the volume attachment is done correctly",
			status:  http.StatusOK,
			content: "{\"id\":\"volume attachment id\", \"name\":\"volume attachment\", \"device\": {\"id\":\"xvdc\"}, \"volume\": {\"id\":\"volume-id\",\"name\":\"volume-name\",\"capacity\":10,\"iops\":3000,\"status\":\"pending\"}}",
			verify: func(t *testing.T, volumeAttachment *models.VolumeAttachment, err error) {
				if assert.NotNil(t, volumeAttachment) {
					assert.Equal(t, "volume attachment id", volumeAttachment.ID)
					assert.Equal(t, "xvdc", volumeAttachment.Device.ID)
				}
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {

			template := &models.VolumeAttachment{
				Name:       "volume attachment",
				InstanceID: &instanceID,
				Volume: &models.Volume{
					ID:       "volume-id",
					Name:     "volume-name",
					Capacity: 10,
					ResourceGroup: &models.ResourceGroup{
						ID: "rg1",
					},
					Generation: models.GenerationType("gc"),
					Zone:       &models.Zone{Name: "test-1"},
				},
			}

			mux, client, teardown := test.SetupServer(t)
			test.SetupMuxResponse(t, mux, "/v1/instances/testinstance/volume_attachments", http.MethodPost, nil, testcase.status, testcase.content, nil)

			defer teardown()

			logger.Info("Test case being executed", zap.Reflect("testcase", testcase.name))

			volumeAttachService := instances.New(client)

			volumeAttachment, err := volumeAttachService.AttachVolume(template, logger)
			logger.Info("Volume attachment details", zap.Reflect("volumeAttachment", volumeAttachment))

			if testcase.expectErr != "" && assert.Error(t, err) {
				assert.Equal(t, testcase.expectErr, err.Error())
				assert.Nil(t, volumeAttachment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, volumeAttachment)
			}

		})
	}
}

func TestIKSAttachVolume(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	instanceID := "testinstance"
	clusterID := "testcluster"
	// IKS tests
	mux, client, teardown := test.SetupServer(t)

	content := "{\"id\":\"volume attachment id\", \"name\":\"volume attachment\", \"device\": {\"id\":\"xvdc\"}, \"volume\": {\"id\":\"volume-id\",\"name\":\"volume-name\",\"capacity\":10,\"iops\":3000,\"status\":\"pending\"}}"

	test.SetupMuxResponse(t, mux, "/v2/storage/createAttachment", http.MethodPost, nil, http.StatusOK, content, nil)
	volumeAttachService := instances.NewIKSVolumeAttachmentManager(client)

	template := &models.VolumeAttachment{
		Name:       "volume attachment",
		ClusterID:  &clusterID,
		InstanceID: &instanceID,
		Volume: &models.Volume{
			ID:       "volume-id",
			Name:     "volume-name",
			Capacity: 10,
			ResourceGroup: &models.ResourceGroup{
				ID: "rg1",
			},
			Generation: models.GenerationType("gc"),
			Zone:       &models.Zone{Name: "test-1"},
		},
	}
	defer teardown()

	volumeAttachment, err := volumeAttachService.AttachVolume(template, logger)
	logger.Info("Volume attachment details", zap.Reflect("volumeAttachment", volumeAttachment))

	assert.NoError(t, err)
	assert.NotNil(t, volumeAttachment)
}

func GetTestContextLogger() (*zap.Logger, zap.AtomicLevel) {
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	traceLevel := zap.NewAtomicLevel()
	traceLevel.SetLevel(zap.InfoLevel)
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleDebugging, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return (lvl >= traceLevel.Level()) && (lvl < zapcore.ErrorLevel)
		})),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleErrors, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	)
	logger := zap.New(core, zap.AddCaller())
	return logger, traceLevel
}
