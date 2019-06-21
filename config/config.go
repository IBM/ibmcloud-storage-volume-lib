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
	"github.com/BurntSushi/toml"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getEnv(key string) string {
	return os.Getenv(strings.ToUpper(key))
}

// GetGoPath inspects the environment for the GOPATH variable
func GetGoPath() string {
	if goPath := getEnv("GOPATH"); goPath != "" {
		return goPath
	}
	return ""
}

// Config is the parent struct for all the configuration information for -cluster
type Config struct {
	Server    *ServerConfig  `required:"true"`
	Bluemix   *BluemixConfig //`required:"true"`
	Softlayer *SoftlayerConfig
	Gen2      *Gen2Config
	VPC       *VPCProviderConfig
	IKS       *IKSConfig
}

//ReadConfig loads the config from file
func ReadConfig(confPath string, logger *zap.Logger) (*Config, error) {
	// load the default config, if confPath not provided
	if confPath == "" {
		confPath = GetDefaultConfPath()
	}

	// Parse config file
	conf := Config{
		IKS: &IKSConfig{}, // IKS block may not be populated in secrete toml. Make sure its not nil
	}
	logger.Info("parsing conf file", zap.String("confpath", confPath))
	err := ParseConfig(confPath, &conf, logger)
	return &conf, err
}

// GetConfPath get configuration file path
func GetConfPath() string {
	if confPath := getEnv("SECRET_CONFIG_PATH"); confPath != "" {
		return filepath.Join(confPath, "libconfig.toml")
	}
	//Get default conf path
	return GetDefaultConfPath()
}

// GetConfPathDir get configuration  dir path
func GetConfPathDir() string {
	if confPath := getEnv("SECRET_CONFIG_PATH"); confPath != "" {
		return confPath
	}
	//Get default conf path
	return GetEtcPath()
}

// GetDefaultConfPath get default config file path
func GetDefaultConfPath() string {
	return filepath.Join(GetEtcPath(), "libconfig.toml")
}

// ParseConfig ...
func ParseConfig(filePath string, conf interface{}, logger *zap.Logger) error {
	_, err := toml.DecodeFile(filePath, conf)
	if err != nil {
		logger.Error("Failed to parse config file", zap.Error(err))
	}
	// Read environment variables
	err = envconfig.Process("", conf)
	if err != nil {
		logger.Error("Failed to gather environment config variable", zap.Error(err))
	}
	return err
}

// ServerConfig configuration options for the provider server itself
type ServerConfig struct {
	// DebugTrace is a flag to enable the debug level trace within the provider code.
	DebugTrace bool `toml:"debug_trace" envconfig:"DEBUG_TRACE"`
}

// BluemixConfig ...
type BluemixConfig struct {
	IamURL          string `toml:"iam_url"`
	IamClientID     string `toml:"iam_client_id"`
	IamClientSecret string `toml:"iam_client_secret" json:"-"`
	IamAPIKey       string `toml:"iam_api_key" json:"-"`
	RefreshToken    string `toml:"refresh_token" json:"-"`
	APIEndpointURL  string `toml:"containers_api_route"`
	Encryption      bool   `toml:"encryption"`
}

