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

/*func Test_ForRefreshToken(t *testing.T) {

	httpSetup()
	// IAM endpoint
	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"access_token": "at_success","refresh_token": "rt_success", "ims_user_id": 123, "ims_token": "ims_token_1"}`)
		},
	)

	refreshToken := "abc12345"
	endpointURL := "http://myEndpointUrl"

	ccf, err := NewContextCredentialsFactory(
		&config.BluemixConfig{
			IamURL: server.URL,
		},
		&config.SoftlayerConfig{
			SoftlayerEndpointURL:    "Something else",
			SoftlayerIMSEndpointURL: endpointURL,
		})
	assert.NoError(t, err)

	contextCredentials, err := ccf.ForRefreshToken(refreshToken, *logger)

	assert.Nil(t, err)

	assert.Equal(t, provider.ContextCredentials{
		AuthType:   IMSToken,
		UserID:     "123",
		Credential: "ims_token_1",
	}, contextCredentials)

}

func Test_ForIAMAPIKey(t *testing.T) {

	httpSetup()
	// IAM endpoint
	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"access_token": "at_success","refresh_token": "rt_success", "ims_user_id": 123, "ims_token": "ims_token_1"}`)
		},
	)

	iamAPIKey := "12345abc"
	endpointURL := "http://myEndpointUrl"

	ccf, err := NewContextCredentialsFactory(
		&config.BluemixConfig{
			IamURL: server.URL,
		},
		&config.SoftlayerConfig{
			SoftlayerEndpointURL:    "Something else",
			SoftlayerIMSEndpointURL: endpointURL,
		})
	assert.NoError(t, err)

	contextCredentials, err := ccf.ForIAMAPIKey("account1", iamAPIKey, *logger)

	assert.Nil(t, err)

	assert.Equal(t, provider.ContextCredentials{
		AuthType:     IMSToken,
		IAMAccountID: "account1",
		UserID:       "123",
		Credential:   "ims_token_1",
	}, contextCredentials)

}

func Test_ContextCredentialsFromRefreshToken_fail_to_retrieve_access_token_from_refresh_token(t *testing.T) {

	refreshToken := "abc12345"
	endpointURL := "http://myEndpointUrl"

	httpSetup()
	// IAM endpoint
	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			// Fail on first attempt to retrieve token
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		},
	)

	ccf, err := NewContextCredentialsFactory(
		&config.BluemixConfig{
			IamURL: server.URL,
		},
		&config.SoftlayerConfig{
			SoftlayerIMSEndpointURL: endpointURL,
		})
	assert.NoError(t, err)

	contextCredentials, err := ccf.ForRefreshToken(refreshToken, *logger)

	assert.Equal(t, provider.ContextCredentials{}, contextCredentials)

	if assert.NotNil(t, err) {
		assert.Equal(t,
			util.NewError("ErrorUnclassified", "IAM token exchange request failed",
				&rest.ErrorResponse{Message: http.StatusText(http.StatusBadRequest) + "\n", StatusCode: http.StatusBadRequest}),
			err)
	}
}

func Test_ContextCredentialsFromRefreshToken_fail_to_retrieve_ims_token_from_refresh_token(t *testing.T) {

	refreshToken := "abc12345"
	endpointURL := "http://myEndpointUrl"

	httpSetup()
	// IAM endpoint
	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			// Only return a response when the access token is empty
			if (*r).FormValue("access_token") == "" {
				fmt.Fprint(w, `{"access_token": "at_success","refresh_token": "rt_success", "ims_user_id": 123, "ims_token": "ims_token_1"}`)
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
		},
	)

	ccf, err := NewContextCredentialsFactory(
		&config.BluemixConfig{
			IamURL: server.URL,
		},
		&config.SoftlayerConfig{
			SoftlayerIMSEndpointURL: endpointURL,
		})
	assert.NoError(t, err)

	contextCredentials, err := ccf.ForRefreshToken(refreshToken, logger)

	assert.Equal(t, provider.ContextCredentials{}, contextCredentials)

	if assert.NotNil(t, err) {
		assert.Equal(t,
			util.NewError("ErrorUnclassified", "IAM token exchange request failed",
				&rest.ErrorResponse{Message: http.StatusText(http.StatusBadRequest) + "\n", StatusCode: http.StatusBadRequest}),
			err)
	}
}

func Test_ContextCredentialsFromIAMAPIKey_fail_to_retrieve_ims_token_from_api_key(t *testing.T) {

	iamAPIKey := "12345abc"
	endpointURL := "http://myEndpointUrl"

	httpSetup()
	// IAM endpoint
	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		},
	)

	ccf, err := NewContextCredentialsFactory(
		&config.BluemixConfig{
			IamURL: server.URL,
		},
		&config.SoftlayerConfig{
			SoftlayerIMSEndpointURL: endpointURL,
		})
	assert.NoError(t, err)

	contextCredentials, err := ccf.ForIAMAPIKey("account1", iamAPIKey, logger)

	assert.Equal(t, provider.ContextCredentials{}, contextCredentials)

	if assert.NotNil(t, err) {
		assert.Equal(t,
			util.NewError("ErrorUnclassified", "IAM token exchange request failed",
				&rest.ErrorResponse{Message: http.StatusText(http.StatusBadRequest) + "\n", StatusCode: http.StatusBadRequest}),
			err)
	}
}

func Test_ContextCredentialsFromIAMAPIKey_fail_to_retrieve_ims_token_account_locked(t *testing.T) {

	iamAPIKey := "12345abc"
	endpointURL := "http://myEndpointUrl"

	httpSetup()
	// IAM endpoint
	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `
				{"errorMessage": "OpenID Connect exception",
				"errorCode": "BXNIM0400E",
				"errorDetails" : "Failed external authentication.",
				"requirements" : { "error": "Account has been locked for 30 minutes", "code":"SoftLayer_Exception_User_Customer_AccountLocked" }
				}`)
		},
	)

	ccf, err := NewContextCredentialsFactory(
		&config.BluemixConfig{
			IamURL: server.URL,
		},
		&config.SoftlayerConfig{
			SoftlayerIMSEndpointURL: endpointURL,
		})
	assert.NoError(t, err)

	contextCredentials, err := ccf.ForIAMAPIKey("account1", iamAPIKey, logger)

	assert.Equal(t, provider.ContextCredentials{}, contextCredentials)

	if assert.NotNil(t, err) {
		assert.Equal(t,
			util.NewError("ErrorProviderAccountTemporarilyLocked",
				"Infrastructure account is temporarily locked",
				util.NewError("ErrorFailedTokenExchange", "IAM token exchange request failed: OpenID Connect exception",
					errors.New("Failed external authentication. SoftLayer_Exception_User_Customer_AccountLocked: Account has been locked for 30 minutes")),
			),
			err)
	}
}*/
