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
	//"fmt"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

/* wrtie utility method for slack*/

//VolumeRequest ...
type VolumeRequest struct {
	Volume       provider.Volume
	TestName     string                 // e.g "VPC Create Volume with Custom profile"
	AssertError  func(error)            // Special error assertion for this request
	AssertResult func(*provider.Volume) // Special result assertion for this request

}

func (volReq *VolumeRequest) Clone() VolumeRequest {
	srcVol := volReq.Volume
	newVol := provider.Volume{
		VolumeID:    srcVol.VolumeID,
		VolumeType:  srcVol.VolumeType,
		VolumeNotes: srcVol.VolumeNotes,
		VPCVolume:   volReq.Volume.VPCVolume,
		Name:        srcVol.Name,
		Capacity:    srcVol.Capacity,
	}
	newVolReq := VolumeRequest{
		Volume:       newVol,
		TestName:     volReq.TestName,
		AssertError:  volReq.AssertError,
		AssertResult: volReq.AssertResult,
	}
	return newVolReq
}

type VolumeAttachmentRequest struct {
	VolumeRequest provider.VolumeAttachmentRequest
	TestName      string                                  // e.g "VPC  Volume  attachment "
	AssertError   func(error)                             // Special error assertion for this request
	AssertResult  func(provider.VolumeAttachmentResponse) // Special result assertion for this request
}
type ProviderE2ETest interface {
	SetUp()
	GetName() string
	GetVolumeRequests() []VolumeRequest // Rename
	GetVolumeAttachmentRequests(*provider.Volume) []VolumeAttachmentRequest
	TestCreateVolume(VolumeRequest) *provider.Volume
	TestAttachVolume(VolumeAttachmentRequest)
	TestDetachVolume(VolumeAttachmentRequest)
	TestDeleteVolume(provider.Volume)
	TestAuthorizeVolume(provider.Volume)
	TestDeAuthorizeVolume(provider.Volume)
	TearDown()
}

var _ ProviderE2ETest = &BaseE2ETest{}

type BaseE2ETest struct {
	session provider.Session
	Name    string
}

func (b *BaseE2ETest) SetUp() {

}
func (b *BaseE2ETest) GetName() string {
	return b.Name
}
func (b *BaseE2ETest) TestCreateVolume(volumeRequest VolumeRequest) *provider.Volume {
	vol, err := b.session.CreateVolume(volumeRequest.Volume)
	volumeRequest.AssertError(err)
	volumeRequest.AssertResult(vol)
	return vol
}
func (b *BaseE2ETest) TestAttachVolume(VolumeAttachmentRequest) {

}
func (b *BaseE2ETest) TestDetachVolume(VolumeAttachmentRequest) {

}
func (b *BaseE2ETest) TestDeleteVolume(volume provider.Volume) {
	err := b.session.DeleteVolume(&volume)
	Expect(err).NotTo(HaveOccurred())

}
func (b *BaseE2ETest) TestAuthorizeVolume(volumes provider.Volume) {

}
func (b *BaseE2ETest) TestDeAuthorizeVolume(volumes provider.Volume) {

}

func (b *BaseE2ETest) GetVolumeRequests() []VolumeRequest {
	requestList := []VolumeRequest{}
	/*	volumeRequest := VolumeRequest{
			TestName: "My demo test",
		}
		volumeRequest.AssertError = func(err error) {
			Expect(err).NotTo(HaveOccurred())

		}
		volumeRequest.AssertResult = func(volume *provider.Volume) {
			Expect(volumeRequest.TestName).To(Equal("My demo test"))

		}
		requestList = append(requestList, volumeRequest)*/
	return requestList

}

//GetVolumeAttachmentRequests ...
func (b *BaseE2ETest) GetVolumeAttachmentRequests(volume *provider.Volume) []VolumeAttachmentRequest {
	requestList := []VolumeAttachmentRequest{}
	volumeRequest := VolumeAttachmentRequest{
		TestName: "My demo test",
	}
	volumeRequest.AssertError = func(err error) {
		Expect(err).NotTo(HaveOccurred())

	}
	volumeRequest.AssertResult = func(volume provider.VolumeAttachmentResponse) {

	}
	requestList = append(requestList, volumeRequest)
	return requestList

}