// SoftlayerConfig ...
type SoftlayerConfig struct {
	SoftlayerBlockEnabled        bool   `toml:"softlayer_block_enabled" envconfig:"SOFTLAYER_BLOCK_ENABLED"`
	SoftlayerBlockProviderName   string `toml:"softlayer_block_provider_name" envconfig:"SOFTLAYER_BLOCK_PROVIDER_NAME"`
	SoftlayerFileEnabled         bool   `toml:"softlayer_file_enabled" envconfig:"SOFTLAYER_FILE_ENABLED"`
	SoftlayerFileProviderName    string `toml:"softlayer_file_provider_name" envconfig:"SOFTLAYER_FILE_PROVIDER_NAME"`
	SoftlayerUsername            string `toml:"softlayer_username" json:"-"`
	SoftlayerAPIKey              string `toml:"softlayer_api_key" json:"-"`
	SoftlayerEndpointURL         string `toml:"softlayer_endpoint_url"`
	SoftlayerDataCenter          string `toml:"softlayer_datacenter"`
	SoftlayerTimeout             string `toml:"softlayer_api_timeout" envconfig:"SOFTLAYER_API_TIMEOUT"`
	SoftlayerVolProvisionTimeout string `toml:"softlayer_vol_provision_timeout" envconfig:"SOFTLAYER_VOL_PROVISION_TIMEOUT"`
	SoftlayerRetryInterval       string `toml:"softlayer_api_retry_interval" envconfig:"SOFTLAYER_API_RETRY_INTERVAL"`

	//Configuration values for JWT tokens
	SoftlayerJWTKID       string `toml:"softlayer_jwt_kid"`
	SoftlayerJWTTTL       int    `toml:"softlayer_jwt_ttl"`
	SoftlayerJWTValidFrom int    `toml:"softlayer_jwt_valid"`

	SoftlayerIMSEndpointURL string `toml:"softlayer_iam_endpoint_url"`
	SoftlayerAPIDebug       bool
}

// Gen2Config ...
type Gen2Config struct {
	Gen2ProviderEnabled bool   `toml:"genesis_provider_enabled"`
	Gen2Username        string `toml:"genesis_user_name"`
	Gen2APIKey          string `toml:"genesis_api_key"`
	Gen2URL             string `toml:"genesis_url"`
}

// VPCProviderConfig configures a specific instance of a VPC provider (e.g. GT/GC/Z)
type VPCProviderConfig struct {
	Enabled              bool   `toml:"vpc_enabled" envconfig:"VPC_ENABLED"`
	VPCBlockProviderName string `toml:"vpc_block_provider_name" envconfig:"VPC_BLOCK_PROVIDER_NAME"`
	EndpointURL          string `toml:"gc_riaas_endpoint_url"`
	TokenExchangeURL     string `toml:"gc_token_exchange_endpoint_url"`
	APIKey               string `toml:"gc_api_key" json:"-"`
	Encryption           bool   `toml:"encryption"`
	ResourceGroupID      string `toml:"gc_resource_group_id"`
	Timeout              string `toml:"vpc_timeout" envconfig:"VPC_TIMEOUT"`
	MaxRetryAttempt      int    `toml:"max_retry_attempt"`
	MaxRetryGap          int    `toml:"max_retry_gap" envconfig:"VPC_RETRY_INTERVAL"`
	APIVersion           string `toml:"api_version" envconfig:"VPC_API_VERSION"`
}

//IKSConfig config
type IKSConfig struct {
	Enabled              bool   `toml:"iks_enabled" envconfig:"IKS_ENABLED"`
	IKSBlockProviderName string `toml:"iks_block_provider_name" envconfig:"IKS_BLOCK_PROVIDER_NAME"`
}

// GetEtcPath returns the path to the etc directory
func GetEtcPath() string {
	goPath := GetGoPath()
	srcPath := filepath.Join("src", "github.com", "IBM",
		"ibmcloud-storage-volume-lib")
	return filepath.Join(goPath, srcPath, "etc")
}

//GetTimeOutParameters retrives the parameteer to implement retry logic
// maxTimeout - Maximum time out for the operations
//retryGapDuration - The time interval for next attempt
// maxRetryAttempt - maxmum retry attempts derived based on  maxTimeout and retryGapDuration
func (config *VPCProviderConfig) GetTimeOutParameters() (int, int, time.Duration) {
	maxTimeoutConfig, _ := time.ParseDuration(config.Timeout)
	maxTimeout := int(maxTimeoutConfig.Seconds())
	maxRetryAttempt := maxTimeout / config.MaxRetryGap
	retryGapDuration := time.Duration(config.MaxRetryGap) * time.Second
	return maxTimeout, maxRetryAttempt, retryGapDuration
}
