/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package client

import (
	"errors"
)

// ErrAuthenticationRequired is returned if a request is made before an authentication
// token has been provided to the client
var ErrAuthenticationRequired = errors.New("Authentication token required")

type authenticationHandler struct {
	authToken     string
	resourceGroup string
}

// Before is called before each request
func (a *authenticationHandler) Before(request *Request) error {
	request.resourceGroup = a.resourceGroup

	if a.authToken == "" {
		return ErrAuthenticationRequired
	}
	request.headers.Set("Authorization", "Bearer "+a.authToken)
	return nil
}
