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

				By("Test Create Volume")
				fmt.Println(testCase)
				volumeCreated, err = createVolume(testCase)

				if volumeCreated != nil {

					//This case would enhance for creating file access points per VPC, as of now we will do it for one VPC
					if testCase.Input.VPCID != nil && testCase.Input.VPCID[0] != "" {

						//This case for VPC File library to test create/delete access point
						By("Test Create Volume Access Point")

						volumeAccessPointCreated, err = createVolumeAccessPoint(testCase, volumeCreated.VolumeID)

						if volumeAccessPointCreated != nil {

							By("Test Delete Volume Access Point")
							err = deleteVolumeAccessPoint(testCase, volumeAccessPointCreated)

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
