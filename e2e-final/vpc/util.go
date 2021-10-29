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
	"go.uber.org/zap"
	"time"
)

func createVolume(testCase TestCaseData) (*provider.Volume, error) {
	volume := getVolume(testCase)
	startTime = time.Now()
	volumeObj, err := sess.CreateVolume(*volume)
	if err == nil {
		ctxLogger.Info("Successfully created volume...", zap.Reflect("volumeObj", volumeObj))
	} else {
		err = updateRequestID(err, requestID)
		ctxLogger.Info("Failed to create volume...", zap.Reflect("StorageType", volume.ProviderType), zap.Reflect("Error", err))
	}

	ctxLogger.Info("Test Create Volume", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
	fmt.Printf("\n\n")
	return volumeObj, err
}

func deleteVolume(testCase TestCaseData, volumeObj *provider.Volume) error {

	startTime = time.Now()
	err = sess.DeleteVolume(volumeObj)
	if err == nil {
		ctxLogger.Info("Successfully deleted volume...", zap.Reflect("volumeObj", volumeObj))
	} else {
		err = updateRequestID(err, requestID)
		ctxLogger.Info("Failed to delete volume...", zap.Reflect("StorageType", volumeObj.VolumeID), zap.Reflect("Error", err))
	}

	ctxLogger.Info("Test Delete Volume", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
	fmt.Printf("\n\n")

	return err
}

func createVolumeAccessPoint(testCase TestCaseData, volumeID string) (*provider.VolumeAccessPointResponse, error) {
	volumeAccessPointRequest := provider.VolumeAccessPointRequest{VolumeID: volumeID,
		VPCID: testCase.Input.VPCID[0],
	}
	startTime = time.Now()
	response, err := sess.CreateVolumeAccessPoint(volumeAccessPointRequest)
	if err != nil {
		err = updateRequestID(err, requestID)
		ctxLogger.Info("Failed to create volume access point...", zap.Reflect("Error", err))
		return nil, err
	}

	volumeAccessPointRequest.AccessPointID = response.AccessPointID

	ctxLogger.Info("Waiting for CreateVolumeAccessPoint...")

	response, err = sess.WaitForCreateVolumeAccessPoint(volumeAccessPointRequest)

	if err != nil {
		err = updateRequestID(err, requestID)
		ctxLogger.Info("Failed to create volume access point...", zap.Reflect("Error", err))
	}

	ctxLogger.Info("VolumeAccessPoint created successfully")
	ctxLogger.Info("Test Create Volume Access Point", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))

	fmt.Printf("\n\n")
	return response, err
}

func deleteVolumeAccessPoint(testCase TestCaseData, volumeAccessPointResponse *provider.VolumeAccessPointResponse) error {
	volumeAccessPointRequest := provider.VolumeAccessPointRequest{VolumeID: volumeAccessPointResponse.VolumeID,
		AccessPointID: volumeAccessPointResponse.AccessPointID,
	}
	startTime = time.Now()
	repsonse, err := sess.DeleteVolumeAccessPoint(volumeAccessPointRequest)
	if err != nil {
		err = updateRequestID(err, requestID)
		ctxLogger.Info("Failed to delete volume access point...", zap.Reflect("Error", err))
		return err
	}

	ctxLogger.Info("Initiated deletion of  volume access point ...", zap.Reflect("repsonse", repsonse))

	err = sess.WaitForDeleteVolumeAccessPoint(volumeAccessPointRequest)
	if err != nil {
		err = updateRequestID(err, requestID)
		ctxLogger.Info("Failed to delete volume access point...", zap.Reflect("Error", err))
	}

	ctxLogger.Info("VolumeAccessPoint deleted successfully")
	ctxLogger.Info("Test Delete Volume Access Point", zap.Reflect("Elapsed time:", fmt.Sprintf("%s", time.Since(startTime))))
	fmt.Printf("\n\n")

	return err
}

func getVolume(testCase TestCaseData) *provider.Volume {
	volume := &provider.Volume{}

	volume.VPCVolume.ResourceGroup = &provider.ResourceGroup{}
	volume.VPCVolume.Profile = &provider.Profile{Name: testCase.Input.Volume.Profile}
	volume.Name = &testCase.Input.Volume.Name
	volume.Capacity = &testCase.Input.Volume.Capacity
	volume.Iops = &testCase.Input.Volume.Iops
	volume.VPCVolume.ResourceGroup.ID = resourceGroupID
	volume.Az = vpcZone
	volume.VPCVolume.Tags = []string{testCase.Input.Volume.Tags}
	if testCase.Input.VolumeEncryptionKey != "" {
		volume.VPCVolume.VolumeEncryptionKey = &provider.VolumeEncryptionKey{}
		volume.VPCVolume.VolumeEncryptionKey.CRN = testCase.Input.VolumeEncryptionKey
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
