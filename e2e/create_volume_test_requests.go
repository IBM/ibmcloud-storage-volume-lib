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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	. "github.com/onsi/gomega"
)

//GetVolumeRequests ...
func (vpc *VPCClassicE2E) GetTestVolumeRequests() []VolumeRequest {
	logger.Info("GetVolumeRequests")
	requestList := []VolumeRequest{}
	volName := volumeName + "5iops"
	volSize := volumeSize
	Iops := iops
	volume := &provider.Volume{}
	volume.VolumeType = volumeType
	volume.VPCVolume.Generation = generation
	volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}
	volume.VPCVolume.Profile = &provider.Profile{Name: "5iops-tier"}
	volume.Name = &volName
	volume.Capacity = &volSize
	volume.Iops = &Iops
	volume.VPCVolume.ResourceGroup.ID = resourceGroupID
	volume.Az = vpcZone

	volume.VPCVolume.Tags = []string{"Testing create VPC volume"}

	volumeRequest := VolumeRequest{}
	volumeRequest.Volume = *volume
	volumeRequest.TestName = "Testing create VPC volume with 5iops-tier profile"
	volumeRequest.AssertError = func(err error) {
		Expect(err).NotTo(HaveOccurred())
	}
	volumeRequest.AssertResult = func(volume *provider.Volume) {
		Expect(volume).ShouldNot(BeNil())
		Expect("3000").To(Equal(*volume.Iops))
		Expect(volumeRequest.Volume.Capacity).To(Equal(volume.Capacity))
	}
	requestList = append(requestList, volumeRequest)

	// Check if volume create fails with invalid iops
	volumeRequest = volumeRequest.Clone()
	volumeRequest.TestName = "5-iops-tier with explicit iops"
	iops := "100"
	volumeRequest.Volume.Iops = &iops
	volumeRequest.AssertResult = func(volume *provider.Volume) {
		Expect(volume).Should(BeNil())
	}
	volumeRequest.AssertError = func(err error) {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("VolumeProfileIopsInvalid"))
	}
	requestList = append(requestList, volumeRequest)

	volumeRequest = volumeRequest.Clone()
	volumeRequest.Volume.Az = vpcZone
	volumeRequest.Volume.Iops = &Iops
	volumeRequest.Volume.VPCVolume.ResourceGroup.ID = resourceGroupID
	volumeRequest.Volume.VPCVolume.Profile = &provider.Profile{Name: "10iops-tier"}
	volumeRequest.TestName = "Testing create VPC volume with 10iops-tier profile"
	volumeRequest.AssertError = func(err error) {
		Expect(err).NotTo(HaveOccurred())
	}
	volumeRequest.AssertResult = func(volume *provider.Volume) {
		Expect(volume).ShouldNot(BeNil())
		Expect("3000").To(Equal(*volume.Iops))
		Expect(volumeRequest.Volume.Capacity).To(Equal(volume.Capacity))
	}
	requestList = append(requestList, volumeRequest)

	volumeRequest = volumeRequest.Clone()
	volumeRequest.Volume.Az = vpcZone
	volumeRequest.Volume.Iops = &Iops
	volumeRequest.Volume.VPCVolume.ResourceGroup.ID = resourceGroupID
	volumeRequest.Volume.VPCVolume.Profile = &provider.Profile{Name: "custom"}
	volumeRequest.TestName = "Testing create VPC volume with custom profile"
	volumeRequest.AssertError = func(err error) {
		Expect(err).NotTo(HaveOccurred())
	}
	volumeRequest.AssertResult = func(volume *provider.Volume) {
		Expect(volume).ShouldNot(BeNil())
		Expect("100").To(Equal(*volume.Iops))
		Expect(volumeRequest.Volume.Capacity).To(Equal(volume.Capacity))
	}
	requestList = append(requestList, volumeRequest)

	return requestList

}

func (vpc *SLBlockE2E) GetTestVolumeRequests() []VolumeRequest {
	requestList := []VolumeRequest{}
	volSize := 20
	tier := "0.25"
	Iops := iops
	volume := &provider.Volume{}
	// Create volume with endurance
	volume.VolumeType = volumeType
	volume.ProviderType = provider.VolumeProviderType("endurance")
	volume.Tier = &tier
	volume.Capacity = &volSize
	volume.Iops = &Iops
	volume.Az = vpcZone

	volume.VolumeNotes = map[string]string{"note": "ibm-volume-lib-test"}
	volumeRequest := VolumeRequest{}
	testName := fmt.Sprintf("%s_%s_%d", volume.ProviderType, tier, volume.Capacity)
	volumeRequest.Volume = *volume
	volumeRequest.TestName = testName

	volumeRequest.AssertError = func(err error) {
		Expect(err).NotTo(HaveOccurred())
	}
	volumeRequest.AssertResult = func(volume *provider.Volume) {
		Expect("0.25").To(Equal(*volume.Iops))
		Expect(volumeRequest.Volume.Capacity).To(Equal(volume.Capacity))
	}
	requestList = append(requestList, volumeRequest)
	// Create volume with performance
	volumeRequest1 := volumeRequest.Clone()
	volumeRequest1.Volume.ProviderType = provider.VolumeProviderType("performance")
	requestList = append(requestList, volumeRequest)
	return requestList
}
