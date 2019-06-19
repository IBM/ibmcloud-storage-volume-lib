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
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func init() {
	logger, _ = zap.NewDevelopment()
}

func TestNewContextCredentialsFactory(t *testing.T) {
	bluemixConfig := &config.BluemixConfig{
		IamURL:    "http://myEndpointUrl",
		IamAPIKey: "test",
	}

	softlayerConfig := &config.SoftlayerConfig{
		SoftlayerAPIKey: "test",
	}

	vpcProviderConfig := &config.VPCProviderConfig{
		EndpointURL: "http://myEndpointUrl",
	}

	contextCredentials, err := NewContextCredentialsFactory(bluemixConfig, softlayerConfig, vpcProviderConfig)
	assert.NoError(t, err)
	assert.NotNil(t, contextCredentials)
}
