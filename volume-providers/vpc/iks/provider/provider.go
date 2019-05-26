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
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/auth"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/iam"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	vpcprovider "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/provider"
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

type IksVpcBlockProvider struct {
	vpcprovider.VPCBlockProvider
	conf *config.Config
}

var _ local.Provider = &IksVpcBlockProvider{}

func NewIksVpcProvider(conf *config.Config, logger *zap.Logger) (local.Provider, error) {

	provider, _ := vpcprovider.NewProvider(conf, logger)
	vpcBlockProvider, _ := provider.(*vpcprovider.VPCBlockProvider)
	iksVpcBlockProvider := &IksVpcBlockProvider{
		VPCBlockProvider: *vpcBlockProvider,
	}
	iksVpcBlockProvider.VPCBlockProvider.getAPIConfig = func() riaas.Config {
		return riaas.Config{
			BaseURL: conf.Bluemix.IamURL,
			//HTTPClient: httpClient,
			//APIVersion: conf.VPC.APIVersion,
		}
	}
	return iksVpcBlockProvider, nil

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

	apiConfig := vpcp.getAPIConfig()

	if vpcp.serverConfig.DebugTrace {
		apiConfig.DebugWriter = os.Stdout
	}

	if vpcp.ClientProvider == nil {
		vpcp.ClientProvider = riaas.DefaultRegionalAPIClientProvider{}
	}
	ctxLogger.Debug("", zap.Reflect("apiConfig.BaseURL", apiConfig.BaseURL))

	client, err := vpcp.ClientProvider.New(apiConfig)
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
		VPCAccountID:       contextCredentials.IAMAccountID,
		Config:             vpcp.config,
		ContextCredentials: contextCredentials,
		VolumeType:         "vpc-block",
		Provider:           VPC,
		Apiclient:          client,
		Logger:             ctxLogger,
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
