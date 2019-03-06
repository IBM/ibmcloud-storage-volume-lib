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
	"io"
	"net/http"
)

type handler interface {
	Before(request *Request) error
}

// Client provides an interface for a REST API client
//go:generate counterfeiter -o fakes/client.go --fake-name Client . Client
type Client interface {
	NewRequest(operation *Operation) *Request
	WithDebug(writer io.Writer) Client
	WithAuthToken(authToken string) Client
	WithPathParameter(name, value string) Client
}

type client struct {
	baseURL       string
	httpClient    *http.Client
	pathParams    Params
	headers       http.Header
	authenHandler handler
	debugWriter   io.Writer
	resourceGroup string
}

// New creates a new instance of a Client
func New(baseURL string, httpClient *http.Client) Client {
	return &client{
		baseURL:       baseURL,
		httpClient:    httpClient,
		pathParams:    Params{},
		authenHandler: &authenticationHandler{},
	}
}

// NewRequest creates a request and configures it with the supplied operation
func (c *client) NewRequest(operation *Operation) *Request {
	if c.headers == nil {
		c.headers = http.Header{}
	}
	c.headers.Set("Accept", "application/json")
	return &Request{
		httpClient:    c.httpClient,
		baseURL:       c.baseURL,
		operation:     operation,
		pathParams:    c.pathParams.Copy(),
		authenHandler: c.authenHandler,
		headers:       c.headers,
		debugWriter:   c.debugWriter,
		resourceGroup: c.resourceGroup,
	}
}

// WithDebug enables debug for this Client, outputting to the supplied writer
func (c *client) WithDebug(writer io.Writer) Client {
	c.debugWriter = writer
	return c
}

// WithAuthToken supplies the authentication token to use for all requests made
// by this client
func (c *client) WithAuthToken(authToken string) Client {
	c.authenHandler = &authenticationHandler{
		authToken: authToken,
	}
	return c
}

// WithPathParameter adds a path parameter to the request
func (c *client) WithPathParameter(name, value string) Client {
	c.pathParams[name] = value
	return c
}
