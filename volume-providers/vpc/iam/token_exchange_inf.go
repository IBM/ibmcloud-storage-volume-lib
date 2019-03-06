/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package iam

import (
	"context"

	"go.uber.org/zap"
)

// Config is for internal use by armada-cluster
type Config struct {
	IamURL          string `envconfig:"IAM_URL"`
	IamClientID     string `envconfig:"IAM_CLIENT_ID"`
	IamClientSecret string `envconfig:"IAM_CLIENT_SECRET"`
}

// IMSToken is for internal use by armada-cluster
type IMSToken struct {
	UserID int    // Numerical ID is safe to trace
	Token  string `json:"-"` // Do not trace
}

// AccessToken is for internal use by armada-cluster
type AccessToken struct {
	Token string `json:"-"` // Do not trace
}

// TokenExchangeService is for internal use by armada-cluster
//go:generate counterfeiter -o fakes/token_exchange_service.go --fake-name TokenExchangeService . TokenExchangeService
type TokenExchangeService interface {

	// ExchangeRefreshTokenForAccessToken is for internal use by armada-cluster
	// TODO Deprecate when no longer reliant on refresh token authentication
	ExchangeRefreshTokenForAccessToken(ctx context.Context, refreshToken string, logger *zap.Logger) (*AccessToken, error)

	// ExchangeAccessTokenForIMSToken is for internal use by armada-cluster
	ExchangeAccessTokenForIMSToken(ctx context.Context, accessToken AccessToken, logger *zap.Logger) (*IMSToken, error)

	// ExchangeIAMAPIKeyForIMSToken is for internal use by armada-cluster
	ExchangeIAMAPIKeyForIMSToken(ctx context.Context, iamAPIKey string, logger *zap.Logger) (*IMSToken, error)

	// ExchangeIAMAPIKeyForAccessToken is for internal use by armada-cluster
	ExchangeIAMAPIKeyForAccessToken(ctx context.Context, iamAPIKey string, logger *zap.Logger) (*AccessToken, error)
}
