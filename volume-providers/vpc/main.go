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
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	//softlayer_block "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/block"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	//"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	//util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	//"github.com/IBM/ibmcloud-storage-volume-lib/provider/registry"
	provider_util "github.com/IBM/ibmcloud-storage-volume-lib/provider/utils"
)

func main() {

	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	), zap.AddCaller()).With(zap.String("Provider", "VPC"))

	defer logger.Sync()

	atom.SetLevel(zap.InfoLevel)

	// Load config file
	conf, err := config.ReadConfig("/root/gopath/src/github.com/IBM/ibmcloud-storage-volume-lib/etc/libconfig.toml", logger)
	if err != nil {
		logger.Fatal("Error loading configuration")
	}

	// Prepare provider registry
	providerRegistry, err := provider_util.InitProviders(conf, logger)
	if err != nil {
		logger.Fatal("Error configuring providers", local.ZapError(err))
	}

	providerName := conf.VPC.VPCBlockProviderName
	if conf.VPC.Enabled {
		providerName = conf.VPC.VPCBlockProviderName
	}
	providerTest, err := providerRegistry.Get("VPC")
	logger.Info("Provider registry testing", zap.Reflect("providerTest", providerTest))
	logger.Info("In main before openProviderSession call", zap.Reflect("providerRegistry", providerRegistry))

	sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
	if err != nil {
		logger.Error("Failed to get session", zap.Reflect("Error", err))
		return
	}
	logger.Info("In main after openProviderSession call", zap.Reflect("sess", sess))

	defer sess.Close()

	logger.Info("Currently you are using provider ....", zap.Reflect("ProviderName", sess.ProviderName()))

	volumeID := "80587c7b-ffc2-4aa8-9fea-438248ac1313"
	volume, errr := sess.GetVolume(volumeID)

	if errr == nil {
		logger.Info("Successfully got volume details==================>", zap.Reflect("Volume", volume))
	} else {
		logger.Info("Failed to get volume details ================>", zap.Reflect("VolumeID", volumeID), zap.Reflect("Error", errr))
	}
}
