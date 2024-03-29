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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/IBM/ibmcloud-volume-interface/lib/provider"
)

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	var (
		volume *provider.Volume
	)
	BeforeEach(func() {
		RefreshSession()
	})

	AfterEach(func() {
		sess.DeleteVolume(volume)
		CloseSession()
	})
	It("VPC: Create and delete VPC volume[profile: custom]", func() {
		volName := volumeName + "-custom-vol"
		volSize := volumeSize
		Iops := iops

		volume = &provider.Volume{}

		volume.VolumeType = volumeType
		volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}
		volume.VPCVolume.Profile = &provider.Profile{Name: "custom"}
		volume.Name = &volName
		volume.Capacity = &volSize
		volume.Iops = &Iops
		volume.VPCVolume.ResourceGroup.ID = resourceGroupID
		volume.Az = vpcZone

		volume.VPCVolume.Tags = []string{"Testing VPC volume from library with custom profile"}
		By("Test Create Volume")
		startTime = time.Now()
		volumeObj, err := sess.CreateVolume(*volume)
		if err == nil {
			Expect(err).NotTo(HaveOccurred())
			ctxLogger.Info("Successfully created volume...", zap.Reflect("volumeObj", volumeObj))
		} else {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to create volume...", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", err))
			Expect(err).To(HaveOccurred())
		}
		ctxLogger.Info("Test Create Volume", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
		fmt.Printf("\n\n")

		volume = &provider.Volume{}
		volume.VolumeID = volumeObj.VolumeID
		By("Test Delete Volume")
		startTime = time.Now()
		err = sess.DeleteVolume(volume)
		if err == nil {
			Expect(err).NotTo(HaveOccurred())
			ctxLogger.Info("Successfully deleted volume...", zap.Reflect("volumeObj", volume))
		} else {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to delete volume...", zap.Reflect("StorageType", volume.VolumeID), zap.Reflect("Error", err))
			Expect(err).To(HaveOccurred())
		}
		ctxLogger.Info("Test Delete Volume", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
		fmt.Printf("\n\n")
	})
})
