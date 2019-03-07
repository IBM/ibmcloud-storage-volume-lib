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
	"context"
	"errors"
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/auth"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/iam"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/client/riaas"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

const (
	armadaDisplayName      = "IBM Cloud container service"
	vpcProviderDisplayName = "IBM Cloud infrastructure"
	vpcExceptionPrefix     = "IBM Cloud infrastructure exception"

	timeoutDefault = "120s"
)

// VPCBlockProvider implements provider.Provider
type VPCBlockProvider struct {
	timeout        time.Duration
	serverConfig   *config.ServerConfig
	config         *config.VPCProviderConfig
	tokenGenerator *tokenGenerator
	contextCF      local.ContextCredentialsFactory

	ClientProvider riaas.RegionalAPIClientProvider
	httpClient     *http.Client
}

var _ local.Provider = &VPCBlockProvider{}

// NewProvider initialises an instance of an IaaS provider.
func NewProvider(conf *config.Config, logger *zap.Logger) (local.Provider, error) {
	logger.Info("Entering NewProvider")

	if conf.Bluemix == nil || conf.VPC == nil {
		return nil, errors.New("Incomplete config for VPCBlockProvider")
	}

	contextCF, err := auth.NewContextCredentialsFactory(conf.Bluemix, nil, conf.VPC)
	if err != nil {
		return nil, err
	}
	timeoutString := conf.VPC.Timeout
	if timeoutString == "" || timeoutString == "0s" {
		logger.Info("Using VPC default timeout")
		timeoutString = "120s"
	}
	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		return nil, err
	}

	httpClient, err := config.GeneralCAHttpClientWithTimeout(timeout)
	if err != nil {
		logger.Error("Failed to prepare HTTP client", util.ZapError(err))
		return nil, err
	}

	provider := &VPCBlockProvider{
		timeout:        timeout,
		serverConfig:   conf.Server,
		config:         conf.VPC,
		tokenGenerator: &tokenGenerator{config: conf.VPC},
		contextCF:      contextCF,
		httpClient:     httpClient,
	}
	logger.Info("", zap.Reflect("Provider config", provider.config))
	return provider, nil
}

// ContextCredentialsFactory ...
func (vpcp *VPCBlockProvider) ContextCredentialsFactory(zone *string) (local.ContextCredentialsFactory, error) {
	//  Datacenter hint not required by VPC provider implementation
	return vpcp.contextCF, nil
}

// OpenSession opens a session on the provider
func (vpcp *VPCBlockProvider) OpenSession(ctx context.Context, contextCredentials provider.ContextCredentials, logger *zap.Logger) (provider.Session, error) {
	logger.Info("Entering OpenSession")

	defer func() {
		logger.Debug("Exiting OpenSession")
	}()

	// validate that we have what we need - i.e. valid credentials
	if contextCredentials.Credential == "" {
		return nil, util.NewError("Error Insufficient Authentication", "No authentication credential provided")
	}

	// Attempt to build an API client
	apiConfig := riaas.Config{
		BaseURL:    vpcp.config.EndpointURL,
		HTTPClient: vpcp.httpClient,
	}

	if vpcp.serverConfig.DebugTrace {
		apiConfig.DebugWriter = os.Stdout
	}

	if vpcp.ClientProvider == nil {
		vpcp.ClientProvider = riaas.DefaultRegionalAPIClientProvider{}
	}
	logger.Debug("", zap.Reflect("apiConfig.BaseURL", apiConfig.BaseURL))

	client, err := vpcp.ClientProvider.New(apiConfig)
	if err != nil {
		return nil, err
	}

	// Create a token for all other API calls
	token, err := getAccessToken(contextCredentials, logger)
	if err != nil {
		return nil, err
	}
	logger.Debug("", zap.Reflect("Token", token.Token))

	err = client.Login(token.Token)
	if err != nil {
		return nil, err
	}

	vpcSession := &VPCSession{
		VPCAccountID:       contextCredentials.IAMAccountID,
		Config:             vpcp.config,
		ContextCredentials: contextCredentials,
		VolumeType:         "vpc-block",
		Provider:           VPC,
		Apiclient:          client,
		Logger:             logger,
	}

	return vpcSession, nil
}

func getAccessToken(creds provider.ContextCredentials, logger *zap.Logger) (token *iam.AccessToken, err error) {
	switch creds.AuthType {
	case provider.IAMAccessToken:
		token = &iam.AccessToken{Token: creds.Credential}
	default:
		err = errors.New("Unknown AuthType")
	}
	return
}
