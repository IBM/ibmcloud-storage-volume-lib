/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/
package e2e

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/registry"
	provider_util "github.com/IBM/ibmcloud-storage-volume-lib/provider/utils"
	uid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sess provider.Session
var logger *zap.Logger
var ctxLogger *zap.Logger
var traceLevel zap.AtomicLevel
var requestID string
var providerRegistry registry.Providers
var providers []ProviderE2ETest

func TestVPCE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ibmcloud-storage-volume-lib VPC e2e test suite")
}

var _ = BeforeSuite(func() {
	//initSuite()
})

var _ = AfterSuite(func() {
	//defer sess.Close()
	//defer ctxLogger.Sync()
})

func initSuite() {
	// Setup new style zap logger
	logger, traceLevel = getContextLogger()
	defer logger.Sync()
	// Load config file
	conf, err := config.ReadConfig("", logger)
	if err != nil {
		logger.Fatal("Error loading configuration")
		Expect(err).To(HaveOccurred())
	}

	// Check if debug log level enabled or not
	if conf.Server != nil && conf.Server.DebugTrace {
		traceLevel.SetLevel(zap.DebugLevel)
	}
	logger.Info("Config", zap.Reflect("Config", conf))

	// Prepare provider registry
	providerRegistry, err = provider_util.InitProviders(conf, logger)
	if err != nil {
		logger.Fatal("Error configuring providers", local.ZapError(err))
		Expect(err).To(HaveOccurred())
	}

	populateE2EProviders(conf)
	logger.Info("No. of providers enabled", zap.Int("No.of Providers", len(providers)))

}

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

func populateE2EProviders(conf *config.Config) {
	providerName := ""
	ctxLogger, _ = getContextLogger()
	requestID = uid.NewV4().String()
	ctxLogger = logger.With(zap.String("RequestID", requestID))
	if conf.VPC.Enabled {
		providerName = conf.VPC.VPCBlockProviderName
		sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, ctxLogger)
		if err != nil {
			Expect(err).To(HaveOccurred())
		}
		e2eProvider := &VpcClassicE2E{
			BaseE2ETest{
				session: sess,
				Name:    providerName,
			},
		}
		providers = append(providers, e2eProvider)
	}

	if conf.Softlayer.SoftlayerBlockEnabled {
		providerName = conf.Softlayer.SoftlayerBlockProviderName
		sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, ctxLogger)
		if err != nil {
			Expect(err).To(HaveOccurred())
		}
		e2eProvider := &SLBlockE2E{
			BaseE2ETest{
				session: sess,
				Name:    providerName,
			},
		}
		providers = append(providers, e2eProvider)
	}

}
