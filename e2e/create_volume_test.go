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

	volume.VPCVolume.Tags = []string{"Testing VPC volume from library with 10iops-tier profile"}

	volumeRequest := VolumeRequest{}
	volumeRequest.Volume = *volume
	volumeRequest.TestName = *volume.Name
	volumeRequest.AssertError = func(err error) {
		Expect(err).NotTo(HaveOccurred())

	}
	volumeRequest.AssertResult = func(volume *provider.Volume) {
		Expect(volume.Name).To(Equal(volumeRequest.Volume.Name))
		Expect(volume.Iops).To(Equal(volumeRequest.Volume.Iops))
		Expect(volume.Az).To(Equal(volumeRequest.Volume.Az))
		Expect(volume.Capacity).To(Equal(volumeRequest.Volume.Capacity))
		Expect(volume.VPCVolume.Generation).To(Equal(volumeRequest.Volume.VPCVolume.Generation))
		Expect(volume.VPCVolume.Profile).To(Equal(volumeRequest.Volume.VPCVolume.Profile))

	}
	requestList = append(requestList, volumeRequest)
	volumeRequest1 := VolumeRequest{}
	copy(volumeRequest1, volumeRequest)
	return requestList

}
