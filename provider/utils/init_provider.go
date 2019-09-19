/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package utils

import (
	"errors"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	softlayer_block "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/block"
	softlayer_file "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/file"
	iks_vpc_provider "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/iks/provider"
	vpc_provider "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/provider"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/registry"
)

const maxTimeout = 300  //secondsPerDay
const retryInterval = 5 //seconds
const maxRetryAttempts = maxTimeout / retryInterval

// InitProviders initialization for all providers as per configurations
func InitProviders(conf *config.Config, logger *zap.Logger) (registry.Providers, error) {
	var haveProviders bool
	providerRegistry := &registry.ProviderRegistry{}
	// BLOCK volume registration
	if conf.Softlayer != nil && conf.Softlayer.SoftlayerBlockEnabled {
		prov, err := softlayer_block.NewProvider(conf, logger)
		if err != nil {
			return nil, err
		}
		providerRegistry.Register(conf.Softlayer.SoftlayerBlockProviderName, prov)
		logger.Info("Block softlayer provider volume registry done!")

		haveProviders = true
	}

	// FILE volume registration
	if conf.Softlayer != nil && conf.Softlayer.SoftlayerFileEnabled {
		prov, err := softlayer_file.NewProvider(conf, logger)
		if err != nil {
			return nil, err
		}
		providerRegistry.Register(conf.Softlayer.SoftlayerFileProviderName, prov)
		logger.Info("File softlayer provider volume registry done!")

		haveProviders = true
	}

	// Genesis provider registration
	if conf.Gen2 != nil && conf.Gen2.Gen2ProviderEnabled {
		logger.Info("Configuring provider for Gen2")
		//TODO:~ Need to implement methods
		haveProviders = true
	}

	// VPC provider registration
	if conf.VPC != nil && conf.VPC.Enabled {
		logger.Info("Configuring VPC Block Provider")
		prov, err := vpc_provider.NewProvider(conf, logger)
		if err != nil {
			logger.Info("VPC block provider error!")
			return nil, err
		}
		providerRegistry.Register(conf.VPC.VPCBlockProviderName, prov)
		haveProviders = true
	}

	// IKS provider registration
	if conf.IKS != nil && conf.IKS.Enabled {
		logger.Info("Configuring IKS-VPC Block Provider")
		prov, err := iks_vpc_provider.NewProvider(conf, logger)
		if err != nil {
			logger.Info("VPC block provider error!")
			return nil, err
		}
		providerRegistry.Register(conf.IKS.IKSBlockProviderName, prov)
		haveProviders = true
	}

	if haveProviders {
		logger.Info("Provider registration done!!!")
		return providerRegistry, nil
	}

	return nil, errors.New("No providers registered")
}

// isEmptyStringValue ...
func isEmptyStringValue(value *string) bool {
	return value == nil || *value == ""
}

// OpenProviderSession ...
func OpenProviderSession(conf *config.Config, providers registry.Providers, providerID string, ctxLogger *zap.Logger) (session provider.Session, fatal bool, err error) {
	return OpenProviderSessionWithContext(nil, conf, providers, providerID, ctxLogger)
}

// OpenProviderSessionWithContext ...
func OpenProviderSessionWithContext(ctx context.Context, conf *config.Config, providers registry.Providers, providerID string, ctxLogger *zap.Logger) (session provider.Session, fatal bool, err error) {
	prov, err := providers.Get(providerID)
	if err != nil {
		ctxLogger.Error("Not able to get the said provider, might be its not registered", local.ZapError(err))
		fatal = true
		return
	}

	ccf, err := prov.ContextCredentialsFactory(&conf.Softlayer.SoftlayerDataCenter)
	if err != nil {
		fatal = true
		return
	}
	contextCredentials, err := GenerateContextCredentials(conf, providerID, ccf, ctxLogger)
	if err == nil {
		session, err = prov.OpenSession(ctx, contextCredentials, ctxLogger)
	}

	if err != nil {
		fatal = true
		ctxLogger.Error("Failed to open provider session", local.ZapError(err), zap.Bool("Fatal", fatal))
	}
	return
}

// GenerateContextCredentials ...
func GenerateContextCredentials(conf *config.Config, providerID string, contextCredentialsFactory local.ContextCredentialsFactory, ctxLogger *zap.Logger) (provider.ContextCredentials, error) {
	ctxLogger.Info("Generating generateContextCredentials for ", zap.String("Provider ID", providerID))

	AccountID := conf.Bluemix.IamClientID
	slUser := conf.Softlayer.SoftlayerUsername
	slAPIKey := conf.Softlayer.SoftlayerAPIKey
	iamAPIKey := conf.Bluemix.IamAPIKey

	// Select appropriate authentication strategy
	isSLProvider := providerID == conf.Softlayer.SoftlayerBlockProviderName || providerID == conf.Softlayer.SoftlayerFileProviderName
	switch {
	case isSLProvider && !isEmptyStringValue(&slUser) && !isEmptyStringValue(&slAPIKey):
		return contextCredentialsFactory.ForIaaSAPIKey(util.SafeStringValue(&AccountID), slUser, slAPIKey, ctxLogger)

	case (conf.VPC != nil && providerID == conf.VPC.VPCBlockProviderName):
		return contextCredentialsFactory.ForIAMAccessToken(conf.VPC.APIKey, ctxLogger)

	case isSLProvider && !isEmptyStringValue(&iamAPIKey):
		return contextCredentialsFactory.ForIAMAPIKey(AccountID, iamAPIKey, ctxLogger)

	case (conf.IKS != nil && providerID == conf.IKS.IKSBlockProviderName):
		return provider.ContextCredentials{}, nil // Get credentials  in OpenSession method

	default:
		return provider.ContextCredentials{}, util.NewError("ErrorInsufficientAuthentication",
			"Insufficient authentication credentials")
	}
}
