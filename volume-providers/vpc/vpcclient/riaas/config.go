/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package riaas

import (
	"context"
	"io"
	"net/http"
)

// Config for the Session
type Config struct {
	BaseURL       string
	AccountID     string
	Username      string
	APIKey        string
	ResourceGroup string
	Password      string
	ContextID     string

	DebugWriter io.Writer
	HTTPClient  *http.Client
	Context     context.Context
}

func (c Config) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}

	return http.DefaultClient
}

func (c Config) baseURL() string {
	return c.BaseURL
}
