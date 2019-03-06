/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package iam

import (
	"context"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
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

			config := Config{
				IamURL:          server.URL,
				IamClientID:     "test",
				IamClientSecret: "secret",
			}

			tes, err := NewTokenExchangeService(&config, http.DefaultClient)
			assert.NoError(t, err)

			accountID, err := tes.GetIAMAccountIDFromAccessToken(context.Background(), AccessToken{Token: testcase.token}, logger)
			assert.Equal(t, testcase.expectedAccountID, accountID)
			assert.NoError(t, err)

		})
	}
}
