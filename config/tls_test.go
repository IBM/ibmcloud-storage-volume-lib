/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGeneralCAHttpClient(t *testing.T) {
	t.Log("Testing GeneralCAHttpClient")

	client, _ := GeneralCAHttpClient()

	assert.NotNil(t, client)
}

func TestGeneralCAHttpClientWithTimeout(t *testing.T) {
	t.Log("Testing GeneralCAHttpClientWithTimeout")

	client, _ := GeneralCAHttpClientWithTimeout(120)

	assert.NotNil(t, client)
	assert.Equal(t, client.Timeout, time.Duration(120))
}
