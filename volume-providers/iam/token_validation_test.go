/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package iam

import (
	"testing"

	"github.com/golang-jwt/jwt"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_GetIAMAccountIDFromAccessToken(t *testing.T) {
	logger, _ := zap.NewDevelopment(zap.AddCaller())

	fakeAccountID := "12345"
	fakeSigningKey := []byte("aabbccdd")

	fakeToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"account": map[string]interface{}{"bss": fakeAccountID}}).SignedString(fakeSigningKey)

	testcases := []struct {
		name              string
		token             string
		expectedAccountID string
	}{{
		name:              "fake_token",
		token:             fakeToken,
		expectedAccountID: fakeAccountID,
	}}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {

			httpSetup()

			config := config.BluemixConfig{
				IamURL:          server.URL,
				IamClientID:     "test",
				IamClientSecret: "secret",
			}

			tes, err := NewTokenExchangeService(&config)
			assert.NoError(t, err)

			accountID, err := tes.GetIAMAccountIDFromAccessToken(AccessToken{Token: testcase.token}, logger)
			assert.Equal(t, testcase.expectedAccountID, accountID)
			assert.NoError(t, err)

		})
	}
}
