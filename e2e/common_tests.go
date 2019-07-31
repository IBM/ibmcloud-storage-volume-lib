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
	"go.uber.org/zap"
	"time"
)

/* wrtie utility method for slack */

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
	GetTestVolumeRequests() []VolumeRequest // Rename
	GetTestVolumeAttachmentRequests(*provider.Volume) []VolumeAttachmentRequest
	TestCreateVolume(VolumeRequest) *provider.Volume
	TestAttachVolume(VolumeAttachmentRequest)
	TestDetachVolume(VolumeAttachmentRequest)
	TestDeleteVolume(provider.Volume)
	DeleteVolume(string)
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

func (b *BaseE2ETest) DeleteVolume(volumeName string) {
	volume, _ := b.session.GetVolumeByName(volumeName)
	if volume != nil {
		b.session.DeleteVolume(volume)
	}
}

func (b *BaseE2ETest) TestAuthorizeVolume(volumes provider.Volume) {

}

func (b *BaseE2ETest) TestDeAuthorizeVolume(volumes provider.Volume) {

}

func (b *BaseE2ETest) GetTestVolumeRequests() []VolumeRequest {
	// Blank implementation
	requestList := []VolumeRequest{}
	return requestList
}

//GetVolumeAttachmentRequests ...
func (b *BaseE2ETest) GetTestVolumeAttachmentRequests(volume *provider.Volume) []VolumeAttachmentRequest {
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

// SLFileE2E ...
type SLFileE2E struct {
	BaseE2ETest
}

// SLBlockE2E ...
type SLBlockE2E struct {
	BaseE2ETest
}

// VpcClassicE2E ...
type VPCClassicE2E struct {
	BaseE2ETest
}

// IksVPCClassicE2E ...
type IksVPCClassicE2E struct {
	BaseE2ETest
}

//VPCNextGenE2E ...
type VPCNextGenE2E struct {
	BaseE2ETest
}

// SharedVolume Used to pass volume to differnt spec
type SharedVolume struct {
	Volume *provider.Volume
}

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	var sharedVolume SharedVolume
	AssertCreateVolume := func(providere2e ProviderE2ETest, volRequest VolumeRequest) SharedVolume {
		sharedVolume = SharedVolume{}
		sharedVolume.Volume = providere2e.TestCreateVolume(volRequest)
		return sharedVolume
	}

	AssertAttachVolume := func(providere2e ProviderE2ETest, volAttachReq VolumeAttachmentRequest) {
		providere2e.TestAttachVolume(volAttachReq)
	}

	AssertDetachVolume := func(providere2e ProviderE2ETest, volAttachReq VolumeAttachmentRequest) {
		providere2e.TestDetachVolume(volAttachReq)
	}

	AssertDeleteVolume := func(providere2e ProviderE2ETest) {
		if sharedVolume.Volume != nil {
			providere2e.TestDeleteVolume(*sharedVolume.Volume)
		}
	}

	It("Provider E2E", func() {
		for _, e2eProvider := range providers {
			ctxLogger.Info("Provider", zap.Reflect("provider", e2eProvider.GetName()))
			// Get volume create request for each provider
			volumeRequests := e2eProvider.GetTestVolumeRequests()
			for _, volumeRequest := range volumeRequests {
				time.Sleep(10 * time.Second)
				// Delete the volume it was not deleted Successfully in the last run
				e2eProvider.DeleteVolume(*volumeRequest.Volume.Name)

				// For each volume create request perform following steps
				step := volumeRequest.TestName + ", Test Create Volume"
				By(step)
				sharedVolume = AssertCreateVolume(e2eProvider, volumeRequest)

				// Get volume attachment requests for this volume
				volumeAttachmentRequests := e2eProvider.GetTestVolumeAttachmentRequests(sharedVolume.Volume)
				for _, volAttachReq := range volumeAttachmentRequests {
					step = volumeRequest.TestName + ", Test Attach Volume"
					By(step)
					AssertAttachVolume(e2eProvider, volAttachReq)

					step = volumeRequest.TestName + ", Test Detach Volume"
					By(step)
					AssertDetachVolume(e2eProvider, volAttachReq)
				}

				step = volumeRequest.TestName + ", Test Delete Volume"
				By(step)
				AssertDeleteVolume(e2eProvider)
			}
		}
	})
})
