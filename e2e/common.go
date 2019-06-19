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
)


// wrtie utility method for slack


type VolumeRequest struct {
	Volume       provider.Volume
	TestName     string                // e.g "VPC Create Volume with Custom profile"
	AssertError  func(error)           // Special error assertion for this request
	AssertResult func(provider.Volume) // Special result assertion for this request
  SetUp()
  TearDown()

}
type VolumeAttachmentRequest struct {
	VolumeRequest provider.VolumeAttachmentRequest
	TestName      string                                  // e.g "VPC  Volume  attachment "
	AssertError   func(error)                             // Special error assertion for this request
	AssertResult  func(provider.VolumeAttachmentResponse) // Special result assertion for this request
}
type ProviderE2ETest interface {
	SetUp()
	GetVolumeRequests() []VolumeRequest // Rename
	GetVolumeAttachmentRequests() []VolumeAttachmentRequest
	TestCreateVolume([]VolumeRequest) []provider.Volume
	TestAttachVolume([]provider.VolumeAttachmentRequest)
	TestDetachVolume([]provider.VolumeAttachmentRequest)
	TestDeleteVolume([]provider.Volume)
	TestAuthorizeVolume([]provider.Volume)
	TestDeAuthorizeVolume([]provider.Volume)
	TearDown()
}

var _ ProviderE2ETest = BaseE2ETest{}

type BaseE2ETest struct {
  session provider.Session
	// Unimplemented all ther method
}

func (b *BaseE2ETest) SetUp() {

}
func (b *BaseE2ETest) TestCreateVolume(volumeRequests []provider.VolumeRequest) []provider.Volume {
	var volumeRequesList []provider.Volume
  for volRequest  := range volumeRequests {
      vol, err := session.CreateVolume()
      It()
      volumeRequests.AssertError(err)
      volumeRequests.AssertResult(vol)
  }
	return volumeRequesList
}
func (b *BaseE2ETest) TestAttachVolume([]provider.VolumeAttachmentRequest) {

}
func (b *BaseE2ETest) TestDetachVolume([]provider.VolumeAttachmentRequest) {

}
func (b *BaseE2ETest) TestDeleteVolume([]provider.Volume) {

}
func (b *BaseE2ETest) TestAuthorizeVolume([]provider.Volume) {

}
func (b *BaseE2ETest) TestDeAuthorizeVolume([]provider.Volume) {

}

func (b *BaseE2ETest) GetVolumeRequests([]provider.Volume) {
  volumeRequest :=  &VolumeRequest{

  }
  volumeRequest.AssertError= func(error){

  }
  volumeRequest.AssertResult= func(volume){

  }
}

func (b *BaseE2ETest) TearDown() {

}

type SLFileE2E struct {
  BaseE2ETest
}
type SLBlockE2E struct {
  BaseE2ETest
}

type VpcClassicE2E struct {
  BaseE2ETest
}
type IksVpcClassicE2E struct {
  BaseE2ETest
}

type VpcNextGenE2E struct {
  BaseE2ETest
}

type E2ETestRunner struct {
  config := readConfig() // differnt toml file for test
	globalTestSetup()                 // Create VM, create cluster// Use before suite after sute
	providers := initProvider(config) // Initialse all the enabled provider in config file
var _ = Describe("ibmcloud-storage-volume-lib", func() {
  for provider := range providers {
    Context(provider.Name, func() {
		volumeRequests := provider.GetVolumeRequests()
    Context("Create Volume", func() {
		    provider.TestCreateVolume(volumeRequests)
    }
		volumeAttachmentRequests := provider.GetVolumeAttachmentRequests()
    Context("Attach Volume", func() {
		 provider.TestAttachVolume(volumeAttachmentRequests)
   }
     Context("Detach Volume", func() {
		provider.DetachAttachVolume(volumeAttachmentRequests)
    }
      Context("Delete Volume", func() {
		provider.TestDeleteVolume(volumeRequests)
    }
   }
	}
}
}
