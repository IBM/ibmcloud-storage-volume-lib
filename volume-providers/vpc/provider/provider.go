/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/auth"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/iam"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/riaas"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

const (
	// displayName ...
	displayName = "IBM Cloud container service"
	// vpcProviderDisplayName ...
	vpcProviderDisplayName = "IBM Cloud infrastructure"
	// vpcExceptionPrefix ...
	vpcExceptionPrefix = "IBM Cloud infrastructure exception"
	// timeoutDefault ...
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
	APIConfig      riaas.Config
}

var _ local.Provider = &VPCBlockProvider{}

// NewProvider initialises an instance of an IaaS provider.
func NewProvider(conf *config.Config, logger *zap.Logger) (local.Provider, error) {
	logger.Info("Entering NewProvider")

	if conf.Bluemix == nil || conf.VPC == nil {
		return nil, errors.New("Incomplete config for VPCBlockProvider")
	}

	//Do config validation and enable only one generationType (i.e VPC-Classic | VPC-NG)
	gcConfigFound := (conf.VPC.EndpointURL != "") && (conf.VPC.TokenExchangeURL != "") && (conf.VPC.APIKey != "") && (conf.VPC.ResourceGroupID != "")
	g2ConfigFound := (conf.VPC.G2EndpointURL != "") && (conf.VPC.G2TokenExchangeURL != "") && (conf.VPC.G2APIKey != "") && (conf.VPC.G2ResourceGroupID != "")
	//if both config found, look for VPCTypeEnabled, otherwise default to GC
	//Incase of NG configurations, override the base properties.
	if (gcConfigFound && g2ConfigFound && conf.VPC.VPCTypeEnabled == "g2") || (!gcConfigFound && g2ConfigFound) {

		conf.VPC.EndpointURL = conf.VPC.G2EndpointURL
		conf.VPC.TokenExchangeURL = conf.VPC.G2TokenExchangeURL
		conf.VPC.APIKey = conf.VPC.G2APIKey
		conf.VPC.ResourceGroupID = conf.VPC.G2ResourceGroupID

		//Set API Generation As 2 (if unspecified in config/ENV-VAR)
		if conf.VPC.G2VPCAPIGeneration <= 0 {
			conf.VPC.G2VPCAPIGeneration = 2
		}
		conf.VPC.VPCAPIGeneration = conf.VPC.G2VPCAPIGeneration

		//Set the APIVersion Date, it can be diffrent in GC and NG
		if conf.VPC.G2APIVersion != "" {
			conf.VPC.APIVersion = conf.VPC.G2APIVersion
		}
	}

	// VPC provider use differnt APIkey and Auth Endpoint
	authConfig := &config.BluemixConfig{
		IamURL:          conf.VPC.TokenExchangeURL,
		IamAPIKey:       conf.VPC.APIKey,
		IamClientID:     conf.Bluemix.IamClientID,
		IamClientSecret: conf.Bluemix.IamClientSecret,
	}
	contextCF, err := auth.NewContextCredentialsFactory(authConfig, nil, conf.VPC)
	if err != nil {
		return nil, err
	}
	timeoutString := conf.VPC.VPCTimeout
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

	// SetRetryParameters sets the retry logic parameters
	SetRetryParameters(conf.VPC.MaxRetryAttempt, conf.VPC.MaxRetryGap)
	provider := &VPCBlockProvider{
		timeout:        timeout,
		serverConfig:   conf.Server,
		config:         conf.VPC,
		tokenGenerator: &tokenGenerator{config: conf.VPC},
		contextCF:      contextCF,
		httpClient:     httpClient,
		APIConfig: riaas.Config{
			BaseURL:       conf.VPC.EndpointURL,
			HTTPClient:    httpClient,
			APIVersion:    conf.VPC.APIVersion,
			APIGeneration: conf.VPC.VPCAPIGeneration,
		},
	}
	// Update VPC config for IKS deployment
	provider.config.IsIKS = conf.IKS != nil && conf.IKS.Enabled

	logger.Info("", zap.Reflect("Provider config", provider.config))

	userError.MessagesEn = messages.InitMessages()
	return provider, nil
}

// ContextCredentialsFactory ...
func (vpcp *VPCBlockProvider) ContextCredentialsFactory(zone *string) (local.ContextCredentialsFactory, error) {
	//  Datacenter hint not required by VPC provider implementation
	return vpcp.contextCF, nil
}

// OpenSession opens a session on the provider
func (vpcp *VPCBlockProvider) OpenSession(ctx context.Context, contextCredentials provider.ContextCredentials, ctxLogger *zap.Logger) (provider.Session, error) {
	ctxLogger.Info("Entering OpenSession")

	defer func() {
		ctxLogger.Debug("Exiting OpenSession")
	}()

	// validate that we have what we need - i.e. valid credentials
	if contextCredentials.Credential == "" {
		return nil, util.NewError("Error Insufficient Authentication", "No authentication credential provided")
	}

	if vpcp.serverConfig.DebugTrace {
		vpcp.APIConfig.DebugWriter = os.Stdout
	}

	if vpcp.ClientProvider == nil {
		vpcp.ClientProvider = riaas.DefaultRegionalAPIClientProvider{}
	}
	ctxLogger.Debug("", zap.Reflect("apiConfig.BaseURL", vpcp.APIConfig.BaseURL))

	if ctx != nil && ctx.Value(provider.RequestID) != nil {
		// set ContextID only of speicifed in the context
		vpcp.APIConfig.ContextID = fmt.Sprintf("%v", ctx.Value(provider.RequestID))
		ctxLogger.Info("", zap.Reflect("apiConfig.ContextID", vpcp.APIConfig.ContextID))
	}
	client, err := vpcp.ClientProvider.New(vpcp.APIConfig)
	if err != nil {
		return nil, err
	}

	// Create a token for all other API calls
	token, err := getAccessToken(contextCredentials, ctxLogger)
	if err != nil {
		return nil, err
	}
	ctxLogger.Debug("", zap.Reflect("Token", token.Token))

	err = client.Login(token.Token)
	if err != nil {
		return nil, err
	}

	// Update retry logic default values
	if vpcp.config.MaxRetryAttempt > 0 {
		ctxLogger.Debug("", zap.Reflect("MaxRetryAttempt", vpcp.config.MaxRetryAttempt))
		maxRetryAttempt = vpcp.config.MaxRetryAttempt
	}
	if vpcp.config.MaxRetryGap > 0 {
		ctxLogger.Debug("", zap.Reflect("MaxRetryGap", vpcp.config.MaxRetryGap))
		maxRetryGap = vpcp.config.MaxRetryGap
	}

	vpcSession := &VPCSession{
		VPCAccountID:          contextCredentials.IAMAccountID,
		Config:                vpcp.config,
		ContextCredentials:    contextCredentials,
		VolumeType:            "vpc-block",
		Provider:              VPC,
		Apiclient:             client,
		APIClientVolAttachMgr: client.VolumeAttachService(),
		Logger:                ctxLogger,
		APIRetry:              NewFlexyRetryDefault(),
	}
	return vpcSession, nil
}

// getAccessToken ...
func getAccessToken(creds provider.ContextCredentials, logger *zap.Logger) (token *iam.AccessToken, err error) {
	switch creds.AuthType {
	case provider.IAMAccessToken:
		token = &iam.AccessToken{Token: creds.Credential}
	default:
		err = errors.New("Unknown AuthType")
	}
	return
}
