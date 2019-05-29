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

func main() {
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

	//dc_name := "mex01"
	providerName := conf.Softlayer.SoftlayerBlockProviderName
	if conf.IKS != nil && conf.IKS.Enabled {
		providerName = conf.IKS.IKSBlockProviderName
	} else if conf.Softlayer.SoftlayerFileEnabled {
		providerName = conf.Softlayer.SoftlayerFileProviderName
	} else if conf.VPC.Enabled {
		providerName = conf.VPC.VPCBlockProviderName
	}

	valid := true
	for valid {

		fmt.Println("\n\nSelect your choice\n 1- Get volume details \n 2- Create snapshot \n 3- list snapshot \n 4- Create volume \n 5- Snapshot details \n 6- Snapshot Order \n 7- Create volume from snapshot\n 8- Delete volume \n 9- Delete Snapshot \n 10- List all Snapshot \n 12- Authorize volume \n 13- Create VPC Volume \n 14- Create VPC Snapshot \n 15- Attach VPC volume \n 16- Detach VPC volume \n 17- Get volume by name \n Your choice?:")

		var choiceN int
		var volumeID string
		var snapshotID string
		_, er11 := fmt.Scanf("%d", &choiceN)
		if er11 != nil {
			fmt.Printf("Wrong input, please provide option in int: ")
			fmt.Printf("\n\n")
			continue
		}
		ctxLogger, _ := getContextLogger()
		requestID := uid.NewV4().String()
		ctxLogger = ctxLogger.With(zap.String("RequestID", requestID))
		sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, providerName, ctxLogger)
		if err != nil {
			ctxLogger.Error("Failed to get session", zap.Reflect("Error", err))
			continue
		}
		defer sess.Close()
		defer ctxLogger.Sync()
		if choiceN == 1 {
			fmt.Println("You selected choice to get volume details")
			fmt.Printf("Please enter volume ID: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume, errr := sess.GetVolume(volumeID)
			if errr == nil {
				ctxLogger.Info("SUCCESSFULLY get volume details ================>", zap.Reflect("VolumeDetails", volume))
			} else {
				ctxLogger.Info("Provider error is ================>", zap.Reflect("ErrorType", userError.GetErrorType(errr)))
				errr = updateRequestID(errr, requestID)
				ctxLogger.Info("FAILED to get volume details ================>", zap.Reflect("VolumeID", volumeID), zap.Reflect("Error", errr))
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
				ctxLogger.Info("Successfully created snapshot on ================>", zap.Reflect("VolumeID", volumeID))
				ctxLogger.Info("Snapshot details: ", zap.Reflect("Snapshot", snapshot))
			} else {
				errr = updateRequestID(errr, requestID)
				ctxLogger.Info("Failed to create snapshot on ================>", zap.Reflect("VolumeID", volumeID), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 3 {
			fmt.Println("You selected choice to list snapshot from volume")
			fmt.Printf("Please enter volume ID to get the snapshots: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			fmt.Printf("\n")
			snapshots, errr := sess.ListAllSnapshots(volumeID)
			if errr == nil {
				ctxLogger.Info("Successfully get snapshot details ================>", zap.Reflect("Snapshot ID", volumeID))
				ctxLogger.Info("List of snapshots ", zap.Reflect("Snapshots are->", snapshots))
			} else {
				errr = updateRequestID(errr, requestID)
				ctxLogger.Info("Failed to get snapshot details ================>", zap.Reflect("Snapshot ID", volumeID), zap.Reflect("Error", errr))
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
				ctxLogger.Info("Successfully ordered volume ================>", zap.Reflect("volumeObj", volumeObj))
			} else {
				errr = updateRequestID(errr, requestID)
				ctxLogger.Info("Failed to order volume ================>", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 5 {
			fmt.Println("You selected choice to get snapshot details")
			fmt.Printf("Please enter Snapshot ID: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			snapdetails, errr := sess.GetSnapshot(volumeID)
			fmt.Printf("\n\n")
			if errr == nil {
				ctxLogger.Info("Successfully retrieved the snapshot details ================>", zap.Reflect("Snapshot ID", volumeID))
				ctxLogger.Info("Snapshot details ================>", zap.Reflect("SnapshotDetails", snapdetails))
			} else {
				ctxLogger.Info("Failed to get snapshot details ================>", zap.Reflect("Snapshot ID", volumeID), zap.Reflect("Error", errr))
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
				ctxLogger.Info("Successfully ordered snapshot space ================>", zap.Reflect("Volume ID", volumeID))
			} else {
				er11 = updateRequestID(er11, requestID)
				ctxLogger.Info("failed to order snapshot space================>", zap.Reflect("Volume ID", volumeID), zap.Reflect("Error", er11))
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
				ctxLogger.Info("Successfully Created volume from snapshot ================>", zap.Reflect("OriginalVolumeID", volumeID), zap.Reflect("SnapshotID", snapshotID))
				ctxLogger.Info("New volume from snapshot================>", zap.Reflect("New Volume->", vol))
			} else {
				errr = updateRequestID(errr, requestID)
				ctxLogger.Info("Failed to create volume from snapshot ================>", zap.Reflect("OriginalVolumeID", volumeID), zap.Reflect("SnapshotID", snapshotID), zap.Reflect("Error", errr))
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
				ctxLogger.Info("SUCCESSFULLY deleted volume ================>", zap.Reflect("Volume ID", volumeID))
			} else {
				er11 = updateRequestID(er11, requestID)
				ctxLogger.Info("FAILED volume deletion================>", zap.Reflect("Volume ID", volumeID), zap.Reflect("Error", er11))
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
				ctxLogger.Info("Successfully deleted snapshot ================>", zap.Reflect("Snapshot ID", snapshotID))
			} else {
				er11 = updateRequestID(er11, requestID)
				ctxLogger.Info("failed snapshot deletion================>", zap.Reflect("Snapshot ID", snapshotID), zap.Reflect("Error", er11))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 10 {
			fmt.Println("You selected choice to list all snapshot")
			list, errr := sess.ListSnapshots()
			if errr == nil {
				ctxLogger.Info("SUCCESSFULLY got all snapshots ================>", zap.Reflect("Snapshots", list))
			} else {
				errr = updateRequestID(errr, requestID)
				ctxLogger.Info("FAILED All snapshots ================>", zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 11 {
			fmt.Println("Get volume ID by using order ID")
			fmt.Printf("Please enter volume order ID to get volume ID:")
			_, er11 = fmt.Scanf("%s", &volumeID)
			_, error1 := sess.ListAllSnapshots(volumeID)
			if error1 != nil {
				error1 = updateRequestID(error1, requestID)
				ctxLogger.Info("Failed to get volumeID", zap.Reflect("Error", error1))
			}
		} else if choiceN == 12 {
			fmt.Println("Authorize volume")
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
				error1 = updateRequestID(error1, requestID)
				ctxLogger.Info("Failed to authorize", zap.Reflect("Error", error1))
			}
		} else if choiceN == 13 {
			fmt.Println("You selected choice to Create VPC volume")
			volume := &provider.Volume{}
			volumeName := ""
			volume.VolumeType = "vpc-block"

			resiurceGType := 0
			resourceGroup := "default resource group"
			zone := "us-south-1"
			volSize := 0
			Iops := "0"

			volume.Az = zone
			volume.VPCVolume.Generation = "gt"

			volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}

			profile := "general-purpose"
			fmt.Printf("\nPlease enter profile name supported profiles are [general-purpose, custom, 10iops-tier, 5iops-tier]: ")
			_, er11 = fmt.Scanf("%s", &profile)
			volume.VPCVolume.Profile = &provider.Profile{Name: profile}

			fmt.Printf("\nPlease enter volume name: ")
			_, er11 = fmt.Scanf("%s", &volumeName)
			volume.Name = &volumeName

			fmt.Printf("\nPlease enter volume size (Specify 10 GB - 2 TB of capacity in 1 GB increments): ")
			_, er11 = fmt.Scanf("%d", &volSize)
			volume.Capacity = &volSize

			fmt.Printf("\nPlease enter iops (Only custom profiles require iops): ")
			_, er11 = fmt.Scanf("%s", &Iops)
			volume.Iops = &Iops

			fmt.Printf("\nPlease enter resource group info type : 1- for ID and 2- for Name: ")
			_, er11 = fmt.Scanf("%d", &resiurceGType)
			if resiurceGType == 1 {
				fmt.Printf("\nPlease enter resource group ID:")
				_, er11 = fmt.Scanf("%s", &resourceGroup)
				volume.VPCVolume.ResourceGroup.ID = resourceGroup
			} else if resiurceGType == 2 {
				fmt.Printf("\nPlease enter resource group Name:")
				_, er11 = fmt.Scanf("%s", &resourceGroup)
				volume.VPCVolume.ResourceGroup.Name = resourceGroup
			} else {
				fmt.Printf("\nWrong resource group type\n")
				continue
			}

			fmt.Printf("\nPlease enter zone: ")
			_, er11 = fmt.Scanf("%s", &zone)
			volume.Az = zone

			volume.SnapshotSpace = &volSize
			volume.VPCVolume.Tags = []string{"Testing VPC Volume"}
			volumeObj, errr := sess.CreateVolume(*volume)
			if errr == nil {
				ctxLogger.Info("SUCCESSFULLY created volume...", zap.Reflect("volumeObj", volumeObj))
			} else {
				errr = updateRequestID(errr, requestID)
				ctxLogger.Info("FAILED to create volume...", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", errr))
			}
			fmt.Printf("\n\n")

		} else if choiceN == 14 {
			fmt.Println("You selected choice to order VPC snapshot")
			volume := &provider.Volume{}
			fmt.Printf("Please enter volume ID to create the snapshot space: ")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume.VolumeID = volumeID
			er11 := sess.OrderSnapshot(*volume)
			if er11 == nil {
				ctxLogger.Info("Successfully ordered snapshot space ================>", zap.Reflect("Volume ID", volumeID))
			} else {
				er11 = updateRequestID(er11, requestID)
				ctxLogger.Info("failed to order snapshot space================>", zap.Reflect("Volume ID", volumeID), zap.Reflect("Error", er11))
			}
			fmt.Printf("\n\n")
		} else if choiceN == 15 {
			fmt.Println("Enter the volume id to attach")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume := &provider.Volume{}
			volume.VolumeID = volumeID
			var instanceID string
			fmt.Println("Enter the instance id to attach")
			_, er11 = fmt.Scanf("%s", &instanceID)
			volumeAttachmentReq := provider.VolumeAttachmentRequest{
				VolumeID:   volumeID,
				InstanceID: instanceID,
				VPCVolumeAttachment: &provider.VolumeAttachment{
					DeleteVolumeOnInstanceDelete: false,
				},
			}
			response, err := sess.AttachVolume(volumeAttachmentReq)
			if err != nil {
				updateRequestID(err, requestID)
				ctxLogger.Error("Failed to attach the volume", zap.Error(err))
				return
			}
			volumeAttachmentReq.VPCVolumeAttachment = &provider.VolumeAttachment{
				ID: response.VPCVolumeAttachment.ID,
			}
			response, err = sess.WaitForAttachVolume(volumeAttachmentReq)
			if err != nil {
				updateRequestID(err, requestID)
				ctxLogger.Error("Failed to complete volume attach", zap.Error(err))
			}
			fmt.Println("Volume attachment", response, err)
		} else if choiceN == 16 {
			fmt.Println("Enter the volume id to detach")
			_, er11 = fmt.Scanf("%s", &volumeID)
			volume := &provider.Volume{}
			volume.VolumeID = volumeID
			var instanceID string
			fmt.Println("Enter the instance id to detach")
			_, er11 = fmt.Scanf("%s", &instanceID)
			volumeDetachmentReq := provider.VolumeAttachmentRequest{
				VolumeID:   volumeID,
				InstanceID: instanceID,
			}
			response, err := sess.DetachVolume(volumeDetachmentReq)
			if err != nil {
				updateRequestID(err, requestID)
				ctxLogger.Error("Failed to detach the volume", zap.Error(err))
				return
			}
			err = sess.WaitForDetachVolume(volumeDetachmentReq)
			if err != nil {
				updateRequestID(err, requestID)
				ctxLogger.Error("Failed to complete volume detach", zap.Error(err))
			}
			fmt.Println("Volume detach", response, err)
		} else if choiceN == 17 {
			fmt.Println("You selected get VPC volume by name")
			volumeName := ""
			fmt.Printf("Please enter volume Name to get the details: ")
			_, er11 = fmt.Scanf("%s", &volumeName)
			volumeobj1, er11 := sess.GetVolumeByName(volumeName)
			if er11 == nil {
				ctxLogger.Info("Successfully got VPC volume details ================>", zap.Reflect("VolumeDetail", volumeobj1))
			} else {
				er11 = updateRequestID(er11, requestID)
				ctxLogger.Info("failed to order snapshot space================>", zap.Reflect("VolumeName", volumeName), zap.Reflect("Error", er11))
			}
			fmt.Printf("\n\n")
		} else {
			fmt.Println("No right choice")
			return
		}
	}
}
