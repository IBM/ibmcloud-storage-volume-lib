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
	"go.uber.org/zap"
)

// IMSToken ...
type IMSToken struct {
	UserID int    // Numerical ID is safe to trace
	Token  string `json:"-"` // Do not trace
}

// AccessToken ...
type AccessToken struct {
	Token string `json:"-"` // Do not trace
}

// TokenExchangeService ...
type TokenExchangeService interface {

	// ExchangeRefreshTokenForAccessToken ...
	// TODO Deprecate when no longer reliant on refresh token authentication
	ExchangeRefreshTokenForAccessToken(refreshToken string, logger *zap.Logger) (*AccessToken, error)

	// ExchangeAccessTokenForIMSToken ...
	ExchangeAccessTokenForIMSToken(accessToken AccessToken, logger *zap.Logger) (*IMSToken, error)

	// ExchangeIAMAPIKeyForIMSToken ...
	ExchangeIAMAPIKeyForIMSToken(iamAPIKey string, logger *zap.Logger) (*IMSToken, error)

	// ExchangeIAMAPIKeyForAccessToken ...
	ExchangeIAMAPIKeyForAccessToken(iamAPIKey string, logger *zap.Logger) (*AccessToken, error)

	// GetIAMAccountIDFromAccessToken ...
	GetIAMAccountIDFromAccessToken(accessToken AccessToken, logger *zap.Logger) (string, error)
}
