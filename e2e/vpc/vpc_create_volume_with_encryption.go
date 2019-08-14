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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
)

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	var (
		volume *provider.Volume
	)
	AfterEach(func() {
		sess.DeleteVolume(volume)
	})
	It("VPC: Create and delete VPC volume[with encryption]", func() {
		volName := volumeName + "with-encryption"
		volSize := volumeSize
		Iops := iops

		volume = &provider.Volume{}

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
		volume.VPCVolume.VolumeEncryptionKey = &provider.VolumeEncryptionKey{}
		volume.VPCVolume.VolumeEncryptionKey.CRN = volumeEncryptionKeyCRN

		volume.VPCVolume.Tags = []string{"Testing VPC volume from library with encryption"}
		By("Test Create Volume")
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
		By("Test Delete Volume")
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
