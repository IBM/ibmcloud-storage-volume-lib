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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	provider_util "github.com/IBM/ibmcloud-storage-volume-lib/provider/utils"
	uid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
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

	if conf.VPC != nil && conf.VPC.VPCTypeEnabled == "g2" && conf.VPC.G2ResourceGroupID != "" {
		resourceGroupID = conf.VPC.G2ResourceGroupID
	} else if conf.VPC != nil && conf.VPC.ResourceGroupID != "" {
		resourceGroupID = conf.VPC.ResourceGroupID
	}

	// Check if debug log level enabled or not
	if conf.Server != nil && conf.Server.DebugTrace {
		traceLevel.SetLevel(zap.DebugLevel)
	}

	// Load zone info
	vpcZone = os.Getenv("VPC_ZONE")

	// Load EncryptionKeyCRN info
	volumeEncryptionKeyCRN = os.Getenv("ENCRYPTION_KEY_CRN")

	// Prepare provider registry
	providerRegistry, err := provider_util.InitProviders(conf, logger)
	if err != nil {
		logger.Fatal("Error configuring providers", local.ZapError(err))
		Expect(err).To(HaveOccurred())
	}

	providerName := ""
	if conf.VPC.Enabled {
		providerName = conf.VPC.VPCBlockProviderName
	} else {
		providerName = conf.IKS.IKSBlockProviderName
	}

	ctxLogger, _ = getContextLogger()
	requestID = uid.NewV4().String()
	ctxLogger = logger.With(zap.String("RequestID", requestID))
	sess, _, err = provider_util.OpenProviderSession(conf, providerRegistry, providerName, ctxLogger)
	if err != nil {
		Expect(err).To(HaveOccurred())
	}

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
