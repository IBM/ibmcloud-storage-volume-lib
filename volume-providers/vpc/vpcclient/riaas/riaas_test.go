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
	"bytes"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client/fakes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	client := &fakes.SessionClient{}

	riaas := Session{
		client: client,
	}

	err := riaas.Login("token")

	if assert.Equal(t, 1, client.WithAuthTokenCallCount()) {
		assert.Equal(t, "token", client.WithAuthTokenArgsForCall(0))
	}

	assert.NoError(t, err)
}

func TestNewSession(t *testing.T) {
	var b bytes.Buffer
	cfg := Config{
		BaseURL:       "http://gc",
		AccountID:     "test account ID",
		Username:      "tester",
		APIKey:        "tester",
		ResourceGroup: "test resource group",
		Password:      "tester",
		ContextID:     "tester",
		APIVersion:    "2019-06-05",
		APIGeneration: 2,
		HTTPClient:    &http.Client{},
		DebugWriter:   io.Writer(&b),
	}

	session, err := New(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, session)

	d := DefaultRegionalAPIClientProvider{}
	regionalAPI, err := d.New(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, regionalAPI)

	noAPIVerAndGen := Config{
		BaseURL:       "http://gc",
		AccountID:     "test account ID",
		Username:      "tester",
		APIKey:        "tester",
		ResourceGroup: "test resource group",
		Password:      "tester",
		ContextID:     "tester",
		HTTPClient:    &http.Client{},
		DebugWriter:   io.Writer(&b),
	}
	sessionAPI, err := New(noAPIVerAndGen)
	assert.Nil(t, err)
	assert.NotNil(t, sessionAPI)
}

func TestVolumeService(t *testing.T) {
	volumeManager := (&Session{}).VolumeService()
	assert.NotNil(t, volumeManager)
}

func TestSnapshotService(t *testing.T) {
	snapshotManager := (&Session{}).SnapshotService()
	assert.NotNil(t, snapshotManager)
}

func TestVolumeAttachService(t *testing.T) {
	volumeAttachService := (&Session{}).VolumeAttachService()
	assert.NotNil(t, volumeAttachService)
}
