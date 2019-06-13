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
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	provider_util "github.com/IBM/ibmcloud-storage-volume-lib/provider/utils"
	uid "github.com/satori/go.uuid"
)

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	It("Create and delete VPC volume", func() {
		var sess provider.Session
		var logger *zap.Logger
		var ctxLogger *zap.Logger
		var traceLevel zap.AtomicLevel
		var requestID string

		volName := volumeName
		volSize := volumeSize
		Iops := iops

		// Setup new style zap logger
		logger, traceLevel = getContextLogger()
		defer logger.Sync()
		// Load config file
		goPath := os.Getenv("GOPATH")
		conf, err := config.ReadConfig(goPath+"/src/github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/e2e/config/config.toml", logger)
		if err != nil {
			logger.Fatal("Error loading configuration")
			Expect(err).To(HaveOccurred())
		}

		// Check if debug log level enabled or not
		if conf.Server != nil && conf.Server.DebugTrace {
			traceLevel.SetLevel(zap.DebugLevel)
		}

		// Prepare provider registry
		providerRegistry, err := provider_util.InitProviders(conf, logger)
		if err != nil {
			logger.Fatal("Error configuring providers", local.ZapError(err))
			Expect(err).To(HaveOccurred())
		}

		providerName := ""
		if conf.VPC.Enabled {
			providerName = conf.VPC.VPCBlockProviderName
		}

		ctxLogger, _ = getContextLogger()
		requestID = uid.NewV4().String()
		ctxLogger = logger.With(zap.String("RequestID", requestID))
		sess, _, err = provider_util.OpenProviderSession(conf, providerRegistry, providerName, ctxLogger)
		if err != nil {
			Expect(err).To(HaveOccurred())
		}
		defer sess.Close()
		defer ctxLogger.Sync()
		volume := &provider.Volume{}

		volume.VolumeType = volumeType
		volume.VPCVolume.Generation = generation
		volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}
		profile := vpcProfile
		volume.VPCVolume.Profile = &provider.Profile{Name: profile}
		volume.Name = &volName
		volume.Capacity = &volSize
		volume.Iops = &Iops
		volume.VPCVolume.ResourceGroup.ID = resourceGroupID
		volume.Az = vpcZone

		volume.VPCVolume.Tags = []string{"Testing VPC Volume"}
		volumeObj, err := sess.CreateVolume(*volume)
		if err == nil {
			Expect(err).NotTo(HaveOccurred())
			ctxLogger.Info("Successfully created volume...", zap.Reflect("volumeObj", volumeObj))
		} else {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to create volume...", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", err))
			Expect(err).To(HaveOccurred())
		}
		fmt.Printf("\n\n")

		volume = &provider.Volume{}
		volume.VolumeID = volumeObj.VolumeID
		err = sess.DeleteVolume(volume)
		if err == nil {
			Expect(err).NotTo(HaveOccurred())
			ctxLogger.Info("Successfully deleted volume...", zap.Reflect("volumeObj", volume))
		} else {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to delete volume...", zap.Reflect("StorageType", volume.VolumeID), zap.Reflect("Error", err))
			Expect(err).To(HaveOccurred())
		}
		fmt.Printf("\n\n")
	})
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
