/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2020 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package iam

import (
	"errors"
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"go.uber.org/zap"
	"time"
)

// tokenExchangeIKSService ...
type tokenExchangeIKSService struct {
	tokenExchangeService
}

// TokenExchangeService ...
var _ TokenExchangeService = &tokenExchangeIKSService{}

// NewTokenExchangeIKSService ...
func NewTokenExchangeIKSService(bluemixConf *config.BluemixConfig) (TokenExchangeService, error) {
	httpClient, err := config.GeneralCAHttpClient()
	if err != nil {
		return nil, err
	}
	return &tokenExchangeIKSService{
		tokenExchangeService: tokenExchangeService{
			bluemixConf: bluemixConf,
			httpClient:  httpClient,
		},
	}, nil
}

// tokenExchangeIKSRequest ...
type tokenExchangeIKSRequest struct {
	tes          *tokenExchangeIKSService
	request      *rest.Request
	client       *rest.Client
	logger       *zap.Logger
	errorRetrier *util.ErrorRetrier
}

// tokenExchangeIKSResponse ...
type tokenExchangeIKSResponse struct {
	AccessToken string `json:"token"`
	//ImsToken    string `json:"ims_token"`
}

// ExchangeRefreshTokenForAccessToken ...
func (tes *tokenExchangeIKSService) ExchangeRefreshTokenForAccessToken(refreshToken string, logger *zap.Logger) (*AccessToken, error) {
	r := tes.newTokenExchangeRequest(logger)
	return r.exchangeForAccessToken()
}

// ExchangeIAMAPIKeyForAccessToken ...
func (tes *tokenExchangeIKSService) ExchangeIAMAPIKeyForAccessToken(iamAPIKey string, logger *zap.Logger) (*AccessToken, error) {
	r := tes.newTokenExchangeRequest(logger)
	return r.exchangeForAccessToken()
}

// newTokenExchangeRequest ...
func (tes *tokenExchangeIKSService) newTokenExchangeRequest(logger *zap.Logger) *tokenExchangeIKSRequest {
	client := rest.NewClient()
	client.HTTPClient = tes.httpClient
	retyrInterval, _ := time.ParseDuration("3s")
	return &tokenExchangeIKSRequest{
		tes:          tes,
		request:      rest.PostRequest(fmt.Sprintf("%s/v1/iam/apikey", tes.bluemixConf.PrivateAPIRoute)),
		client:       client,
		logger:       logger,
		errorRetrier: util.NewErrorRetrier(40, retyrInterval, logger),
	}
}

func (tes *tokenExchangeIKSService) GetIAMAccountIDFromAccessToken(accessToken AccessToken, logger *zap.Logger) (accountID string, err error) {
	return "Not required to implement", nil
}

// exchangeForAccessToken ...
func (r *tokenExchangeIKSRequest) exchangeForAccessToken() (*AccessToken, error) {
	var iamResp *tokenExchangeIKSResponse
	var err error
	err = r.errorRetrier.ErrorRetry(func() (error, bool) {
		iamResp, err = r.sendTokenExchangeRequest()
		return err, !isConnectionError(err) // Skip retry if its not connection error
	})
	if err != nil {
		return nil, err
	}
	return &AccessToken{Token: iamResp.AccessToken}, nil
}

// sendTokenExchangeRequest ...
func (r *tokenExchangeIKSRequest) sendTokenExchangeRequest() (*tokenExchangeIKSResponse, error) {
	r.logger.Info("In tokenExchangeIKSRequest's sendTokenExchangeRequest()")
	// Set headers
	r.request = r.request.Add("X-CSRF-TOKEN", r.tes.bluemixConf.CSRFToken)
	// Setting body
	var apikey = struct {
		APIKey string `json:"apikey"`
	}{
		APIKey: r.tes.bluemixConf.IamAPIKey,
	}
	r.request = r.request.Body(&apikey)

	var successV tokenExchangeIKSResponse
	var errorV = struct {
		ErrorCode        string `json:"code"`
		ErrorDescription string `json:"description"`
		ErrorType        string `json:"type"`
		IncidentID       string `json:"incidentID"`
	}{}

	r.logger.Info("Sending IAM token exchange request to container api server")
	resp, err := r.client.Do(r.request, &successV, &errorV)
	if err != nil {
		r.logger.Error("IAM token exchange request failed", zap.Reflect("Response", resp), zap.Error(err))
		return nil,
			util.NewError("ErrorUnclassified",
				"IAM token exchange request failed", err)
	}

	if resp != nil && resp.StatusCode == 200 {
		r.logger.Debug("IAM token exchange request successful")
		return &successV, nil
	}

	if errorV.ErrorDescription != "" {
		r.logger.Error("IAM token exchange request failed with message",
			zap.Int("StatusCode", resp.StatusCode), zap.Reflect("API IncidentID", errorV.IncidentID),
			zap.Reflect("Error", errorV))

		err := util.NewError("ErrorFailedTokenExchange",
			"IAM token exchange request failed: "+errorV.ErrorDescription,
			errors.New(errorV.ErrorCode+" "+errorV.ErrorType+", Description: "+errorV.ErrorDescription+", API IncidentID:"+errorV.IncidentID))
		return nil, err
	}

	r.logger.Error("Unexpected IAM token exchange response",
		zap.Int("StatusCode", resp.StatusCode), zap.Reflect("Response", resp))

	return nil,
		util.NewError("ErrorUnclassified",
			"Unexpected IAM token exchange response")
}