func (b *BaseE2ETest) TearDown() {

}

type SLFileE2E struct {
	BaseE2ETest
}
type SLBlockE2E struct {
	BaseE2ETest
}

//VpcClassicE2E ...
type VpcClassicE2E struct {
	BaseE2ETest
}
type IksVpcClassicE2E struct {
	BaseE2ETest
}

type VpcNextGenE2E struct {
	BaseE2ETest
}
type SharedVolume struct {
	Volume *provider.Volume
}

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	initSuite()
	var vol *provider.Volume
	var sharedVolume SharedVolume
	AssertCreateVolume := func(providere2e ProviderE2ETest, volRequest VolumeRequest) SharedVolume {
		sharedVolume = SharedVolume{}
		It(volRequest.TestName, func() {
			sharedVolume.Volume = providere2e.TestCreateVolume(volRequest)
		})
		return sharedVolume
	}
	AssertAttachVolume := func(providere2e ProviderE2ETest, volAttachReq VolumeAttachmentRequest) {
		It("DeleteVolume", func() {
			providere2e.TestAttachVolume(volAttachReq)
		})
	}
	AssertDetachVolume := func(providere2e ProviderE2ETest, volAttachReq VolumeAttachmentRequest) {
		It("DeleteVolume", func() {
			providere2e.TestDetachVolume(volAttachReq)
		})
	}

	AssertDeleteVolume := func(providere2e ProviderE2ETest) {
		It("DeleteVolume", func() {
			if sharedVolume.Volume != nil {
				providere2e.TestDeleteVolume(*sharedVolume.Volume)
			}
		})
	}
	for _, e2eProvider := range providers {
		Context("Context :"+e2eProvider.GetName(), func() {

			volumeRequests := e2eProvider.GetVolumeRequests()
			for _, volumeRequest := range volumeRequests {
				sharedVolume = AssertCreateVolume(e2eProvider, volumeRequest)

				volumeAttachmentRequests := e2eProvider.GetVolumeAttachmentRequests(vol)
				for _, volAttachReq := range volumeAttachmentRequests {
					AssertAttachVolume(e2eProvider, volAttachReq)
					AssertDetachVolume(e2eProvider, volAttachReq)
				}
				AssertDeleteVolume(e2eProvider)

			}

		})
	}

	/*for _, providere2e := range providers {
		Describe("Initialising the provider e2e", func() {

			Context("When initialization is successfull", func() {

				//volumeRequests := providere2e.GetVolumeRequests()
				It(providere2e.GetName()+" Create Volume", func() {
					//volumes = providere2e.TestCreateVolume(volumeRequests)
				})
				volumeAttachmentRequests := providere2e.GetVolumeAttachmentRequests()
				It(providere2e.GetName()+" Attach Volume", func() {
					providere2e.TestAttachVolume(volumeAttachmentRequests)
				})
				It(providere2e.GetName()+" Detach Volume", func() {
					providere2e.TestDetachVolume(volumeAttachmentRequests)
				})
				It(providere2e.GetName()+" Delete Volume", func() {
					providere2e.TestDeleteVolume(volumes)
				})

			})

		})
	}*/
	/*var entries []TableEntry
	for _, providere2e := range providers {
		entries = append(entries, Entry(providere2e.GetName(), providere2e))
	}
	DescribeTable("Providers", func(providere2e ProviderE2ETest) {

		By(" Create Volume", func() {
			for _, volumeRequest := range volumeRequests {
				vol := AssertCreateVolumes(providere2e, volumeRequest)
				if vol != nil {
					volumes = append(volumes, *vol)
				}
			}
			//	volumes = providere2e.TestCreateVolume(volumeRequests)
		})
		volumeAttachmentRequests := providere2e.GetVolumeAttachmentRequests()
		By(" Attach Volume", func() {
			providere2e.TestAttachVolume(volumeAttachmentRequests)
		})
		By(" Detach Volume", func() {
			providere2e.TestDetachVolume(volumeAttachmentRequests)
		})
		By(" Delete Volume", func() {
			providere2e.TestDeleteVolume(volumes)
		})

	}, entries...,
	)*/
})
