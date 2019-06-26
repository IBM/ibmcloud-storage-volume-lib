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
	GetVolumeAttachmentRequests() []VolumeAttachmentRequest
	TestCreateVolume([]VolumeRequest) []provider.Volume
	TestAttachVolume([]VolumeAttachmentRequest)
	TestDetachVolume([]VolumeAttachmentRequest)
	TestDeleteVolume([]provider.Volume)
	TestAuthorizeVolume([]provider.Volume)
	TestDeAuthorizeVolume([]provider.Volume)
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
func (b *BaseE2ETest) TestCreateVolume(volumeRequests []VolumeRequest) []provider.Volume {

	var volumeResultList []provider.Volume
	for _, volRequest := range volumeRequests {
		vol, err := b.session.CreateVolume(volRequest.Volume)
		volRequest.AssertError(err)
		volRequest.AssertResult(vol)
		if vol != nil {
			volumeResultList = append(volumeResultList, *vol)
		}
	}
	return volumeResultList
}
func (b *BaseE2ETest) TestAttachVolume([]VolumeAttachmentRequest) {

}
func (b *BaseE2ETest) TestDetachVolume([]VolumeAttachmentRequest) {

}
func (b *BaseE2ETest) TestDeleteVolume([]provider.Volume) {

}
func (b *BaseE2ETest) TestAuthorizeVolume([]provider.Volume) {

}
func (b *BaseE2ETest) TestDeAuthorizeVolume([]provider.Volume) {

}

func (b *BaseE2ETest) GetVolumeRequests() []VolumeRequest {
	requestList := []VolumeRequest{}
	volumeRequest := VolumeRequest{
		TestName: "My demo test",
	}
	volumeRequest.AssertError = func(err error) {
		Expect(err).NotTo(HaveOccurred())

	}
	volumeRequest.AssertResult = func(volume *provider.Volume) {
		Expect(volumeRequest.TestName).To(Equal("My demo test"))

	}
	requestList = append(requestList, volumeRequest)
	return requestList

}

//GetVolumeAttachmentRequests ...
func (b *BaseE2ETest) GetVolumeAttachmentRequests() []VolumeAttachmentRequest {
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

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	var (
		volumes []provider.Volume
	)
	BeforeEach(func() {
		//initSuite()
		fmt.Printf("Before each All pro %d", len(providers))
	})
	Describe("Initialising the provider e2e", func() {
		initSuite()
		Context("When initialization is successfull", func() {
			for _, providere2e := range providers {
				volumeRequests := providere2e.GetVolumeRequests()
				It("Create Volume", func() {
					volumes = providere2e.TestCreateVolume(volumeRequests)
				})
				volumeAttachmentRequests := providere2e.GetVolumeAttachmentRequests()
				It("Attach Volume", func() {
					providere2e.TestAttachVolume(volumeAttachmentRequests)
				})
				It("Detach Volume", func() {
					providere2e.TestDetachVolume(volumeAttachmentRequests)
				})
				It("Delete Volume", func() {
					providere2e.TestDeleteVolume(volumes)
				})
			}

		})

	})
})
