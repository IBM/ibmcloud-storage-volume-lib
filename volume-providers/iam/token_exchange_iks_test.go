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
	"fmt"
	"net/http"
	"testing"

	"github.com/softlayer/softlayer-go/sl"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
)

func Test_IKSExchangeRefreshTokenForAccessToken_Success(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)
	httpSetup()

	// IAM endpoint
	mux.HandleFunc("/v1/iam/apikey",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{"token": "at_success"}`)
		},
	)

	bluemixConf := config.BluemixConfig{
		PrivateAPIRoute: server.URL,
	}

	tes, err := NewTokenExchangeIKSService(&bluemixConf)
	assert.NoError(t, err)

	r, err := tes.ExchangeRefreshTokenForAccessToken("testrefreshtoken", logger)
	assert.Nil(t, err)
	if assert.NotNil(t, r) {
		assert.Equal(t, (*r).Token, "at_success")
	}
}

func Test_IKSExchangeRefreshTokenForAccessToken_FailedDuringRequest(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()
	mux.HandleFunc("/v1/iam/apikey",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"description": "did not work",
				"code": "bad news",
				"type" : "more details",
				"incidentID" : "1000"
				}`)
		},
	)

	bluemixConf := config.BluemixConfig{
		PrivateAPIRoute: server.URL,
	}

	tes, err := NewTokenExchangeIKSService(&bluemixConf)
	assert.NoError(t, err)

	r, err := tes.ExchangeRefreshTokenForAccessToken("badrefreshtoken", logger)
	assert.Nil(t, r)
	if assert.NotNil(t, err) {
		assert.Equal(t, "IAM token exchange request failed: did not work", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorFailedTokenExchange"), util.ErrorReasonCode(err))
	}
}

func Test_IKSExchangeRefreshTokenForAccessToken_FailedDuringRequest_no_message(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()
	mux.HandleFunc("/v1/iam/apikey",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{}`)
		},
	)

	bluemixConf := config.BluemixConfig{
		PrivateAPIRoute: server.URL,
	}

	tes, err := NewTokenExchangeIKSService(&bluemixConf)
	assert.NoError(t, err)

	r, err := tes.ExchangeRefreshTokenForAccessToken("badrefreshtoken", logger)
	assert.Nil(t, r)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Unexpected IAM token exchange response", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorUnclassified"), util.ErrorReasonCode(err))
	}
}

func Test_IKSExchangeRefreshTokenForAccessToken_FailedWrongApiUrl(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/v1/iam/apikey",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{}`)
		},
	)

	bluemixConf := config.BluemixConfig{
		PrivateAPIRoute: "wrongProtocolURL",
	}

	tes, err := NewTokenExchangeIKSService(&bluemixConf)
	assert.NoError(t, err)

	r, err := tes.ExchangeRefreshTokenForAccessToken("testrefreshtoken", logger)
	assert.Nil(t, r)

	if assert.NotNil(t, err) {
		assert.Equal(t, "IAM token exchange request failed", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorUnclassified"), util.ErrorReasonCode(err))
		assert.Equal(t, []string{"Post \"wrongProtocolURL/v1/iam/apikey\": unsupported protocol scheme \"\""},
			util.ErrorDeepUnwrapString(err))
	}
}

func Test_IKSExchangeRefreshTokenForAccessToken_FailedRequesting_unclassified_error(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/v1/iam/apikey",
		func(w http.ResponseWriter, r *http.Request) {
			// Leave response empty
		},
	)

	bluemixConf := config.BluemixConfig{
		PrivateAPIRoute: server.URL,
	}

	tes, err := NewTokenExchangeService(&bluemixConf)
	assert.NoError(t, err)

	r, err := tes.ExchangeRefreshTokenForAccessToken("badrefreshtoken", logger)
	assert.Nil(t, r)

	if assert.NotNil(t, err) {
		assert.Equal(t, "IAM token exchange request failed", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorUnclassified"), util.ErrorReasonCode(err))
	}
}

func Test_IKSExchangeIAMAPIKeyForAccessToken(t *testing.T) {
	var testCases = []struct {
		name               string
		apiHandler         func(w http.ResponseWriter, r *http.Request)
		expectedToken      string
		expectedError      *string
		expectedReasonCode string
	}{
		{
			name: "client error",
			apiHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(400)
			},
			expectedError:      sl.String("IAM token exchange request failed"),
			expectedReasonCode: "ErrorUnclassified",
		},
		{
			name: "success 200",
			apiHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				fmt.Fprint(w, `{ "token": "access_token_123" }`)
			},
			expectedToken: "access_token_123",
			expectedError: nil,
		},
		{
			name: "unauthorised",
			apiHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(401)
				fmt.Fprint(w, `{"description": "not authorised",
					"code": "authorisation",
					"type" : "more details",
					"incidentID" : "1000"
					}`)
			},
			expectedError:      sl.String("IAM token exchange request failed: not authorised"),
			expectedReasonCode: "ErrorFailedTokenExchange",
		},
		{
			name: "no error message",
			apiHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(400)
				fmt.Fprint(w, `{"code" : "ErrorUnclassified",
					"incidentID" : "10000"
					}`)
			},
			expectedError:      sl.String("Unexpected IAM token exchange response"),
			expectedReasonCode: "ErrorUnclassified",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			logger := zap.New(
				zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
				zap.AddCaller(),
			)
			httpSetup()

			// ResourceController endpoint
			mux.HandleFunc("/v1/iam/apikey", testCase.apiHandler)

			bluemixConf := config.BluemixConfig{
				//IamURL: server.URL,
				PrivateAPIRoute: server.URL,
			}

			tes, err := NewTokenExchangeIKSService(&bluemixConf)
			assert.NoError(t, err)

			r, actualError := tes.ExchangeIAMAPIKeyForAccessToken("apikey1", logger)
			if testCase.expectedError == nil {
				assert.NoError(t, actualError)
				if assert.NotNil(t, r) {
					assert.Equal(t, testCase.expectedToken, r.Token)
				}
			} else {
				if assert.Error(t, actualError) {
					assert.Equal(t, *testCase.expectedError, actualError.Error())
					assert.Equal(t, reasoncode.ReasonCode(testCase.expectedReasonCode), util.ErrorReasonCode(actualError))
				}
				assert.Nil(t, r)
			}

		})
	}
}
