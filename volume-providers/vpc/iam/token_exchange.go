/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.ibm.com/narkarum/ibmcloud-storage-volume-lib/volume-providers/vpc/config"
	"github.ibm.com/narkarum/ibmcloud-storage-volume-lib/volume-providers/vpc/reasoncode"
	"github.ibm.com/narkarum/ibmcloud-storage-volume-lib/volume-providers/vpc/util"
)

type tokenExchangeService struct {
	config     *Config
	httpClient *http.Client
}

var _ TokenExchangeService = &tokenExchangeService{}

// NewTokenExchangeService is for internal use by armada-cluster
func NewTokenExchangeService(config *Config, httpClient *http.Client) (TokenExchangeService, error) {
	return &tokenExchangeService{
		config:     config,
		httpClient: httpClient,
	}, nil
}

type tokenExchangeRequest struct {
	httpClient *http.Client
	context    context.Context
	request    impl.Request
	logger     *zap.Logger
}

type tokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ImsToken    string `json:"ims_token"`
	ImsUserID   int    `json:"ims_user_id"`
}

// ExchangeRefreshTokenForAccessToken is for internal use by armada-cluster
func (tes *tokenExchangeService) ExchangeRefreshTokenForAccessToken(ctx context.Context, refreshToken string, logger *zap.Logger) (*AccessToken, error) {
	r, err := tes.newTokenExchangeRequest(ctx, logger)
	if err != nil {
		return nil, err
	}

	r.request.AddFormField("grant_type", "refresh_token")
	r.request.AddFormField("refresh_token", refreshToken)

	return r.exchangeForAccessToken()
}

// ExchangeAccessTokenForIMSToken is for internal use by armada-cluster
func (tes *tokenExchangeService) ExchangeAccessTokenForIMSToken(ctx context.Context, accessToken AccessToken, logger *zap.Logger) (*IMSToken, error) {
	r, err := tes.newTokenExchangeRequest(ctx, logger)
	if err != nil {
		return nil, err
	}

	r.request.AddFormField("grant_type", "urn:ibm:params:oauth:grant-type:derive")
	r.request.AddFormField("response_type", "ims_portal")
	r.request.AddFormField("access_token", accessToken.Token)

	return r.exchangeForIMSToken()
}

// ExchangeIAMAPIKeyForIMSToken is for internal use by armada-cluster
func (tes *tokenExchangeService) ExchangeIAMAPIKeyForIMSToken(ctx context.Context, iamAPIKey string, logger *zap.Logger) (*IMSToken, error) {
	r, err := tes.newTokenExchangeRequest(ctx, logger)
	if err != nil {
		return nil, err
	}

	r.request.AddFormField("grant_type", "urn:ibm:params:oauth:grant-type:apikey")
	r.request.AddFormField("response_type", "ims_portal")
	r.request.AddFormField("apikey", iamAPIKey)

	return r.exchangeForIMSToken()
}

// ExchangeIAMAPIKeyForAccessToken is for internal use by armada-cluster
func (tes *tokenExchangeService) ExchangeIAMAPIKeyForAccessToken(ctx context.Context, iamAPIKey string, logger *zap.Logger) (*AccessToken, error) {
	r, err := tes.newTokenExchangeRequest(ctx, logger)
	if err != nil {
		return nil, err
	}

	r.request.AddFormField("grant_type", "urn:ibm:params:oauth:grant-type:apikey")
	r.request.AddFormField("apikey", iamAPIKey)

	return r.exchangeForAccessToken()
}

func (r *tokenExchangeRequest) exchangeForAccessToken() (*AccessToken, error) {
	iamResp, err := r.sendTokenExchangeRequest()
	if err != nil {
		return nil, err
	}
	return &AccessToken{Token: iamResp.AccessToken}, nil
}

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

func (tes *tokenExchangeService) newTokenExchangeRequest(ctx context.Context, logger *zap.Logger) (*tokenExchangeRequest, error) {
	request, err := impl.NewRequest("POST", fmt.Sprintf("%s/oidc/token", tes.config.IamURL)) // TODO Could be hardened
	if err != nil {
		return nil, err
	}

	// Set authentication
	request.BasicAuth(tes.config.IamClientID, tes.config.IamClientSecret)

	return &tokenExchangeRequest{
		httpClient: tes.httpClient,
		context:    ctx,
		request:    request,
		logger:     logger,
	}, nil
}

func (r *tokenExchangeRequest) sendTokenExchangeRequest() (response *tokenExchangeResponse, err error) {
	// Send our request
	r.logger.Debug("Sending IAM token exchange request")
	resp, err := r.request.Do(r.context, r.httpClient)
	if err != nil {
		r.logger.Error("IAM token exchange request failed", zap.Reflect("resp", resp), util.ZapError(err))

		err = util.NewError(util.ErrorReasonCode(err), "IAM token exchange request failed", err)
		return
	}
	defer resp.Close()

	statusCode := resp.StatusCode()

	// Handle success
	if statusCode == 200 {
		r.logger.Debug("IAM token exchange request successful")

		response = &tokenExchangeResponse{}
		err = resp.Decode(response)
		if err != nil {
			r.logger.Error("Failed to parse IAM token exchange success response",
				zap.Int("statusCode", statusCode), zap.Reflect("resp", resp), util.ZapError(err))

			err = util.NewError(reasoncode.ErrorFailedTokenExchange,
				fmt.Sprintf("Unexpected IAM token exchange response (%v)", statusCode), err)
			return
		}
		return
	}

	// Parse failure
	var errorValue = &struct {
		ErrorMessage string `json:"errorMessage"`
		ErrorType    string `json:"errorCode"`
		ErrorDetails string `json:"errorDetails"`
		Requirements struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		} `json:"requirements"`
	}{}
	err = resp.Decode(errorValue)
	if err != nil {
		r.logger.Error("Failed to parse IAM token exchange failure response",
			zap.Int("statusCode", statusCode), zap.Reflect("resp", resp), util.ZapError(err))

		err = util.NewError(reasoncode.ErrorFailedTokenExchange,
			fmt.Sprintf("Unexpected IAM token exchange response (%v)", statusCode), err)
		return
	}

	// Examine error message
	if errorValue.ErrorMessage != "" {
		r.logger.Error("IAM token exchange request failed with message",
			zap.Int("statusCode", statusCode),
			zap.String("ErrorMessage:", errorValue.ErrorMessage),
			zap.String("ErrorType:", errorValue.ErrorType),
			zap.Reflect("Error", errorValue))

		err = util.NewError(reasoncode.ErrorFailedTokenExchange,
			"IAM token exchange request failed: "+errorValue.ErrorMessage,
			errors.New(errorValue.ErrorDetails+" "+errorValue.Requirements.Code+": "+errorValue.Requirements.Error))

		if errorValue.Requirements.Code == "SoftLayer_Exception_User_Customer_AccountLocked" {
			err = util.NewError(reasoncode.ErrorProviderAccountTemporarilyLocked,
				"Infrastructure account is temporarily locked", err)
		}

		return
	}

	r.logger.Error("Unexpected IAM token exchange response",
		zap.Int("statusCode", statusCode), zap.Reflect("resp", resp))

	err = util.NewError(reasoncode.ErrorFailedTokenExchange,
		fmt.Sprintf("IAM token exchange request failed (%v)", statusCode))
	return
}
