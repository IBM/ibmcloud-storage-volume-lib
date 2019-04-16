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
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	provider_util "github.com/IBM/ibmcloud-storage-volume-lib/provider/utils"
	uid "github.com/satori/go.uuid"
)

func main() {
	// Setup new style zap logger
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

	//dc_name := "mex01"
	providerName := conf.Softlayer.SoftlayerBlockProviderName
	if conf.Softlayer.SoftlayerFileEnabled {
		providerName = conf.Softlayer.SoftlayerFileProviderName
	} else if conf.VPC.Enabled {
		providerName = conf.VPC.VPCBlockProviderName
	}

	valid := true
	for valid {
		fmt.Println("\n\nSelect your choice\n 1- Get volume details \n 2- Create snapshot \n 3- list snapshot \n 4- Create volume \n 5- Snapshot details \n 6- Snapshot Order \n 7- Create volume from snapshot\n 8- Delete volume \n 9- Delete Snapshot \n 10- List all Snapshot \n 12- Authorize volume \n 13- Create VPC Volume \n 14- Create VPC Snapshot \n Your choice?:")
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
			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				continue
			}
			defer sess.Close()

			fmt.Printf("Please enter volume ID: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume, errr := sess.GetVolume(volumeID)
			if errr == nil {
				logger.Info("SUCCESSFULLY get volume details ================>", zap.Reflect("VolumeDetails", volume))
			} else {
				logger.Info("FAILED to get volume details ================>", zap.Reflect("VolumeID", volumeID), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 2 {
			fmt.Println("You selected choice to create snapshot")
			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

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

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

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

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

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

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

			fmt.Printf("Please enter Snapshot ID: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			snapdetails, errr := sess.GetSnapshot(volumeID)
			fmt.Printf("\n\n")
			if errr == nil {
				logger.Info("Successfully retrieved the snapshot details ================>", zap.Reflect("Snapshot ID", volumeID))
				logger.Info("Snapshot details ================>", zap.Reflect("SnapshotDetails", snapdetails))
			} else {
				logger.Info("Failed to get snapshot details ================>", zap.Reflect("Snapshot ID", volumeID), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 6 {
			fmt.Println("You selected choice to order snapshot")

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

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

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

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

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

			volume := &provider.Volume{}
			fmt.Printf("Please enter volume ID for delete:")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume.VolumeID = volumeID
			er11 = sess.DeleteVolume(volume)
			if er11 == nil {
				logger.Info("SUCCESSFULLY deleted volume ================>", zap.Reflect("Volume ID", volumeID))
			} else {
				logger.Info("FAILED volume deletion================>", zap.Reflect("Volume ID", volumeID), zap.Reflect("Error", er11))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 9 {
			fmt.Println("You selected choice to delete snapshot")

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

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

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

			list, _ := sess.ListSnapshots()
			logger.Info("All snapshots ================>", zap.Reflect("Snapshots", list))
			fmt.Printf("\n\n")
		} else if choiceN == 11 {
			fmt.Println("Get volume ID by using order ID")

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

			fmt.Printf("Please enter volume order ID to get volume ID:")
			_, er11 = fmt.Scanf("%s", &volumeID)
			_, error1 := sess.ListAllSnapshots(volumeID)
			if error1 != nil {
				logger.Info("Failed to get volumeID", zap.Reflect("Error", error1))
			}
		} else if choiceN == 12 {
			fmt.Println("Authorize volume")

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

			fmt.Printf("Please enter volume ID:")
			_, er11 = fmt.Scanf("%s", &volumeID)
			var subnetIDs string
			fmt.Printf("Please enter subnet IDs comma seperated, default[]")
			_, er11 = fmt.Scanf("%s", &subnetIDs)
			var hostIPs string
			fmt.Printf("Please enter host IPs comma seperated, default[]")
			_, er11 = fmt.Scanf("%s", &hostIPs)
			splitFn := func(c rune) bool {
				return c == ','
			}
			subnetIDList := strings.FieldsFunc(subnetIDs, splitFn)
			hostIPList := strings.FieldsFunc(strings.TrimSpace(hostIPs), splitFn)
			fmt.Printf("lengnt:%d", len(hostIPList))
			volumeObj, _ := sess.GetVolume(volumeID)
			authRequest := provider.VolumeAuthorization{
				Volume:  *volumeObj,
				Subnets: subnetIDList,
				HostIPs: hostIPList,
			}
			error1 := sess.AuthorizeVolume(authRequest)
			if error1 != nil {
				logger.Info("Failed to authorize", zap.Reflect("Error", error1))
			}
		} else if choiceN == 13 {
			fmt.Println("You selected choice to Create VPC volume")

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

			volume := &provider.Volume{}
			volumeName := ""
			volume.VolumeType = "vpc-block"

			resourceGroup := "default resource group"
			zone := "us-south-1"
			volSize := 0
			Iops := "0"

			volume.Az = zone
			volume.VPCVolume.Generation = "gt"

			volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{ID: resourceGroup}

			volume.VPCVolume.Profile = &provider.Profile{Name: "general-purpose"}

			fmt.Printf("\nPlease enter volume name: ")
			_, er11 = fmt.Scanf("%s", &volumeName)
			volume.Name = &volumeName

			fmt.Printf("\nPlease enter volume size (Specify 10 GB - 2 TB of capacity in 1 GB increments): ")
			_, er11 = fmt.Scanf("%d", &volSize)
			volume.Capacity = &volSize

			fmt.Printf("\nPlease enter iops (For general purpose profiles, leave it empty): ")
			_, er11 = fmt.Scanf("%s", &Iops)
			volume.Iops = &Iops

			fmt.Printf("\nPlease enter resource group: ")
			_, er11 = fmt.Scanf("%s", &resourceGroup)
			volume.VPCVolume.ResourceGroup.ID = resourceGroup

			fmt.Printf("\nPlease enter zone: ")
			_, er11 = fmt.Scanf("%s", &zone)
			volume.Az = zone

			volume.SnapshotSpace = &volSize
			volume.VPCVolume.Tags = []string{"Testing VPC Volume"}
			volumeObj, errr := sess.CreateVolume(*volume)
			if errr == nil {
				logger.Info("SUCCESSFULLY created volume...", zap.Reflect("volumeObj", volumeObj))
			} else {
				logger.Info("FAILED to create volume...", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")

		} else if choiceN == 14 {
			fmt.Println("You selected choice to order VPC snapshot")

			logger = logger.With(zap.String("RequestID", uid.NewV4().String()))
			sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, logger)
			if err != nil {
				logger.Error("Failed to get session", zap.Reflect("Error", err))
				return
			}
			defer sess.Close()

			volume := &provider.Volume{}
			fmt.Printf("Please enter volume ID to create the snapshot space: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume.VolumeID = volumeID
			er11 := sess.OrderSnapshot(*volume)
			if er11 == nil {
				logger.Info("Successfully ordered snapshot space ================>", zap.Reflect("Volume ID", volumeID))
			} else {
				logger.Info("failed to order snapshot space================>", zap.Reflect("Volume ID", volumeID), zap.Reflect("Error", er11))
			}
			fmt.Printf("\n\n")
		} else {
			fmt.Println("No right choice")
			return
		}
	}
}
