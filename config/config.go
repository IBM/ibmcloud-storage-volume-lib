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
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
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
	Bluemix   *BluemixConfig //`required:"true"`
	Softlayer *SoftlayerConfig
	Gen2      *Gen2Config
}

func ReadConfig(confPath string, logger zap.Logger) *Config {
	// load the default config, if confPath not provided
	if confPath == "" {
		confPath = GetDefaultConfPath()
	}

	// Parse config file
	var conf Config
	logger.Info("parsing conf file", zap.String("confpath", confPath))
	ParseConfig(confPath, &conf, logger)

	return &conf
}

func GetConfPath() string {
	if confPath := getEnv("SECRET_CONFIG_PATH"); confPath != "" {
		return filepath.Join(confPath, "libconfig.toml")
	}
	//Get default conf path
	return GetDefaultConfPath()
}

func GetDefaultConfPath() string {
	return filepath.Join(GetEtcPath(), "libconfig.toml")
}

// ParseConfig ...
func ParseConfig(filePath string, conf interface{}, logger zap.Logger) {
	if _, err := toml.DecodeFile(filePath, conf); err != nil {
		logger.Fatal("error parsing config file", zap.Error(err))
	}
}

// BluemixConfig ...
type BluemixConfig struct {
	IamURL          string `toml:"iam_url"`
	IamClientID     string `toml:"iam_client_id"`
	IamClientSecret string `toml:"iam_client_secret" json:"-"`
	IamAPIKey       string `toml:"iam_api_key" json:"-"`
	RefreshToken    string `toml:"refresh_token" json:"-"`
}

// SoftlayerConfig ...
type SoftlayerConfig struct {
	SoftlayerBlockEnabled      bool   `toml:"softlayer_block_enabled" envconfig:"SOFTLAYER_BLOCK_ENABLED"`
	SoftlayerBlockProviderName string `toml:"softlayer_block_provider_name" envconfig:"SOFTLAYER_BLOCK_PROVIDER_NAME"`
	SoftlayerFileEnabled       bool   `toml:"softlayer_file_enabled" envconfig:"SOFTLAYER_FILE_ENABLED"`
	SoftlayerFileProviderName  string `toml:"softlayer_file_provider_name" envconfig:"SOFTLAYER_FILE_PROVIDER_NAME"`
	SoftlayerUsername          string `toml:"softlayer_username" json:"-"`
	SoftlayerAPIKey            string `toml:"softlayer_api_key" json:"-"`
	SoftlayerEndpointURL       string `toml:"softlayer_endpoint_url"`
	SoftlayerDataCenter        string `toml:"softlayer_datacenter"`
	SoftlayerTimeout           string `toml:"softlayer_api_timeout"`
	//Configuration values for JWT tokens
	SoftlayerJWTKID       string `toml:"softlayer_jwt_kid"`
	SoftlayerJWTTTL       int    `toml:"softlayer_jwt_ttl"`
	SoftlayerJWTValidFrom int    `toml:"softlayer_jwt_valid"`

	SoftlayerIMSEndpointURL string `toml:"softlayer_iam_endpoint_url"`
	SoftlayerAPIDebug       bool
}

type Gen2Config struct {
	Gen2ProviderEnabled bool   `toml:"genesis_provider_enabled"`
	Gen2Username        string `toml:"genesis_user_name"`
	Gen2APIKey          string `toml:"genesis_api_key"`
	Gen2URL             string `toml:"genesis_url"`
}

// GetEtcPath returns the path to the etc directory
func GetEtcPath() string {
	goPath := GetGoPath()
	srcPath := filepath.Join("src", "github.com", "IBM",
		"ibmcloud-storage-volume-lib")
	return filepath.Join(goPath, srcPath, "etc")
}
