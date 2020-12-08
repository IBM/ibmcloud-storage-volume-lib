/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package main

import (
	"fmt"
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

func updateRequestID1(err error, requestID string) error {
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

func main_test() {
	// Setup new style zap logger
	logger, traceLevel := getContextLogger()
	defer logger.Sync()

	// Load config file
	goPath := os.Getenv("GOPATH")
	conf, err := config.ReadConfig(goPath+"/src/github.com/IBM/ibmcloud-storage-volume-lib/etc/libconfig.toml", logger)
	if err != nil {
		logger.Fatal("Error loading configuration")
	}

	// Check if debug log level enabled or not
	if conf.Server != nil && conf.Server.DebugTrace {
		traceLevel.SetLevel(zap.DebugLevel)
	}

	// Prepare provider registry
	providerRegistry, err := provider_util.InitProviders(conf, logger)
	if err != nil {
		logger.Fatal("Error configuring providers", local.ZapError(err))
	}

	providerName := ""
	if conf.VPC.Enabled {
		providerName = conf.VPC.VPCBlockProviderName
	}

	ctxLogger, _ := getContextLogger()
	requestID := uid.NewV4().String()
	ctxLogger = ctxLogger.With(zap.String("RequestID", requestID))
	sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, ctxLogger)
	if err != nil {
		ctxLogger.Error("Failed to get session", zap.Reflect("Error", err))
		os.Exit(3)
	}
	defer sess.Close()
	defer ctxLogger.Sync()

	volume := &provider.Volume{}
	volumeName := ""
	volume.VolumeType = "vpc-block"

	resiurceGType := 0
	resourceGroup := "default resource group"
	zone := "us-south-1"
	volSize := 0
	Iops := "0"

	volume.Az = zone

	volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}

	profile := "general-purpose"
	volume.VPCVolume.Profile = &provider.Profile{Name: profile}

	fmt.Printf("\nPlease enter volume name: ")
	_, err = fmt.Scanf("%s", &volumeName)
	volume.Name = &volumeName

	fmt.Printf("\nPlease enter volume size (Specify 10 GB - 2 TB of capacity in 1 GB increments): ")
	_, err = fmt.Scanf("%d", &volSize)
	volume.Capacity = &volSize

	fmt.Printf("\nPlease enter iops (Only custom profiles require iops): ")
	_, err = fmt.Scanf("%s", &Iops)
	volume.Iops = &Iops

	fmt.Printf("\nPlease enter resource group info type : 1- for ID and 2- for Name: ")
	_, err = fmt.Scanf("%d", &resiurceGType)
	if resiurceGType == 1 {
		fmt.Printf("\nPlease enter resource group ID:")
		_, err = fmt.Scanf("%s", &resourceGroup)
		volume.VPCVolume.ResourceGroup.ID = resourceGroup
	} else if resiurceGType == 2 {
		fmt.Printf("\nPlease enter resource group Name:")
		_, err = fmt.Scanf("%s", &resourceGroup)
		volume.VPCVolume.ResourceGroup.Name = resourceGroup
	} else {
		fmt.Printf("\nWrong resource group type\n")
		os.Exit(3)
	}

	fmt.Printf("\nPlease enter zone: ")
	_, err = fmt.Scanf("%s", &zone)
	volume.Az = zone

	volume.VPCVolume.VolumeEncryptionKey = &provider.VolumeEncryptionKey{}
	fmt.Printf("\nPlease enter encryption key CRN:")
	volumeEncryptionKeyCRN := ""
	_, err = fmt.Scanf("%s", &volumeEncryptionKeyCRN)
	volume.VPCVolume.VolumeEncryptionKey.CRN = volumeEncryptionKeyCRN

	volume.SnapshotSpace = &volSize
	volume.VPCVolume.Tags = []string{"Testing VPC Volume"}
	volumeObj, err := sess.CreateVolume(*volume)
	if err == nil {
		ctxLogger.Info("SUCCESSFULLY created volume...", zap.Reflect("volumeObj", volumeObj))
	} else {
		err = updateRequestID(err, requestID)
		ctxLogger.Info("FAILED to create volume...", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", err))
	}
	fmt.Printf("\n\n")

}
