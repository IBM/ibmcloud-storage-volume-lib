/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/metrics"

	"go.uber.org/zap"
)

// Config is for internal use by armada-cluster
type IAMConfiguration struct {
	IamURL          string `envconfig:"IAM_URL"`
	IamClientID     string `envconfig:"IAM_CLIENT_ID"`
	IamClientSecret string `envconfig:"IAM_CLIENT_SECRET"`
}

// Config is the parent struct for all the configuration information for armada-provider-vpc
type Config struct {
	Server    *ServerConfig `required:"true"`
	Metrics   *metrics.Config
	IAMConfig *IAMConfiguration

	VPCProviders map[string]VPCProviderConfig
}

// ServerConfig configuration options for the provider server itself
type ServerConfig struct {
	// DebugTrace is a flag to enable the debug level trace within the armada-provider-vpc code.
	DebugTrace bool `envconfig:"DEBUG_TRACE"`

	// VPCProviders contains a map of defined VPC provider name to config prefix (space separated)
	// e.g. "gc:GC gt:GT"
	VPCProviders string `envconfig:"VPC_PROVIDERS"`
}

// VPCProviderConfig configures a specific instance of a VPC provider (e.g. GT/GC/Z)
type VPCProviderConfig struct {
	Enabled     bool   `envconfig:"VPC_ENABLED"`
	Generation  string `envconfig:"VPC_GENERATION"`
	EndpointURL string `envconfig:"VPC_ENDPOINT_URL"`
	Timeout     string `envconfig:"VPC_TIMEOUT"`
}

// LoadConfig loads a Config from environment variables
func LoadConfig(logger *zap.Logger) (conf *Config, err error) {
	conf = &Config{
		VPCProviders: map[string]VPCProviderConfig{},
	}

	logger.Info("Reading configuration from environment")
	if err = envconfig.Process("", conf); err != nil {
		logger.Error("Error decoding environment variables", zap.Error(err))
		return nil, err
	}

	// Load individual VPC provider configurations
	err = config.LoadPrefixVarConfigs(conf.Server.VPCProviders, VPCProviderConfig{}, func(name string, value interface{}) {
		conf.VPCProviders[name] = *value.(*VPCProviderConfig)
	})
	if err != nil {
		logger.Error("Error decoding vpc providers", zap.Error(err))
		return nil, err
	}

	return conf, nil
}
