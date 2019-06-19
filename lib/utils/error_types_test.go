/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package util

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetErrorType(t *testing.T) {
	err := errors.New("Infrastructure account is temporarily locked")
	newErr := NewError("ErrorProviderAccountTemporarilyLocked", "Infrastructure account is temporarily locked", err)
	assert.NotNil(t, GetErrorType(newErr))
	newErr = NewError("ProvisioningFailed", "ProvisioningFailed", errors.New("ProvisioningFailed"))
	assert.NotNil(t, GetErrorType(newErr))
}
