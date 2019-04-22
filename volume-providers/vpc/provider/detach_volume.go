/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	//	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
	//	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	//
	"go.uber.org/zap"
)

//const VpcVolumeAttachment = "vpcVolumeAttachment"

// Detach volume
func (vpcs *VPCSession) Detach(volumeDetachRequest provider.VolumeDetachRequest) (provider.VolumeResponse, error) {
	vpcs.Logger.Debug("Entry of Detach method...")
	defer vpcs.Logger.Debug("Exit from Detach method...")
	//var err error
	vpcs.Logger.Info("Validating basic inputs for Detach method...", zap.Reflect("volumeDetachRequest", volumeDetachRequest))
	volumeResponse := provider.VolumeResponse{}
	return volumeResponse, nil
}
