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
	"net/url"

	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
)

type handler interface {
	Before(request *Request) error
}

// ClientSession provides an interface for a REST API client
//go:generate counterfeiter -o fakes/client.go --fake-name ClientSession . ClientSession
type ClientSession interface {
	NewRequest(operation *Operation) *Request
	WithDebug(writer io.Writer) ClientSession
	WithAuthToken(authToken string) ClientSession
	WithPathParameter(name, value string) ClientSession
	WithQueryValue(name, value string) ClientSession
}

type client struct {
	baseURL       string
	httpClient    *http.Client
	pathParams    Params
	queryValues   url.Values
	authenHandler handler
	debugWriter   io.Writer
	resourceGroup string
	contextID     string
}

// New creates a new instance of a ClientSession
func New(baseURL string, httpClient *http.Client, contextID string) ClientSession {
	return &client{
		baseURL:       baseURL,
		httpClient:    httpClient,
		pathParams:    Params{},
		queryValues:   url.Values{"version": []string{models.APIVersion}},
		authenHandler: &authenticationHandler{},
		contextID:     contextID,
	}
}

// NewRequest creates a request and configures it with the supplied operation
func (c *client) NewRequest(operation *Operation) *Request {
	headers := http.Header{}
	headers.Set("Accept", "application/json")
	headers.Set("User-Agent", models.UserAgent)
	if c.contextID != "" {
		headers.Set("X-Request-ID", c.contextID)
	}

	// Copy the query values to a new map
	qv := url.Values{}
	for k, v := range c.queryValues {
		qv[k] = v
	}

	return &Request{
		httpClient:    c.httpClient,
		baseURL:       c.baseURL,
		operation:     operation,
		pathParams:    c.pathParams.Copy(),
		authenHandler: c.authenHandler,
		headers:       headers,
		debugWriter:   c.debugWriter,
		resourceGroup: c.resourceGroup,
		queryValues:   qv,
	}
}

// WithDebug enables debug for this ClientSession, outputting to the supplied writer
func (c *client) WithDebug(writer io.Writer) ClientSession {
	c.debugWriter = writer
	return c
}

// WithAuthToken supplies the authentication token to use for all requests made
// by this session
func (c *client) WithAuthToken(authToken string) ClientSession {
	c.authenHandler = &authenticationHandler{
		authToken: authToken,
	}
	return c
}

// WithPathParameter adds a path parameter to the request
func (c *client) WithPathParameter(name, value string) ClientSession {
	c.pathParams[name] = value
	return c
}

// WithQueryValue adds a query parameter to the request
func (c *client) WithQueryValue(name, value string) ClientSession {
	c.queryValues.Add(name, value)
	return c
}
