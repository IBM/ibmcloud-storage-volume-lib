/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"fmt"
	"testing"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"

	"github.com/stretchr/testify/assert"
)

func Test_Errors(t *testing.T) {

	testCases := []struct {
		testName        string
		errorCode       reasoncode.ReasonCode
		errorMessage    string
		wrappedMessages []string
		properties      map[string]string
	}{
		{
			// ErrorUnclassified - General unclassified error
			testName:  "ErrorUnclassified",
			errorCode: reasoncode.ErrorUnclassified,
		},
		{
			// Default error code
			testName:  "DefaultCode",
			errorCode: "",
		},
		{
			// Example of a specific errorCode
			testName:  "ErrorUnknownProvider",
			errorCode: reasoncode.ErrorUnknownProvider,
		},
		{
			testName:        "Wrapped",
			errorCode:       reasoncode.ErrorUnclassified,
			wrappedMessages: []string{"This is a wrapped exception"},
		},
		{
			testName:        "MultiWrapped",
			errorCode:       reasoncode.ErrorUnclassified,
			wrappedMessages: []string{"This is a wrapped exception", "This is another wrapped exception"},
		},
		{
			testName:   "Properties",
			errorCode:  reasoncode.ErrorUnclassified,
			properties: map[string]string{"prop1": "val1", "prop2": "val2"},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %v", testCase.testName), func(t *testing.T) {
			err := Error{
				Fault: Fault{
					Message:    testCase.errorMessage,
					ReasonCode: testCase.errorCode,
					Wrapped:    testCase.wrappedMessages,
					Properties: testCase.properties,
				},
			}
			assert.Equal(t, testCase.errorMessage, err.Error())
			if testCase.errorCode == "" {
				assert.Equal(t, reasoncode.ErrorUnclassified, err.Code())
			} else {
				assert.Equal(t, testCase.errorCode, err.Code())
			}
			assert.Equal(t, testCase.wrappedMessages, err.Wrapped())
			assert.Equal(t, testCase.properties, err.Properties())
		})
	}
}
