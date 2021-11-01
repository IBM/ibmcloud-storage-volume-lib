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
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"

	userError "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	file_provider_util "github.com/IBM/ibmcloud-volume-file-vpc/file/utils"
	filevpcconfig "github.com/IBM/ibmcloud-volume-file-vpc/file/vpcconfig"
	"github.com/IBM/ibmcloud-volume-interface/config"
	"github.com/IBM/ibmcloud-volume-interface/lib/provider"
	provider_util "github.com/IBM/ibmcloud-volume-vpc/block/utils"
	vpcconfig "github.com/IBM/ibmcloud-volume-vpc/block/vpcconfig"
	"github.com/IBM/ibmcloud-volume-vpc/common/registry"
	uid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var sess provider.Session
var vpcBlockConfig *vpcconfig.VPCBlockConfig
var vpcFileConfig *filevpcconfig.VPCFileConfig
var providerRegistry registry.Providers
var logger *zap.Logger
var ctxLogger *zap.Logger
var traceLevel zap.AtomicLevel
var requestID string
var resourceGroupID string
var vpcZone string
var volumeEncryptionKeyCRN string
var startTime time.Time
var providerName string
var conf *config.Config
var err error

var testCaseList []TestCaseData
var goPath string
var content []byte

type TestCaseData struct {
	TestCase string   `yaml:"testCase"`
	Input    InputDef `yaml:"input"`
	Success  bool     `yaml:"success,omitempty"`
	Skip     bool     `yaml:"skip,omitempty"`
}

type InputDef struct {
	Volume            VolumeDef `yaml:"volume,flow"`
	EncryptionEnabled bool      `yaml:"encryptionEnabled,omitempty"`
	VPCZone           string    `yaml:"vpcZone,omitempty"`
	InstanceID        []string  `yaml:"instanceID,omitempty"`
	VPCID             []string  `yaml:"vpcID,omitempty"`
	SubnetID          []string  `yaml:"subnetID,omitempty"`
	ClusterID         []string  `yaml:"clusterID,omitempty"`
}

type VolumeDef struct {
	Name         string `yaml:"name"`
	Capacity     int    `yaml:"capacity"`
	Profile      string `yaml:"profile"`
	Iops         string `yaml:"iops,omitempty"`
	Tags         string `yaml:"tags,omitempty"`
	InitialOwner bool   `yaml:"initialOwner,omitempty"`
}

func TestVPCE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ibmcloud-storage-volume-lib VPC e2e test suite")
}

var _ = BeforeSuite(func() {
})

func initializeTestCaseData() {
	loadConfig()

	err = yaml.Unmarshal(content, &testCaseList)
	if err != nil {
		ctxLogger.Fatal("Error parsing test case data file for block")
		Expect(err).To(HaveOccurred())
		return
	}

}

func loadConfig() {
	// Setup new style zap logger
	logger, traceLevel = getContextLogger()
	defer logger.Sync()
	// Load config file
	goPath = os.Getenv("GOPATH")
	conf, err = config.ReadConfig(goPath+vpcConfigFilePath, logger)
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

	ctxLogger, _ = getContextLogger()
	requestID = uid.NewV4().String()
	ctxLogger = logger.With(zap.String("RequestID", requestID))

	initializeProvider()
}

func initializeProvider() {
	if conf.VPC.VPCVolumeType == "vpc-share" {
		conf.VPC.G2APIVersion = fileAPIVersion
		initVPCFileProvider()
		content, err = ioutil.ReadFile(goPath + testCasesForFile)

		if err != nil {
			logger.Fatal("Error loading test case data for file")
			Expect(err).To(HaveOccurred())
			return
		}

	} else if conf.VPC.VPCVolumeType == "vpc-block" {
		initVPCBlockProvider()
		content, err = ioutil.ReadFile(goPath + testCasesForBlock)
		if err != nil {
			logger.Fatal("Error loading test case data for block")
			Expect(err).To(HaveOccurred())
			return
		}

	} else {
		logger.Error("Error configuring providers. Please set valid volume type for VPC", local.ZapError(err))
		Expect(err).To(HaveOccurred())
	}
}

func initVPCBlockProvider() {

	vpcBlockConfig = &vpcconfig.VPCBlockConfig{
		VPCConfig:    conf.VPC,
		IKSConfig:    conf.IKS,
		APIConfig:    conf.API,
		ServerConfig: conf.Server,
	}

	ctxLogger, _ = getContextLogger()
	requestID = uid.NewV4().String()
	ctxLogger = logger.With(zap.String("RequestID", requestID))

	// Prepare provider registry
	providerRegistry, err = provider_util.InitProviders(vpcBlockConfig, ctxLogger)
	if err != nil {
		logger.Fatal("Error configuring providers", local.ZapError(err))
		Expect(err).To(HaveOccurred())
	}

	if conf.IKS.Enabled {
		providerName = conf.IKS.IKSBlockProviderName
	} else {
		providerName = conf.VPC.VPCBlockProviderName
	}
}

func initVPCFileProvider() {
	vpcFileConfig = &filevpcconfig.VPCFileConfig{
		VPCConfig:    conf.VPC,
		ServerConfig: conf.Server,
	}

	// Prepare provider registry
	providerRegistry, err = file_provider_util.InitProviders(vpcFileConfig, logger)
	if err != nil {
		logger.Error("Error configuring providers", local.ZapError(err))
		Expect(err).To(HaveOccurred())
	}

	providerName = conf.VPC.VPCVolumeType
}

var _ = AfterSuite(func() {
	defer sess.Close()
	defer ctxLogger.Sync()
})

func RefreshSession() {
	var err error

	if sess != nil {
		sess.Close()
	}

	if conf.VPC.VPCVolumeType == "vpc-share" {
		openVPCFileSession()
	} else if conf.VPC.VPCVolumeType == "vpc-block" {
		openVPCBlockSession()
	} else {
		logger.Error("Error configuring providers. Please set valid volume type for VPC", local.ZapError(err))
		Expect(err).To(HaveOccurred())
	}
}

func openVPCFileSession() {
	sess, _, err = file_provider_util.OpenProviderSession(vpcFileConfig, providerRegistry, providerName, logger)
	if err != nil {
		logger.Error("Failed to get provider session", zap.Reflect("Error", err))
		Expect(err).To(HaveOccurred())
	}
}

func openVPCBlockSession() {
	sess, _, err = provider_util.OpenProviderSession(vpcBlockConfig, providerRegistry, providerName, logger)
	if err != nil {
		logger.Error("Failed to get provider session", zap.Reflect("Error", err))
		Expect(err).To(HaveOccurred())
	}
}

func CloseSession() {
	if sess != nil {
		sess.Close()
	}
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
