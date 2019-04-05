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
)

//AuthorizeVolume allows aceess to volume  based on given authorization
func (vpcs *VPCSession) AuthorizeVolume(volumeAuthorization provider.VolumeAuthorization) error {
	vpcs.Logger.Info("Entry AuthorizeVolume")
	defer vpcs.Logger.Info("Exit AuthorizeVolume")

	return nil
}
