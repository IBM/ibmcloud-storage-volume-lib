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
	"github.com/IBM/ibmcloud-volume-interface/lib/provider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ibmcloud-storage-volume-lib", func() {
	var (
		volumeCreated            *provider.Volume
		volumeAccessPointCreated *provider.VolumeAccessPointResponse
		volumeAttachmentResponse *provider.VolumeAttachmentResponse
		err                      error
	)
	initializeTestCaseData()

	BeforeEach(func() {
		RefreshSession()
	})

	AfterEach(func() {
		sess.DeleteVolume(volumeCreated)
		CloseSession()
	})

	Context("VPC e2e", func() {
		for _, testCase := range testCaseList {
			testCase := testCase //necessary to ensure the correct value is passed to the closure
			It(testCase.TestCase, func() {

				//Skip the test
				if testCase.Skip {
					Skip("Test was skipped, skip flag is true")
				}

				//Skip the IKS based Block storage attach/detach cases if iksEnabled is false
				if !conf.IKS.Enabled && len(testCase.Input.ClusterID) > 0 {
					Skip("Test was skipped, IKS is disable skipping IKS test cases")
				}

				//Skip the non-IKS based Block storage attach/detach cases if iksEnabled is true
				if conf.IKS.Enabled && len(testCase.Input.ClusterID) == 0 && len(testCase.Input.InstanceID) > 0  {
					Skip("Test was skipped, IKS is enabled skipping non-IKS test cases")
				}

				By("Test Create Volume")
				fmt.Println(testCase)
				volumeCreated, err = createVolume(testCase)

				if volumeCreated != nil {

					//This case is for creating file access points per VPC, as of now we will do it for one VPC
					if len(testCase.Input.VPCID) > 0 && testCase.Input.VPCID[0] != "" {
						//File Storage e2e specific handling
						//This case for VPC File library to test create/delete access point
						By("Test Create Volume Access Point")
						volumeAccessPointCreated, err = createVolumeAccessPoint(testCase, volumeCreated.VolumeID)

						if volumeAccessPointCreated != nil {
							By("Test Delete Volume Access Point")
							err = deleteVolumeAccessPoint(testCase, volumeAccessPointCreated)

						}

					}

					//This case is for creating volume attachment, as of now we will do it for one VPC
					if len(testCase.Input.InstanceID) > 0 && testCase.Input.InstanceID[0] != "" {

						By("Test Attach Volume")
						volumeAttachmentResponse, err = attachVolume(testCase, volumeCreated.VolumeID)

						if volumeAttachmentResponse != nil {
							By("Test Detach Volume")
							err = detachVolume(testCase, volumeCreated.VolumeID)
						}

					}

					if testCase.Success {
						Expect(err).NotTo(HaveOccurred())
					} else {
						Expect(err).To(HaveOccurred())
					}

					By("Test Delete Volume")
					err = deleteVolume(testCase, volumeCreated)
				}
			})
		}
	})
})
