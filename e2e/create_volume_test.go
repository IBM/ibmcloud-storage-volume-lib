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
func (vpc *VpcClassicE2E) GetVolumeRequests() []VolumeRequest {
	requestList := []VolumeRequest{}
	volName := volumeName + "10iops"
	volSize := volumeSize
	Iops := iops
	volume := &provider.Volume{}
	volume.VolumeType = volumeType
	volume.VPCVolume.Generation = generation
	volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}
	volume.VPCVolume.Profile = &provider.Profile{Name: "10iops-tier"}
	volume.Name = &volName
	volume.Capacity = &volSize
	volume.Iops = &Iops
	volume.VPCVolume.ResourceGroup.ID = resourceGroupID
	volume.Az = vpcZone

	volume.VPCVolume.Tags = []string{"Testing VPC volume  with 10iops-tier profile"}

	volumeRequest := VolumeRequest{}
	volumeRequest.Volume = *volume
	volumeRequest.TestName = *volume.Name
	volumeRequest.AssertError = func(err error) {
		Expect(err).NotTo(HaveOccurred())

	}
	volumeRequest.AssertResult = func(volume *provider.Volume) {
		//Expect(volume.Name).To(Equal(volumeRequest.Volume.Name))
		Expect("3000").To(Equal(*volume.Iops))
		//Expect(volumeRequest.Volume.Az).To(Equal(volume.Az))
		Expect(volumeRequest.Volume.Capacity).To(Equal(volume.Capacity))
		//Expect(volumeRequest.Volume.VPCVolume.Generation).To(Equal(volume.VPCVolume.Generation))
		//Expect(volumeRequest.Volume.VPCVolume.Profile.Name).To(Equal(volume.VPCVolume.Profile.Name))

	}
	requestList = append(requestList, volumeRequest)
	// Check if volume create fails with invalid iops
	volumeRequest1 := volumeRequest.Clone()
	volumeRequest1.TestName = "10-iops-tier with explicit iops"
	iops := "100"
	volumeRequest1.Volume.Iops = &iops
	volumeRequest1.AssertResult = func(volume *provider.Volume) {
		Expect(volume).Should(BeNil())
	}
	volumeRequest1.AssertError = func(err error) {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("VolumeProfileIopsInvalid"))
	}
	requestList = append(requestList, volumeRequest1)

	return requestList

}

func (vpc *SLBlockE2E) GetVolumeRequests() []VolumeRequest {
	requestList := []VolumeRequest{}
	volSize := 20
	tier := "0.25"
	Iops := iops
	volume := &provider.Volume{}
	// Create volume with endurance
	volume.VolumeType = "block"
	volume.ProviderType = provider.VolumeProviderType("endurance")
	volume.Tier = &tier
	volume.Capacity = &volSize
	volume.Iops = &Iops
	volume.Az = "dal10"

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
