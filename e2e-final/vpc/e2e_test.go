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
		volumesCreated            []*provider.Volume
		snapshotCreated           []*provider.Snapshot
		volumeRestoreCreated      []*provider.Volume
		volumeAccessPointsRequest []provider.VolumeAccessPointRequest
		volumeAttachmentsResponse []*provider.VolumeAttachmentResponse
		volumeAttachmentsRequest  []*provider.VolumeAttachmentRequest
		err                       error
	)
	initializeTestCaseData()

	BeforeEach(func() {
		RefreshSession()
	})

	AfterEach(func() {
		deleteVolumes(volumesCreated)
		CloseSession()
	})

	Context("VPC e2e", func() {
		for _, testCase := range testCaseList {
			testCase := testCase //necessary to ensure the correct value is passed to the closure
			It(testCase.TestCase, func() {

				//Set default value of NumOfVolsRequired to 1 if not set
				//Could not find a way to override this value in yaml struct
				if testCase.Input.NumOfVolsRequired == 0 {
					testCase.Input.NumOfVolsRequired = 1
				}

				//Skip the test
				if testCase.Skip {
					Skip("Test was skipped, skip flag is true")
				}

				//Skip the IKS based Block storage attach/detach cases if iksEnabled is false
				if !conf.IKS.Enabled && len(testCase.Input.ClusterID) > 0 {
					Skip("Test was skipped, IKS is disable skipping IKS test cases")
				}

				//Skip the non-IKS based Block storage attach/detach cases if iksEnabled is true
				if conf.IKS.Enabled && len(testCase.Input.ClusterID) == 0 && len(testCase.Input.InstanceID) > 0 {
					Skip("Test was skipped, IKS is enabled skipping non-IKS test cases")
				}

				By("Test Create Volume")
				fmt.Println(testCase)
				volumesCreated, err = createVolumes(testCase)

				if len(volumesCreated) > 0 {

					//This case is for creating file access points per VPC, as of now we will do it for one VPC
					if len(testCase.Input.VPCID) > 0 && testCase.Input.VPCID[0] != "" {

						/*File Storage e2e specific handling
							  This case for VPC File library to test create/delete access point
						      TBD if we have input for more than one volumes then it would just use the same VPC-ID for
							  creating the access point accross the volumes*/

						By("Test Create Volume Access Point")
						volumeAccessPointsRequest, _, err = createVolumeAccessPoints(testCase, volumesCreated)

						if len(volumeAccessPointsRequest) > 0 {
							By("Test Delete Volume Access Point")
							err = deleteVolumeAccessPoints(volumeAccessPointsRequest)

						}
					}

					//This case is for creating volume attachment, as of now we will do it for one VPC
					if len(testCase.Input.InstanceID) > 0 && testCase.Input.InstanceID[0] != "" {

						By("Test Attach Volume")
						volumeAttachmentsRequest, volumeAttachmentsResponse, err = attachVolumes(testCase, volumesCreated)

						if len(volumeAttachmentsResponse) > 0 {
							By("Test Detach Volume")
							err = detachVolumes(volumeAttachmentsRequest)
						}

					}

					if len(testCase.Input.Volume.SnapshotName) > 0 {
						By("Test snapshot Create")
						snapshotCreated, err = createSnapshot(testCase, volumesCreated)

						By("Test Restore Snapshot")
						testCase.Input.Volume.Name = testCase.Input.Volume.Name + "restore"
						testCase.Input.Volume.SnapshotID = snapshotCreated[0].SnapshotID
						volumeRestoreCreated, err = createVolumes(testCase)

						By("Test Delete Restored Volume")
						err = deleteVolumes(volumeRestoreCreated)

						if len(snapshotCreated) > 0 {
							By("Test Delete Snapshot")
							err = deleteSnapshot(snapshotCreated)
						}

					}

					if testCase.Success {
						Expect(err).NotTo(HaveOccurred())
					} else {
						Expect(err).To(HaveOccurred())
					}

					By("Test Delete Volume")
					err = deleteVolumes(volumesCreated)
				}
			})
		}
	})
})
