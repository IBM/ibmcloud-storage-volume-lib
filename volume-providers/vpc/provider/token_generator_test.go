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
	"github.com/stretchr/testify/assert"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/auth"
	"testing"
)

func TestTokenGenerator(t *testing.T) {
	logger, teardown := GetTestLogger(t)
	defer teardown()

	tg := tokenGenerator{}
	assert.NotNil(t, tg)

	cf := provider.ContextCredentials{
		AuthType:     provider.IAMAccessToken,
		Credential:   TestProviderAccessToken,
		IAMAccountID: TestIKSAccountID,
	}
	signedToken, err := tg.getServiceToken(cf, *logger)
	assert.Nil(t, signedToken)
	assert.NotNil(t, err)

	tg.tokenKID = "sample_key"
	signedToken, err = tg.getServiceToken(cf, *logger)
	assert.Nil(t, signedToken)
	assert.NotNil(t, err)

	cf = provider.ContextCredentials{
		AuthType:     provider.IAMAccessToken,
		Credential:   TestProviderAccessToken,
		IAMAccountID: TestIKSAccountID,
		UserID:       TestIKSAccountID,
	}

	tg.tokenKID = "no_sample_key"
	signedToken, err = tg.getServiceToken(cf, *logger)
	assert.NotNil(t, signedToken)
	assert.Nil(t, err)

	signedToken, err = tg.getServiceToken(cf, *logger)
	assert.NotNil(t, signedToken)
	assert.Nil(t, err)

	tg.tokenKID = "no_sample_key"
	cf = provider.ContextCredentials{
		AuthType:     auth.IMSToken,
		Credential:   TestProviderAccessToken,
		IAMAccountID: TestIKSAccountID,
		UserID:       TestIKSAccountID,
	}
	signedToken, err = tg.getServiceToken(cf, *logger)
	assert.NotNil(t, signedToken)
	assert.Nil(t, err)

	tg.tokenKID = "sample_key_invalid"
	cf = provider.ContextCredentials{
		AuthType:     auth.IMSToken,
		Credential:   TestProviderAccessToken,
		IAMAccountID: TestIKSAccountID,
		UserID:       TestIKSAccountID,
	}
	signedToken, err = tg.getServiceToken(cf, *logger)
	assert.NotNil(t, signedToken)
	assert.Nil(t, err)
}
