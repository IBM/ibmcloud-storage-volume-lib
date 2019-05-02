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
	"bytes"
	"testing"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	TestProviderAccountID   = "test-provider-account"
	TestProviderAccessToken = "test-provider-access-token"
	TestIKSAccountID        = "test-iks-account"
	IamURL                  = "test-iam-url"
	IamClientID             = "test-iam_client_id"
	IamClientSecret         = "test-iam_client_secret"
	IamAPIKey               = "test-iam_api_key"
	RefreshToken            = "test-refresh_token"
)

func GetTestLogger(t *testing.T) (logger *zap.Logger, teardown func()) {

	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.DebugLevel)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	buf := &bytes.Buffer{}

	logger = zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(buf),
			atom,
		),
		zap.AddCaller(),
	)

	teardown = func() {

		logger.Sync()

		if t.Failed() {
			t.Log(buf)
		}
	}

	return

}

func TestNewProvider(t *testing.T) {
	var err error
	conf := &config.Config{
		Server: &config.ServerConfig{
			DebugTrace: true,
		},
		VPC: &config.VPCProviderConfig{
			Enabled:     true,
			EndpointURL: "http://some_endpoint",
			Timeout:     "30s",
		},
	}
	logger, teardown := GetTestLogger(t)
	defer teardown()
	prov, err := NewProvider(conf, logger)
	assert.Nil(t, prov)
	assert.NotNil(t, err)

	conf = &config.Config{
		Server: &config.ServerConfig{
			DebugTrace: true,
		},
		Bluemix: &config.BluemixConfig{
			IamURL:          IamURL,
			IamClientID:     IamClientID,
			IamClientSecret: IamClientSecret,
			IamAPIKey:       IamClientSecret,
			RefreshToken:    RefreshToken,
		},
		VPC: &config.VPCProviderConfig{
			Enabled:     true,
			EndpointURL: "http://some_endpoint",
			Timeout:     "30s",
		},
	}
	prov, err = NewProvider(conf, logger)
	assert.NotNil(t, prov)
	assert.Nil(t, err)

	return
}
