/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpcvolume

import (
	"errors"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
)

// UpdateVolume POSTs to /volumes. Riaas/VPC does have volume update support yet
func (vs *VolumeService) UpdateVolume(volumeTemplate *models.Volume, ctxLogger *zap.Logger) error {

	return errors.New("Unsupported Operation")
}
