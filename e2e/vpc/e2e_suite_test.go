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
	"os"
	"time"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	provider_util "github.com/IBM/ibmcloud-storage-volume-lib/provider/utils"
	uid "github.com/satori/go.uuid"
	"go.uber.org/zap"
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
	if conf.VPC != nil && conf.VPC.ResourceGroupID != "" {
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
