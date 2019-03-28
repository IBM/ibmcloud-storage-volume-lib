/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package client_test

import (
	"testing"

	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/stretchr/testify/assert"
)

func TestParams(t *testing.T) {
	params := client.Params{
		"key":         "value",
		"another-key": "another-value",
	}
	clone := params.Copy()

	assert.Equal(t, params, clone)

}
