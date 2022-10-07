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
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

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

					if len(testCase.Input.Volume.SnapshotName) > 0 && len(testCase.Input.InstanceIP) > 0 {
						By("Test Create VPC Instance ")
						cmd := exec.Command("./../scripts/create_inst.sh")
						var outb, errb bytes.Buffer
						cmd.Stdout = &outb
						cmd.Stderr = &errb
						_ = cmd.Run()
						fmt.Println("out:", outb.String(), "err:", errb.String())
						res := strings.Split(string(outb.String()), "output :")[1]
						res1 := strings.Split(res, ",")
						fmt.Println(res1)
						Instance_id := strings.Split(res1[0], "=")[1]
						Instance_ip := strings.Split(res1[1], "=")[1]
						key_id := strings.Split(res1[2], "=")[1]
						ip_address_id := strings.Split(res1[3], "=")[1]
						testCase.Input.InstanceIP[0] = Instance_ip
						testCase.Input.InstanceID[0] = Instance_id
						By("Test Attach Volume")
						volumeAttachmentsRequests, volumeAttachmentsResponses, _ := attachVolumes(testCase, volumesCreated)
						By("Format and Mount Volume")
						args := []string{Instance_ip, volumeAttachmentsResponses[0].VPCVolumeAttachment.ID, "ext4"}
						fmt.Println(args)
						fmt_mount_cmd := exec.Command("./../scripts/run_fmt_mount.sh", args...)
						var cmd_out, cmd_err bytes.Buffer
						fmt_mount_cmd.Stdout = &cmd_out
						fmt_mount_cmd.Stderr = &cmd_err
						_ = fmt_mount_cmd.Run()
						time.Sleep(15 * time.Second)
						fmt.Println("out:", cmd_out.String(), "err:", cmd_err.String())
						By("Test Create Snapshot")
						snapshotCreated, err = createSnapshot(testCase, volumesCreated)
						time.Sleep(240 * time.Second)

						By("Test Create Volume From Snapshot")
						testCase.Input.Volume.Name = testCase.Input.Volume.Name + "restore"
						testCase.Input.Volume.SnapshotID = snapshotCreated[0].SnapshotID
						volumeRestoreCreated, err = createVolumes(testCase)
						if len(volumeRestoreCreated) > 0 {
							By("Test Attach restored Volume")
							volumeAttachmentsRequest, volumeAttachmentsResponse, err = attachVolumes(testCase, volumeRestoreCreated)
							arg := []string{testCase.Input.InstanceIP[0], volumeAttachmentsResponse[0].VPCVolumeAttachment.ID, "ext4"}
							By("Mount Volume and Validate Data")
							validate_vol_size := exec.Command("./../scripts/run_validate_volume_size.sh", arg...)
							var out, err bytes.Buffer
							validate_vol_size.Stdout = &out
							validate_vol_size.Stderr = &err
							_ = validate_vol_size.Run()
							time.Sleep(15 * time.Second)
							fmt.Println("out:", out.String(), "err:", err.String())

							By("Test Detach Restored Volume")
							_ = detachVolumes(volumeAttachmentsRequest)

							By("Test Delete Restored Volume")
							_ = deleteVolumes(volumeRestoreCreated)
						}

						if len(snapshotCreated) > 0 {
							By("Test Delete Snapshot")
							err = deleteSnapshot(snapshotCreated)
						}

						By("Test Detach Volume")
						_ = detachVolumes(volumeAttachmentsRequests)
						arg := []string{Instance_id, key_id, ip_address_id}
						delete_res_cmd := exec.Command("./../scripts/delete_vpc_resources.sh", arg...)
						var cmdout, cmderr bytes.Buffer
						delete_res_cmd.Stdout = &cmdout
						delete_res_cmd.Stderr = &cmderr
						_ = delete_res_cmd.Run()
						time.Sleep(15 * time.Second)
						fmt.Println("out:", cmdout.String(), "err:", cmderr.String())

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
