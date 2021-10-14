/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpc

import (
	"context"
	"fmt"
	"testing"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"

	userError "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-volume-interface/config"
	"github.com/IBM/ibmcloud-volume-interface/lib/provider"
	"github.com/IBM/ibmcloud-volume-interface/provider/local"
	provider_util "github.com/IBM/ibmcloud-volume-vpc/block/utils"
	vpcconfig "github.com/IBM/ibmcloud-volume-vpc/block/vpcconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sess provider.Session
var logger *zap.Logger
var ctxLogger *zap.Logger
var traceLevel zap.AtomicLevel
var requestID string
var resourceGroupID string
var vpcZone string
var volumeEncryptionKeyCRN string
var startTime time.Time

func TestVPCE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ibmcloud-storage-volume-lib VPC e2e test suite")
}

var _ = BeforeSuite(func() {
	// Setup new style zap logger
	logger, traceLevel = getContextLogger()
	defer logger.Sync()
	// Load config file
	goPath := os.Getenv("GOPATH")
	conf, err := config.ReadConfig(goPath+vpcConfigFilePath, logger)
	if err != nil {
		logger.Fatal("Error loading configuration")
		Expect(err).To(HaveOccurred())
	}

	// Check if debug log level enabled or not
	if conf.Server != nil && conf.Server.DebugTrace {
		traceLevel.SetLevel(zap.DebugLevel)
	}

	// Load zone info
	vpcZone = os.Getenv("VPC_ZONE")

	// Load EncryptionKeyCRN info
	volumeEncryptionKeyCRN = os.Getenv("ENCRYPTION_KEY_CRN")

	// Get only VPC_API_VERSION, in "2019-07-02T00:00:00.000Z" case vpc need only 2019-07-02"
	dateTime, err := time.Parse(time.RFC3339, conf.VPC.APIVersion)
	if err == nil {
		conf.VPC.APIVersion = fmt.Sprintf("%d-%02d-%02d", dateTime.Year(), dateTime.Month(), dateTime.Day())
	} else {
		logger.Warn("Failed to parse VPC_API_VERSION, setting default value")
		conf.VPC.APIVersion = "2020-07-02" // setting default values
	}

	// Update the CSRF  Token
	if conf.Bluemix.PrivateAPIRoute != "" {
		conf.Bluemix.CSRFToken = string([]byte{}) // TODO~ Need to remove it
	}

	if conf.API == nil {
		conf.API = &config.APIConfig{
			PassthroughSecret: string([]byte{}), // // TODO~ Need to remove it
		}
	}
	vpcBlockConfig := &vpcconfig.VPCBlockConfig{
		VPCConfig:    conf.VPC,
		IKSConfig:    conf.IKS,
		APIConfig:    conf.API,
		ServerConfig: conf.Server,
	}
	// Prepare provider registry
	registry, err := provider_util.InitProviders(vpcBlockConfig, logger)
	if err != nil {
		logger.Fatal("Error configuring providers", local.ZapError(err))
		Expect(err).To(HaveOccurred())
	}

	var providerName string
	if conf.IKS.Enabled {
		providerName = conf.IKS.IKSBlockProviderName
	} else if conf.VPC.Enabled {
		providerName = conf.VPC.VPCBlockProviderName
	}

	if conf.API == nil {
		conf.API = &config.APIConfig{
			PassthroughSecret: string([]byte{}), // // TODO~ Need to remove it
		}
	}

	sess, isFatal, err := provider_util.OpenProviderSessionWithContext(context.Background(), vpcBlockConfig, registry, providerName, logger)
	if err != nil || isFatal {
		logger.Error("Failed to get provider session", zap.Reflect("Error", err))
		Expect(err).To(HaveOccurred())
	}

	fmt.Println(sess)

})

var _ = AfterSuite(func() {
	defer sess.Close()
	defer ctxLogger.Sync()
})

func getContextLogger() (*zap.Logger, zap.AtomicLevel) {
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	traceLevel := zap.NewAtomicLevel()
	traceLevel.SetLevel(zap.InfoLevel)
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleDebugging, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return (lvl >= traceLevel.Level()) && (lvl < zapcore.ErrorLevel)
		})),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleErrors, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	)
	logger := zap.New(core, zap.AddCaller())
	return logger, traceLevel
}

func updateRequestID(err error, requestID string) error {
	if err == nil {
		return err
	}
	usrError, ok := err.(userError.Message)
	if !ok {
		return err
	}
	usrError.RequestID = requestID
	return usrError
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value + "-"
}
