/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package softlayer_block

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/auth"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/common"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
)

const (
	armadaDisplayName     = "IBM Cloud container service"
	slProviderDisplayName = "IBM Cloud infrastructure"
	slExceptionPrefix     = "IBM Cloud infrastructure exception"

	timeoutDefault = "120s"
)

// SLBlockProvider implements provider.Provider
type SLBlockProvider struct {
	timeout        time.Duration
	config         *config.SoftlayerConfig
	tokenGenerator *tokenGenerator
	contextCF      local.ContextCredentialsFactory

	NewBackendSession func(url string, credentials provider.ContextCredentials, httpClient *http.Client, debug bool, logger *zap.Logger) backend.Session
}

var _ local.Provider = &SLBlockProvider{}

// NewProvider initialises an instance of an IaaS provider.
func NewProvider(conf *config.Config, logger *zap.Logger) (local.Provider, error) {
	if conf.Bluemix == nil || conf.Softlayer == nil {
		return nil, errors.New("Incomplete config for SLBlockProvider")
	}

	if conf.Softlayer.SoftlayerAPIDebug {
		logger.Warn("SoftlayerAPIDebug is enabled!")
	}

	contextCF, err := auth.NewContextCredentialsFactory(conf.Bluemix, conf.Softlayer)
	if err != nil {
		return nil, err
	}
	timeoutString := conf.Softlayer.SoftlayerTimeout
	if timeoutString == "" || timeoutString == "0s" {
		logger.Debug("Using Softlayer default Timeout")
		timeoutString = "120s"
	}
	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		return nil, err
	}

	provider := &SLBlockProvider{
		timeout:        timeout,
		config:         conf.Softlayer,
		tokenGenerator: &tokenGenerator{config: conf.Softlayer},
		contextCF:      contextCF,
	}

	return provider, nil
}

// ContextCredentialsFactory ...
func (slp *SLBlockProvider) ContextCredentialsFactory(zone *string) (local.ContextCredentialsFactory, error) {
	//  Datacenter hint not required by SL provider implementation
	return slp.contextCF, nil
}

// OpenSession opens a session on the provider
func (slp *SLBlockProvider) OpenSession(ctx context.Context, contextCredentials provider.ContextCredentials, logger *zap.Logger) (provider.Session, error) {

	slSession := &SLBlockSession{
		common.SLSession{
			Config: slp.config,
			//tokenGenerator:     slp.tokenGenerator,
			ContextCredentials: contextCredentials,
			Logger:             logger,
			VolumeType:         "ISCSI",
			Provider:           SoftLayer,
		},
	}
	// For method overriding
	//slSession.SLSession.SLSessionCommonInterface = slSession
	logr := logger.With(
		zap.Reflect("authType", contextCredentials.AuthType),
		zap.Reflect("timeout", slp.timeout),
	)

	switch contextCredentials.AuthType {
	case provider.IaaSAPIKey:
		slSession.Url = slp.config.SoftlayerEndpointURL
	case auth.IMSToken:
		slSession.Url = slp.config.SoftlayerIMSEndpointURL
	default:
		logr.Error("Unrecognised credentials")
		return nil, util.NewError("SLError-Session", "Unrecognised credentials")
	}

	logr = logger.With(
		zap.String("url", slSession.Url),
	)

	logr.Debug("Opening session to SoftLayer account")

	if slp.NewBackendSession == nil {
		// Use a session backed by a real softlayer-go Session
		slp.NewBackendSession = backend.NewSoftLayerSession
	}

	httpClient, err := config.GeneralCAHttpClientWithTimeout(slp.timeout)
	if err != nil {
		logr.Error("A problem occurred creating a generic HTTP Client", local.ZapError(err))
		// return nil, mapSLError(err, nil)  //: TODO: Neeed to map error with SL Error
		return nil, util.NewError("SLError-Session", "Error while creating genneric HTTP client")
	}

	// TODO CAN WE WIRE ctx TO THROUGH TO SOFTLAYER CLIENT?

	slSession.Backend = slp.NewBackendSession(slSession.Url, contextCredentials, httpClient, slp.config.SoftlayerAPIDebug, logger)

	slAccount, err := slSession.Backend.GetAccountService().Mask("id").GetObject()
	if err != nil {
		logr.Error("A problem occurred while retrieving the account ID", local.ZapError(err))
		//return nil, mapSLError(err, nil) //: TODO: Neeed to map error with SL Error
		return nil, util.NewError("SLError-ac", "Problem occurred while retrieving the account ID")
	}

	if slAccount.Id == nil {
		logr.Error("The SoftLayer account ID was not found")
		return nil, util.NewError("SLError-ac", "Provider account ID not found")
	}
	logr = logger.With(zap.Int("slAccountID", *slAccount.Id))
	slSession.SLAccountID = *slAccount.Id

	logr.Info("Opened session to SoftLayer account")

	return slSession, nil
}
