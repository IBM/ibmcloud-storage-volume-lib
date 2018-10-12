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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SafeStringValue(t *testing.T) {
	assert.Equal(t, "", SafeStringValue(nil))
	hello := "hello"
	assert.Equal(t, "hello", SafeStringValue(&hello))
}
