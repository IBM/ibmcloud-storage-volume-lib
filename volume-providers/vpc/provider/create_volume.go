/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"go.uber.org/zap"
)

// CreateVolume Get the volume by using ID
func (vpcs *VPCSession) CreateVolume(volumeRequest provider.Volume) (*provider.Volume, error) {
	vpcs.Logger.Info("Creating volume as per order request .... ", zap.Reflect("Volume", volumeRequest))
	return nil, nil
}
