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
	"strconv"
	"time"

	"github.com/IBM/ibmcloud-volume-interface/lib/provider"
	"go.uber.org/zap"
)

func createVolumes(testCase TestCaseData) ([]*provider.Volume, error) {
	var volumes = make([]*provider.Volume, 0)
	startTime = time.Now()
	for i := 0; i < testCase.Input.NumOfVolsRequired; i++ {
		volume := getVolume(testCase, i)

		volumeObj, err := sess.CreateVolume(*volume)
		if err == nil {
			volumes = append(volumes, volumeObj)
			ctxLogger.Info("Successfully created volume...", zap.Reflect("volumeObj", volumeObj))
		} else {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to create volume...", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", err))
			return volumes, err
		}
	}
	ctxLogger.Info("Test Create Volume", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
	fmt.Printf("\n\n")

	return volumes, err
}

func deleteVolumes(volumes []*provider.Volume) error {
	startTime = time.Now()
	for i := 0; i < len(volumes); i++ {
		err = sess.DeleteVolume(volumes[i])
		if err == nil {
			ctxLogger.Info("Successfully deleted volume...", zap.Reflect("volumeObj", volumes[i]))
		} else {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to delete volume...", zap.Reflect("StorageType", volumes[i].VolumeID), zap.Reflect("Error", err))
		}

	}
	ctxLogger.Info("Test Delete Volume", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
	fmt.Printf("\n\n")

	return err
}

func createVolumeAccessPoints(testCase TestCaseData, volumes []*provider.Volume) ([]provider.VolumeAccessPointRequest, []*provider.VolumeAccessPointResponse, error) {

	startTime = time.Now()
	var volumeAccessPointsResponse = make([]*provider.VolumeAccessPointResponse, 0)
	var volumeAccessPointsRequest = make([]provider.VolumeAccessPointRequest, 0)

	for i := 0; i < len(volumes); i++ {

		volumeAccessPointRequest := provider.VolumeAccessPointRequest{VolumeID: volumes[i].VolumeID,
			VPCID: testCase.Input.VPCID[0],
		}

		response, err := sess.CreateVolumeAccessPoint(volumeAccessPointRequest)
		if err != nil {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to create volume access point...", zap.Reflect("Error", err))
			return volumeAccessPointsRequest, volumeAccessPointsResponse, err
		}

		volumeAccessPointRequest.AccessPointID = response.AccessPointID
		volumeAccessPointsRequest = append(volumeAccessPointsRequest, volumeAccessPointRequest)
		ctxLogger.Info("volumeAccessPointRequest", zap.Reflect("volumeAccessPointRequest", volumeAccessPointRequest))
	}

	ctxLogger.Info("Waiting for CreateVolumeAccessPoint...", zap.Reflect("volumeAccessPointRequest", volumeAccessPointsRequest))

	for i := 0; i < len(volumeAccessPointsRequest); i++ {
		response, err := sess.WaitForCreateVolumeAccessPoint(volumeAccessPointsRequest[i])

		if err != nil {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to create volume access point...", zap.Reflect("Error", err))
			return volumeAccessPointsRequest, volumeAccessPointsResponse, err
		}

		volumeAccessPointsResponse = append(volumeAccessPointsResponse, response)
	}

	ctxLogger.Info("VolumeAccessPoint created successfully")
	ctxLogger.Info("Test Create Volume Access Point", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))

	fmt.Printf("\n\n")
	ctxLogger.Info("volumeAccessPointRequestList", zap.Reflect("volumeAccessPointsRequest", volumeAccessPointsRequest))
	return volumeAccessPointsRequest, volumeAccessPointsResponse, err
}

