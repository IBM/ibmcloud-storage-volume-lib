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

	//softlayer_block "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/block"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
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
	), zap.AddCaller()).With(zap.String("name", "ibm-volume-lib/main")).With(zap.String("VolumeLib", "IKS-VOLUME-LIB"))

	defer logger.Sync()

	atom.SetLevel(zap.InfoLevel)

	// Prepare main logger
	/*loggerLevel := zap.AtomicLevel()
	loggerLevel.SetLevel(zapcore.InfoLevel)
	logger := zap.New(
		zapcore.New.NewJSONEncoder(zap.RFC3339Formatter("ts")),
		zap.AddCaller(),
		loggerLevel,
	).With(zap.String("name", "ibm-volume-lib/main")).With(zap.String("VolumeLib", "IKS-VOLUME-LIB"))
	*/

	// Load config file
	conf, err := config.ReadConfig("", logger)
	if err != nil {
		logger.Fatal("Error loading configuration")
	}

	// Prepare provider registry
	providerRegistry, err := provider_util.InitProviders(conf, logger)
	if err != nil {
		logger.Fatal("Error configuring providers", local.ZapError(err))
	}

	//dc_name := "mex01"
	providerName := conf.Softlayer.SoftlayerBlockProviderName
	if conf.Softlayer.SoftlayerFileEnabled {
		providerName = conf.Softlayer.SoftlayerFileProviderName
	}
	logger.Info("In main before openProviderSession call", zap.Reflect("providerRegistry", providerRegistry))
	sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
	if err != nil {
		logger.Error("Failed to get session", zap.Reflect("Error", err))
		return
	}
	logger.Info("In main after openProviderSession call", zap.Reflect("sess", sess))
	defer sess.Close()
	logger.Info("Currently you are using provider ....", zap.Reflect("ProviderName", sess.ProviderName()))
	valid := true
	for valid {
		fmt.Println("\n\nSelect your choice\n 1- Get volume details \n 2- Create snapshot \n 3- list snapshot \n 4- Create volume \n 5- Snapshot details \n 6- Snapshot Order \n 7- Create volume from snapshot\n 8- Delete volume \n 9- Delete Snapshot \n 10- List all Snapshot \nYour choice?:")
		var choiceN int
		var volumeID string
		var snapshotID string
		_, er11 := fmt.Scanf("%d", &choiceN)
		if er11 != nil {
			fmt.Printf("Wrong input, please provide option in int: ")
			fmt.Printf("\n\n")
			continue
		}

		if choiceN == 1 {
			fmt.Println("You selected choice to get volume details")
			fmt.Printf("Please enter volume ID: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume, errr := sess.GetVolume(volumeID)
			if errr == nil {
				logger.Info("Successfully get volume details ================>", zap.Reflect("Volume ID", volumeID))
				logger.Info("Volume details are: ", zap.Reflect("Volume", volume))
			} else {
				logger.Info("Failed to get volume details ================>", zap.Reflect("VolumeID", volumeID), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 2 {
			fmt.Println("You selected choice to create snapshot")
			fmt.Printf("Please enter volume ID: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume := &provider.Volume{}
			volume.VolumeID = volumeID
			var tags map[string]string
			tags = make(map[string]string)
			tags["tag1"] = "snapshot-tag1"
			snapshot, errr := sess.CreateSnapshot(volume, tags)
			if errr == nil {
				logger.Info("Successfully created snapshot on ================>", zap.Reflect("VolumeID", volumeID))
				logger.Info("Snapshot details: ", zap.Reflect("Snapshot", snapshot))
			} else {
				logger.Info("Failed to create snapshot on ================>", zap.Reflect("VolumeID", volumeID), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 3 {
			fmt.Println("You selected choice to list snapshot from volume")
			fmt.Printf("Please enter volume ID to get the snapshots: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			fmt.Printf("\n")
			snapshots, errr := sess.ListAllSnapshots(volumeID)
			if errr == nil {
				logger.Info("Successfully get snapshot details ================>", zap.Reflect("Snapshot ID", volumeID))
				logger.Info("List of snapshots ", zap.Reflect("Snapshots are->", snapshots))
			} else {
				logger.Info("Failed to get snapshot details ================>", zap.Reflect("Snapshot ID", volumeID), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 4 {
			fmt.Println("You selected choice to Create volume")
			volume := &provider.Volume{}
			volume.VolumeType = "block"
			if conf.Softlayer.SoftlayerFileEnabled {
				volume.VolumeType = "file"
			}
			dcName := ""
			volSize := 0
			Iops := "0"
			tier := ""
			providerType := ""

			var choice int
			fmt.Printf("\nPlease enter storage type choice 1- for endurance  2- for performance: ")
			_, er11 = fmt.Scanf("%d", &choice)
			if choice == 1 {
				providerType = "endurance"
				volume.ProviderType = provider.VolumeProviderType(providerType)
			} else if choice == 2 {
				providerType = "performance"
				volume.ProviderType = provider.VolumeProviderType(providerType)
			}

			fmt.Printf("\nPlease enter datacenter name like dal09, dal10 or mex01  etc: ")
			_, er11 = fmt.Scanf("%s", &dcName)
			volume.Az = dcName

			fmt.Printf("\nPlease enter volume size in GB like 20, 40 80 etc : ")
			_, er11 = fmt.Scanf("%d", &volSize)
			volume.Capacity = &volSize

			if volume.ProviderType == "performance" {
				fmt.Printf("\nPlease enter iops from 1-48000 with multiple of 100: ")
				_, er11 = fmt.Scanf("%s", &Iops)
				volume.Iops = &Iops
			}
			if volume.ProviderType == "endurance" {
				fmt.Printf("\nPlease enter tier like 0.25, 2, 4, 10 iops per GB: ")
				_, er11 = fmt.Scanf("%s", &tier)
				volume.Tier = &tier
			}
			volume.SnapshotSpace = &volSize
			volume.VolumeNotes = map[string]string{"note": "test"}
			volumeObj, errr := sess.CreateVolume(*volume)
			if errr == nil {
				logger.Info("Successfully ordered volume ================>", zap.Reflect("volumeObj", volumeObj))
			} else {
				logger.Info("Failed to order volume ================>", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 5 {
			fmt.Println("You selected choice to get snapshot details")
			fmt.Printf("Please enter Snapshot ID: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			snapdetails, errr := sess.GetSnapshot(volumeID)
			fmt.Printf("\n\n")
			if errr == nil {
				logger.Info("Successfully get snapshot details ================>", zap.Reflect("Snapshot ID", volumeID))
				logger.Info("Snapshot details ================>", zap.Reflect("SnapshotDetails", snapdetails))
			} else {
				logger.Info("Failed to get snapshot details ================>", zap.Reflect("Snapshot ID", volumeID), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 6 {
			fmt.Println("You selected choice to order snapshot")
			volume := &provider.Volume{}
			fmt.Printf("Please enter volume ID to create the snapshot space: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume.VolumeID = volumeID
			var size int
			fmt.Printf("Please enter snapshot space size in GB: ")
			_, er11 = fmt.Scanf("%d", &size)
			volume.SnapshotSpace = &size
			er11 := sess.OrderSnapshot(*volume)
			if er11 == nil {
				logger.Info("Successfully ordered snapshot space ================>", zap.Reflect("Volume ID", volumeID))
			} else {
				logger.Info("failed to order snapshot space================>", zap.Reflect("Volume ID", volumeID), zap.Reflect("Error", er11))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 7 {
			fmt.Println("You selected choice to Create volume from snapshot")
			var snapshotVol provider.Snapshot
			var tags map[string]string
			fmt.Printf("Please enter original volume ID to create the volume from snapshot: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			fmt.Printf("Please enter snapshot ID for creating volume:")
			_, er11 = fmt.Scanf("%s", &snapshotID)
			snapshotVol.SnapshotID = snapshotID
			snapshotVol.Volume.VolumeID = volumeID
			vol, errr := sess.CreateVolumeFromSnapshot(snapshotVol, tags)
			if errr == nil {
				logger.Info("Successfully Created volume from snapshot ================>", zap.Reflect("OriginalVolumeID", volumeID), zap.Reflect("SnapshotID", snapshotID))
				logger.Info("New volume from snapshot================>", zap.Reflect("New Volume->", vol))
			} else {
				logger.Info("Failed to create volume from snapshot ================>", zap.Reflect("OriginalVolumeID", volumeID), zap.Reflect("SnapshotID", snapshotID), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 8 {
			fmt.Println("You selected choice to delete volume")
			volume := &provider.Volume{}
			fmt.Printf("Please enter volume ID for delete:")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume.VolumeID = volumeID
			er11 = sess.DeleteVolume(volume)
			if er11 == nil {
				logger.Info("Successfully deleted volume ================>", zap.Reflect("Volume ID", volumeID))
			} else {
				logger.Info("failed volume deletion================>", zap.Reflect("Volume ID", volumeID), zap.Reflect("Error", er11))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 9 {
			fmt.Println("You selected choice to delete snapshot")
			snapshot := &provider.Snapshot{}
			fmt.Printf("Please enter snapshot ID for delete:")
			_, er11 = fmt.Scanf("%s", &snapshotID)
			snapshot.SnapshotID = snapshotID
			er11 = sess.DeleteSnapshot(snapshot)
			if er11 == nil {
				logger.Info("Successfully deleted snapshot ================>", zap.Reflect("Snapshot ID", snapshotID))
			} else {
				logger.Info("failed snapshot deletion================>", zap.Reflect("Snapshot ID", snapshotID), zap.Reflect("Error", er11))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 10 {
			fmt.Println("You selected choice to list all snapshot")
			list, _ := sess.ListSnapshots()
			logger.Info("All snapshots ================>", zap.Reflect("Snapshots", list))
			fmt.Printf("\n\n")
		} else if choiceN == 11 {
			fmt.Println("Get volume ID by using order ID")
			fmt.Printf("Please enter volume order ID to get volume ID:")
			_, er11 = fmt.Scanf("%s", &volumeID)
			_, error1 := sess.ListAllSnapshots(volumeID)
			if error1 != nil {
				logger.Info("Failed to get volumeID", zap.Reflect("Error", error1))
			}
		} else {
			fmt.Println("No right choice")
			return
		}
	}
}
