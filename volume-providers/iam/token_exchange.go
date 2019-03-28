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
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
)

// tokenExchangeService ...
type tokenExchangeService struct {
	bluemixConf *config.BluemixConfig
	httpClient  *http.Client
}

// TokenExchangeService ...
var _ TokenExchangeService = &tokenExchangeService{}

// NewTokenExchangeServiceWithClient ...
func NewTokenExchangeServiceWithClient(bluemixConf *config.BluemixConfig, httpClient *http.Client) (TokenExchangeService, error) {
	return &tokenExchangeService{
		bluemixConf: bluemixConf,
		httpClient:  httpClient,
	}, nil
}

// NewTokenExchangeService ...
func NewTokenExchangeService(bluemixConf *config.BluemixConfig) (TokenExchangeService, error) {
	httpClient, err := config.GeneralCAHttpClient()
	if err != nil {
		return nil, err
	}
	return &tokenExchangeService{
		bluemixConf: bluemixConf,
		httpClient:  httpClient,
	}, nil
}

// tokenExchangeRequest ...
type tokenExchangeRequest struct {
	tes     *tokenExchangeService
	request *rest.Request
	client  *rest.Client
	logger  *zap.Logger
}

// tokenExchangeResponse ...
type tokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ImsToken    string `json:"ims_token"`
	ImsUserID   int    `json:"ims_user_id"`
}

// ExchangeRefreshTokenForAccessToken ...
func (tes *tokenExchangeService) ExchangeRefreshTokenForAccessToken(refreshToken string, logger *zap.Logger) (*AccessToken, error) {
	r := tes.newTokenExchangeRequest(logger)

	r.request.Field("grant_type", "refresh_token")
	r.request.Field("refresh_token", refreshToken)

	return r.exchangeForAccessToken()
}

// ExchangeAccessTokenForIMSToken ...
func (tes *tokenExchangeService) ExchangeAccessTokenForIMSToken(accessToken AccessToken, logger *zap.Logger) (*IMSToken, error) {
	r := tes.newTokenExchangeRequest(logger)

	r.request.Field("grant_type", "urn:ibm:params:oauth:grant-type:derive")
	r.request.Field("response_type", "ims_portal")
	r.request.Field("access_token", accessToken.Token)

	return r.exchangeForIMSToken()
}

// ExchangeIAMAPIKeyForIMSToken ...
func (tes *tokenExchangeService) ExchangeIAMAPIKeyForIMSToken(iamAPIKey string, logger *zap.Logger) (*IMSToken, error) {
	r := tes.newTokenExchangeRequest(logger)

	r.request.Field("grant_type", "urn:ibm:params:oauth:grant-type:apikey")
	r.request.Field("response_type", "ims_portal")
	r.request.Field("apikey", iamAPIKey)

	return r.exchangeForIMSToken()
}

// ExchangeIAMAPIKeyForAccessToken ...
func (tes *tokenExchangeService) ExchangeIAMAPIKeyForAccessToken(iamAPIKey string, logger *zap.Logger) (*AccessToken, error) {
	r := tes.newTokenExchangeRequest(logger)

	r.request.Field("grant_type", "urn:ibm:params:oauth:grant-type:apikey")
	r.request.Field("apikey", iamAPIKey)

	return r.exchangeForAccessToken()
}

// exchangeForAccessToken ...
func (r *tokenExchangeRequest) exchangeForAccessToken() (*AccessToken, error) {
	iamResp, err := r.sendTokenExchangeRequest()
	if err != nil {
		return nil, err
	}
	return &AccessToken{Token: iamResp.AccessToken}, nil
}

// exchangeForIMSToken ...
func (r *tokenExchangeRequest) exchangeForIMSToken() (*IMSToken, error) {
	iamResp, err := r.sendTokenExchangeRequest()
	if err != nil {
		return nil, err
	}
	return &IMSToken{
		UserID: iamResp.ImsUserID,
		Token:  iamResp.ImsToken,
	}, nil
}

// newTokenExchangeRequest ...
func (tes *tokenExchangeService) newTokenExchangeRequest(logger *zap.Logger) *tokenExchangeRequest {
	client := rest.NewClient()
	client.HTTPClient = tes.httpClient

	return &tokenExchangeRequest{
		tes:     tes,
		request: rest.PostRequest(fmt.Sprintf("%s/oidc/token", tes.bluemixConf.IamURL)),
		client:  client,
		logger:  logger,
	}
}

// sendTokenExchangeRequest ...
func (r *tokenExchangeRequest) sendTokenExchangeRequest() (*tokenExchangeResponse, error) {
	// Set headers
	basicAuth := fmt.Sprintf("%s:%s", r.tes.bluemixConf.IamClientID, r.tes.bluemixConf.IamClientSecret)
	r.request.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(basicAuth))))
	r.request.Set("Accept", "application/json")

	// Make the request
	var successV tokenExchangeResponse
	var errorV = struct {
		ErrorMessage string `json:"errorMessage"`
		ErrorType    string `json:"errorCode"`
		ErrorDetails string `json:"errorDetails"`
		Requirements struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		} `json:"requirements"`
	}{}

	r.logger.Info("Sending IAM token exchange request")
	resp, err := r.client.Do(r.request, &successV, &errorV)

	if err != nil {
		r.logger.Error("IAM token exchange request failed", zap.Reflect("Response", resp), zap.Error(err))

		// TODO Handle timeout here?

		return nil,
			util.NewError("ErrorUnclassified",
				"IAM token exchange request failed", err)
	}

	if resp != nil && resp.StatusCode == 200 {
		r.logger.Debug("IAM token exchange request successful")
		return &successV, nil
	}

	// TODO Check other status code values? (but be careful not to mask the reason codes, below)

	if errorV.ErrorMessage != "" {
		r.logger.Error("IAM token exchange request failed with message",
			zap.Int("StatusCode", resp.StatusCode),
			zap.String("ErrorMessage:", errorV.ErrorMessage),
			zap.String("ErrorType:", errorV.ErrorType),
			zap.Reflect("Error", errorV))

		err := util.NewError("ErrorFailedTokenExchange",
			"IAM token exchange request failed: "+errorV.ErrorMessage,
			errors.New(errorV.ErrorDetails+" "+errorV.Requirements.Code+": "+errorV.Requirements.Error))

		if errorV.Requirements.Code == "SoftLayer_Exception_User_Customer_AccountLocked" {
			err = util.NewError("ErrorProviderAccountTemporarilyLocked",
				"Infrastructure account is temporarily locked", err)
		}

		return nil, err
	}

	r.logger.Error("Unexpected IAM token exchange response",
		zap.Int("StatusCode", resp.StatusCode), zap.Reflect("Response", resp))

	return nil,
		util.NewError("ErrorUnclassified",
			"Unexpected IAM token exchange response")
}
