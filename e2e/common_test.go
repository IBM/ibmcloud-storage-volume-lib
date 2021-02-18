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
	. "github.com/onsi/ginkgo"
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
	// Blank implementation
	requestList := []VolumeRequest{}
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

//SLFileE2E ...
type SLFileE2E struct {
	BaseE2ETest
}

//SLBlockE2E ...
type SLBlockE2E struct {
	BaseE2ETest
}

//VpcClassicE2E ...
type VpcClassicE2E struct {
	BaseE2ETest
}

//IksVpcClassicE2E ...
type IksVpcClassicE2E struct {
	BaseE2ETest
}

//VpcNextGenE2E ...
type VpcNextGenE2E struct {
	BaseE2ETest
}

//SharedVolume Used to pass volume to differnt spec
type SharedVolume struct {
	Volume *provider.Volume
}

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	initSuite() // TODO e2e_suite is not geetting called before calling actual test
	var sharedVolume SharedVolume
	AssertCreateVolume := func(providere2e ProviderE2ETest, volRequest VolumeRequest) SharedVolume {
		sharedVolume = SharedVolume{}
		It(volRequest.TestName, func() {
			sharedVolume.Volume = providere2e.TestCreateVolume(volRequest)
		})
		return sharedVolume
	}
	AssertAttachVolume := func(providere2e ProviderE2ETest, volAttachReq VolumeAttachmentRequest) {
		It("AttachVolume", func() {
			providere2e.TestAttachVolume(volAttachReq)
		})
	}
	AssertDetachVolume := func(providere2e ProviderE2ETest, volAttachReq VolumeAttachmentRequest) {
		It("DetachVolume", func() {
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
	// Iterate over all the initialised providers
	for _, e2eProvider := range providers {
		Context("Context :"+e2eProvider.GetName(), func() {
			// Get volume create request for each provider
			volumeRequests := e2eProvider.GetVolumeRequests()
			for _, volumeRequest := range volumeRequests {
				// For each volume create request perform following steps
				By("Test Create Volume")
				sharedVolume = AssertCreateVolume(e2eProvider, volumeRequest)
				// Get volume attachment requests for this volume
				volumeAttachmentRequests := e2eProvider.GetVolumeAttachmentRequests(sharedVolume.Volume)
				for _, volAttachReq := range volumeAttachmentRequests {
					By("Test Attach  Volume")
					AssertAttachVolume(e2eProvider, volAttachReq)
					By("Test Detach  Volume")
					AssertDetachVolume(e2eProvider, volAttachReq)
				}
				By("Test Delete  Volume")
				AssertDeleteVolume(e2eProvider)
			}
		})
	}
})
