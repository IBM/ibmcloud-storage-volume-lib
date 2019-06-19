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
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := &Config{
		BaseURL:       "http://gc",
		AccountID:     "test account ID",
		Username:      "tester",
		APIKey:        "tester",
		ResourceGroup: "test resource group",
		Password:      "tester",
		ContextID:     "tester",
		APIVersion:    "01-01-2019",
		HTTPClient:    nil,
	}
	assert.NotNil(t, cfg.httpClient())
	cfg.HTTPClient = &http.Client{}
	assert.NotNil(t, cfg.httpClient())
	assert.Equal(t, "http://gc", cfg.baseURL())
}
