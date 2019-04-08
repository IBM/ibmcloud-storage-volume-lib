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
	"go.uber.org/zap"
)

// ListVolumes list all volumes
func (vpcs *VPCSession) ListVolumes(tags map[string]string) ([]*provider.Volume, error) {
	vpcs.Logger.Info("Entry ListVolumes", zap.Reflect("Tags", tags))
	defer vpcs.Logger.Info("Exit ListVolumes", zap.Reflect("Tags", tags))

	//! TODO: we may implement
	return nil, nil
}
