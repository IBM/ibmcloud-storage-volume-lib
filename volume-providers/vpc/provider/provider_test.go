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
	"context"
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/auth"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/riaas/fakes"
	volumeServiceFakes "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/vpcvolume/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	TestProviderAccountID   = "test-provider-account"
	TestProviderAccessToken = "test-provider-access-token"
	TestIKSAccountID        = "test-iks-account"
	TestZone                = "test-zone"
	IamURL                  = "test-iam-url"
	IamClientID             = "test-iam_client_id"
	IamClientSecret         = "test-iam_client_secret"
	IamAPIKey               = "test-iam_api_key"
	RefreshToken            = "test-refresh_token"
	TestEndpointURL         = "http://some_endpoint"
	TestApiVersion          = "2019-01-01"
)

var _ local.ContextCredentialsFactory = &auth.ContextCredentialsFactory{}

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
			EndpointURL: TestEndpointURL,
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
			EndpointURL: TestEndpointURL,
			Timeout:     "",
		},
	}

	prov, err = NewProvider(conf, logger)
	assert.NotNil(t, prov)
	assert.Nil(t, err)

	zone := "Test Zone"
	contextCF, _ := prov.ContextCredentialsFactory(&zone)
	assert.NotNil(t, contextCF)

	return
}

func GetTestProvider(t *testing.T, logger *zap.Logger) (*VPCBlockProvider, error) {
	var cp *fakes.RegionalAPIClientProvider
	var uc, sc *fakes.RegionalAPI

	logger.Info("Getting New test Provider")
	conf := &config.Config{
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
			Enabled:         true,
			EndpointURL:     TestEndpointURL,
			Timeout:         "30s",
			MaxRetryAttempt: 5,
			MaxRetryGap:     10,
			APIVersion:      TestApiVersion,
		},
	}

	p, err := NewProvider(conf, logger)
	assert.NotNil(t, p)
	assert.Nil(t, err)

	timeout, _ := time.ParseDuration(conf.VPC.Timeout)

	// Inject a fake RIAAS API client
	cp = &fakes.RegionalAPIClientProvider{}
	uc = &fakes.RegionalAPI{}
	cp.NewReturnsOnCall(0, uc, nil)
	sc = &fakes.RegionalAPI{}
	cp.NewReturnsOnCall(1, sc, nil)

	volumeService := &volumeServiceFakes.VolumeService{}
	uc.VolumeServiceReturns(volumeService)

	httpClient, err := config.GeneralCAHttpClientWithTimeout(timeout)
	if err != nil {
		logger.Error("Failed to prepare HTTP client", util.ZapError(err))
		return nil, err
	}
	assert.NotNil(t, httpClient)

	provider := &VPCBlockProvider{
		timeout:        timeout,
		serverConfig:   conf.Server,
		config:         conf.VPC,
		tokenGenerator: &tokenGenerator{config: conf.VPC},
		httpClient:     httpClient,
	}
	assert.NotNil(t, provider)
	assert.Equal(t, provider.timeout, timeout)

	return provider, nil
}

func TestGetTestProvider(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	prov, err := GetTestProvider(t, logger)
	assert.NotNil(t, prov)
	assert.Nil(t, err)

	zone := "Test Zone"
	contextCF, _ := prov.ContextCredentialsFactory(&zone)
	assert.Nil(t, contextCF)
	assert.NotNil(t, prov.httpClient)
}

func TestOpenSession(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	vpcp, err := GetTestProvider(t, logger)

	sessn, err := vpcp.OpenSession(context.Background(), provider.ContextCredentials{
		AuthType:     provider.IAMAccessToken,
		Credential:   TestProviderAccessToken,
		IAMAccountID: TestIKSAccountID,
	}, logger)

	require.NoError(t, err)
	assert.NotNil(t, sessn)

	sessn, err = vpcp.OpenSession(context.Background(), provider.ContextCredentials{
		AuthType:     provider.IAMAccessToken,
		IAMAccountID: TestIKSAccountID,
	}, logger)

	require.Error(t, err)
	assert.Nil(t, sessn)

	sessn, err = vpcp.OpenSession(context.Background(), provider.ContextCredentials{
		AuthType:     "WrongType",
		IAMAccountID: TestIKSAccountID,
	}, logger)

	require.Error(t, err)
	assert.Nil(t, sessn)

	return
}

func GetTestOpenSession(t *testing.T, logger *zap.Logger) (sessn provider.Session, uc, sc *fakes.RegionalAPI, err error) {
	vpcp, err := GetTestProvider(t, logger)

	m := http.NewServeMux()
	s := httptest.NewServer(m)
	assert.NotNil(t, s)

	vpcp.httpClient = http.DefaultClient

	// Inject a fake RIAAS API client
	cp := &fakes.RegionalAPIClientProvider{}
	uc = &fakes.RegionalAPI{}
	cp.NewReturnsOnCall(0, uc, nil)
	sc = &fakes.RegionalAPI{}
	cp.NewReturnsOnCall(1, sc, nil)
	vpcp.ClientProvider = cp

	sessn = &VPCSession{
		VPCAccountID: TestIKSAccountID,
		Config:       vpcp.config,
		ContextCredentials: provider.ContextCredentials{
			AuthType:     provider.IAMAccessToken,
			Credential:   TestProviderAccessToken,
			IAMAccountID: TestIKSAccountID,
		},
		VolumeType: "vpc-block",
		Provider:   VPC,
		Apiclient:  uc,
		Logger:     logger,
	}

	return
}

func TestGetTestOpenSession(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	vpcs, uc, sc, err := GetTestOpenSession(t, logger)
	assert.NotNil(t, vpcs)
	assert.NotNil(t, uc)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	providerDisplayName := vpcs.GetProviderDisplayName()
	assert.Equal(t, providerDisplayName, provider.VolumeProvider("VPC"))
	vpcs.Close()

	providerName := vpcs.ProviderName()
	assert.Equal(t, providerName, provider.VolumeProvider("VPC"))

	volumeType := vpcs.Type()
	assert.Equal(t, volumeType, provider.VolumeType("vpc-block"))

	volume, _ := vpcs.GetVolume("test volume")
	assert.Nil(t, volume)
}