func deleteVolumeAccessPoints(volumeAccessPointsRequest []provider.VolumeAccessPointRequest) error {
	startTime = time.Now()
	for i := 0; i < len(volumeAccessPointsRequest); i++ {

		repsonse, err := sess.DeleteVolumeAccessPoint(volumeAccessPointsRequest[i])
		if err != nil {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to delete volume access point...", zap.Reflect("Error", err))
			return err
		}

		ctxLogger.Info("Initiated deletion of  volume access point ...", zap.Reflect("repsonse", repsonse))
	}

	for i := 0; i < len(volumeAccessPointsRequest); i++ {
		err = sess.WaitForDeleteVolumeAccessPoint(volumeAccessPointsRequest[i])
		if err != nil {
			err = updateRequestID(err, requestID)
			ctxLogger.Info("Failed to delete volume access point...", zap.Reflect("Error", err))
		}

	}
	ctxLogger.Info("VolumeAccessPoint deleted successfully")
	ctxLogger.Info("Test Delete Volume Access Point", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
	fmt.Printf("\n\n")

	return err
}

func getVolume(testCase TestCaseData, index int) *provider.Volume {
	volumeName := testCase.Input.Volume.Name + "-" + strconv.Itoa(index+1)
	volume := &provider.Volume{}
	volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}
	volume.VPCVolume.Profile = &provider.Profile{Name: testCase.Input.Volume.Profile}
	volume.Name = &volumeName
	volume.Capacity = &testCase.Input.Volume.Capacity
	volume.Iops = &testCase.Input.Volume.Iops
	volume.VPCVolume.ResourceGroup.ID = resourceGroupID
	volume.Az = testCase.Input.VPCZone
	volume.VPCVolume.Tags = []string{testCase.Input.Volume.Tags}
	if volumeEncryptionKeyCRN != "" && testCase.Input.EncryptionEnabled {
		volume.VPCVolume.VolumeEncryptionKey = &provider.VolumeEncryptionKey{}
		volume.VPCVolume.VolumeEncryptionKey.CRN = volumeEncryptionKeyCRN
	}

	//Case for testing gid/uid for File VPC storage
	if testCase.Input.Volume.InitialOwner {
		volume.VPCFileVolume.InitialOwner = &provider.InitialOwner{
			GroupID: 1000,
			UserID:  1000,
		}
	}

	return volume
}

func attachVolumes(testCase TestCaseData, volumes []*provider.Volume) ([]*provider.VolumeAttachmentRequest, []*provider.VolumeAttachmentResponse, error) {
	startTime = time.Now()
	var attachmentResponses = make([]*provider.VolumeAttachmentResponse, 0)
	volumeAttachRequests := getVolumeAttachmentsRequest(testCase, volumes)

	for i := 0; i < len(volumeAttachRequests); i++ {

		attachResponse, err := sess.AttachVolume(*volumeAttachRequests[i])

		if attachResponse != nil && err == nil {
			attachmentResponses = append(attachmentResponses, attachResponse)
		} else {
			ctxLogger.Error("Error in attaching the volume.", zap.Reflect("err", err))
			return volumeAttachRequests, attachmentResponses, err
		}
	}

	for i := 0; i < len(volumeAttachRequests); i++ {
		attachResponse, err := sess.WaitForAttachVolume(*volumeAttachRequests[i])
		if err != nil {
			ctxLogger.Error("Error in attaching the volume.", zap.Reflect("err", err))
			return volumeAttachRequests, attachmentResponses, err
		}
		ctxLogger.Info("Successfully attached the volume.", zap.Reflect("attachResponse", attachResponse))
	}

	ctxLogger.Info("Test Attach Volume", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
	fmt.Printf("\n\n")
	return volumeAttachRequests, attachmentResponses, err
}

func detachVolumes(volumeAttachRequests []*provider.VolumeAttachmentRequest) error {
	startTime = time.Now()

	for i := 0; i < len(volumeAttachRequests); i++ {
		_, err := sess.DetachVolume(*volumeAttachRequests[i])
		if err != nil {
			ctxLogger.Error("Error in detaching the volume.", zap.Reflect("err", err))
			return err
		}
	}

	for i := 0; i < len(volumeAttachRequests); i++ {
		err := sess.WaitForDetachVolume(*volumeAttachRequests[i])
		if err != nil {
			ctxLogger.Error("Error in detaching the volume.", zap.Reflect("err", err))
			return err
		}

		ctxLogger.Info("Successfully detached the volume.")
	}

	ctxLogger.Info("Test Detach Volume", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
	fmt.Printf("\n\n")
	return err
}

func getVolumeAttachmentsRequest(testCase TestCaseData, volumes []*provider.Volume) []*provider.VolumeAttachmentRequest {
	var volumeAttachRequests = make([]*provider.VolumeAttachmentRequest, 0)

	for i := 0; i < len(volumes); i++ {
		volumeAttachRequest := &provider.VolumeAttachmentRequest{}
		volumeAttachRequest.VolumeID = volumes[i].VolumeID
		volumeAttachRequest.InstanceID = testCase.Input.InstanceID[0]
		volumeAttachRequest.IKSVolumeAttachment = &provider.IKSVolumeAttachment{}

		//This would be populate only for Block IKS based e2e
		if len(testCase.Input.ClusterID) > 0 && testCase.Input.ClusterID[0] != "" {
			volumeAttachRequest.IKSVolumeAttachment.ClusterID = &testCase.Input.ClusterID[0]
		}
		volumeAttachRequests = append(volumeAttachRequests, volumeAttachRequest)
	}
	return volumeAttachRequests
}
