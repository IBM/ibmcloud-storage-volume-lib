/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package riaas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/client/fakes"
)

func TestLogin(t *testing.T) {

	client := &fakes.Client{}

	riaas := Session{
		client: client,
	}

	err := riaas.Login("token")

	if assert.Equal(t, 1, client.WithAuthTokenCallCount()) {
		assert.Equal(t, "token", client.WithAuthTokenArgsForCall(0))
	}

	assert.NoError(t, err)
}

func TestVolume(t *testing.T) {

	volume := (&Session{}).Volume()

	assert.NotNil(t, volume)
}
