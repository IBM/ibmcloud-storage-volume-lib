/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package auth

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/iam"
)

// ContextCredentialsFactory ...
type ContextCredentialsFactory struct {
	softlayerConfig      *config.SoftlayerConfig
	vpcConfig            *config.VPCProviderConfig
	tokenExchangeService iam.TokenExchangeService
}

var _ local.ContextCredentialsFactory = &ContextCredentialsFactory{}

// NewContextCredentialsFactory ...
func NewContextCredentialsFactory(bluemixConfig *config.BluemixConfig, softlayerConfig *config.SoftlayerConfig, vpcConfig *config.VPCProviderConfig) (*ContextCredentialsFactory, error) {
	tokenExchangeService, err := iam.NewTokenExchangeService(bluemixConfig)
	if err != nil {
		return nil, err
	}

	return &ContextCredentialsFactory{
		softlayerConfig:      softlayerConfig,
		vpcConfig:            vpcConfig,
		tokenExchangeService: tokenExchangeService,
	}, nil
}
