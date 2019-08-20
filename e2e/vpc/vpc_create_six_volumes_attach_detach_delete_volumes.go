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
	"os"
	"strconv"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
)

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	var (
		volume                  *provider.Volume
		numberOfVolumesRequired = 6
	)
	It("VPC: Parallel volume attachments, detachments [six volume attachments in parallel]", func() {
		By("Creating test volumes")
		volumes, err := CreateTestVolumes(numberOfVolumesRequired)
		if err == nil {
			ctxLogger.Info("Successfully created the test volumes...")
		} else {
			ctxLogger.Info("Failed to create test volumes.", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", err))
		}

		By("Testing parallel attach volumes")
		attachmentRequests, attachmentResponses, err := CreateVolumeAttachments(volumes)
		for i := 0; i < len(attachmentResponses); i++ {
			sess.WaitForAttachVolume(*attachmentRequests[i])
		}
		if err == nil {
			ctxLogger.Info("Successfully attached the volumes...")
		}

		By("Test parallel detach Volumes")
		for i := 0; i < len(attachmentResponses); i++ {
			httpResponse, err := sess.DetachVolume(*attachmentRequests[i])
			ctxLogger.Info("Successfully detached the volume.", zap.Reflect("httpResponse", httpResponse))
			Expect(err).NotTo(HaveOccurred())
		}
		for i := 0; i < len(attachmentResponses); i++ {
			sess.WaitForDetachVolume(*attachmentRequests[i])
		}
		ctxLogger.Info("Successfully detached the volumes.")

		By("Test Delete Volume")
		err = DeleteTestVolumes(volumes)
		if err == nil {
			ctxLogger.Info("Successfully deleted all the test volumes.")
		} else {
			ctxLogger.Info("Failed to delete volumes.", zap.Reflect("Error", err))
		}
		fmt.Printf("\n\n")
	})
})

func CreateTestVolumes(numberOfVolumesRequired int) ([]*provider.Volume, error) {
	var (
		volName string
		volSize int
		Iops    string
		err     error
		volume  *provider.Volume
	)
	var volumes = make([]*provider.Volume, numberOfVolumesRequired)
	for i := 0; i < numberOfVolumesRequired; i++ {
		volName = volumeName + "-attach-detach-" + strconv.Itoa(i+1)
		volSize = volumeSize
		Iops = iops
		volume = &provider.Volume{}

		volume.VolumeType = volumeType
		volume.VPCVolume.Generation = generation
		volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}
		volume.VPCVolume.Profile = &provider.Profile{Name: "5iops-tier"}
		volume.Name = &volName
		volume.Capacity = &volSize
		volume.Iops = &Iops
		volume.VPCVolume.ResourceGroup.ID = resourceGroupID
		volume.Az = vpcZone
		volume.VPCVolume.Tags = []string{"Testing VPC create volume, attach volume, detach volume, and delete volume"}
		volumeObj, err := sess.CreateVolume(*volume)
		volumes[i] = volumeObj
		if err == nil {
			Expect(err).NotTo(HaveOccurred())
			ctxLogger.Info("Successfully created volume...", zap.Reflect("volumeObj", volumes[i]))
		} else {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to create volume...", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", err))
			Expect(err).To(HaveOccurred())
		}
		fmt.Printf("\n\n")
	}
	return volumes, err
}

func CreateVolumeAttachments(volumes []*provider.Volume) ([]*provider.VolumeAttachmentRequest, []*provider.VolumeAttachmentResponse, error) {
	var volumeAttachRequests = make([]*provider.VolumeAttachmentRequest, len(volumes))
	var attachmentResponses = make([]*provider.VolumeAttachmentResponse, len(volumes))
	var err error
	for i := 0; i < len(volumes); i++ {
		volumeAttachRequests[i] = &provider.VolumeAttachmentRequest{}
		volumeAttachRequests[i].VolumeID = volumes[i].VolumeID
		volumeAttachRequests[i].InstanceID = os.Getenv("INSTANCE_ID")

		attachResponse, err := sess.AttachVolume(*volumeAttachRequests[i])
		attachmentResponses[i] = attachResponse
		Expect(err).NotTo(HaveOccurred())
	}
	return volumeAttachRequests, attachmentResponses, err
}

func DeleteTestVolumes(volumes []*provider.Volume) (err error) {
	for i := 0; i < len(volumes); i++ {
		err = sess.DeleteVolume(volumes[i])
		if err == nil {
			Expect(err).NotTo(HaveOccurred())
			ctxLogger.Info("Successfully deleted volume.", zap.Reflect("volume", volumes[i]))
		} else {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to delete volume.", zap.Reflect("StorageType", volumes[i].VolumeID), zap.Reflect("Error", err))
			Expect(err).To(HaveOccurred())
		}
	}
	return
}