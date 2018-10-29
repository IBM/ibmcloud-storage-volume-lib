/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package auth

import (
	"strconv"

	"go.uber.org/zap"

	"github.com/IBM/ibmcloud-storage-volume-lib/provider/auth"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
)

const (
	// IMSToken is an IMS user ID and token
	IMSToken = provider.AuthType("IMS_TOKEN")
)

// ForRefreshToken ...
func (ccf *contextCredentialsFactory) ForRefreshToken(refreshToken string, logger *zap.Logger) (provider.ContextCredentials, error) {

	accessToken, err := ccf.tokenExchangeService.ExchangeRefreshTokenForAccessToken(refreshToken, logger)
	if err != nil {
		// Must preserve provider error code in the ErrorProviderAccountTemporarilyLocked case
		logger.Error("Unable to retrieve access token from refresh token", local.ZapError(err))
		return provider.ContextCredentials{}, err
	}

	imsToken, err := ccf.tokenExchangeService.ExchangeAccessTokenForIMSToken(*accessToken, logger)
	if err != nil {
		// Must preserve provider error code in the ErrorProviderAccountTemporarilyLocked case
		logger.Error("Unable to retrieve IAM token from access token", local.ZapError(err))
		return provider.ContextCredentials{}, err
	}

	return forIMSToken("", imsToken), nil
}

// ForIAMAPIKey ...
func (ccf *contextCredentialsFactory) ForIAMAPIKey(iamAccountID, apiKey string, logger *zap.Logger) (provider.ContextCredentials, error) {

	imsToken, err := ccf.tokenExchangeService.ExchangeIAMAPIKeyForIMSToken(apiKey, logger)
	if err != nil {
		// Must preserve provider error code in the ErrorProviderAccountTemporarilyLocked case
		logger.Error("Unable to retrieve IMS credentials from IAM API key", local.ZapError(err))
		return provider.ContextCredentials{}, err
	}

	return forIMSToken(iamAccountID, imsToken), nil
}

func forIMSToken(iamAccountID string, imsToken *auth.IMSToken) provider.ContextCredentials {
	return provider.ContextCredentials{
		AuthType:     IMSToken,
		IAMAccountID: iamAccountID,
		UserID:       strconv.Itoa(imsToken.UserID),
		Credential:   imsToken.Token,
	}
}
