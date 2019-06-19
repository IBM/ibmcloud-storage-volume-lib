/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageError(t *testing.T) {
	message := Message{
		Code: "ProvisioningFailed",
		Type: "Invalid",
	}
	assert.NotNil(t, message.Error())
	assert.Equal(t, "{Code:ProvisioningFailed, Type:Invalid, Description:, BackendError:, RC:0}", message.Error())
}
