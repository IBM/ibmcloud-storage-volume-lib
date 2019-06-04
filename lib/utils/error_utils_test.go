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
	"fmt"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewError(t *testing.T) {
	testCases := []struct {
		testName        string
		errorCode       reasoncode.ReasonCode
		errorMessage    string
		wrappedMessages []string
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
			testName:        "Wrapped",
			errorCode:       reasoncode.ErrorUnclassified,
			wrappedMessages: []string{"This is a wrapped exception"},
		},
		{
			testName:        "MultiWrapped",
			errorCode:       reasoncode.ErrorUnclassified,
			wrappedMessages: []string{"This is a wrapped exception", "This is another wrapped exception"},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %v", testCase.testName), func(t *testing.T) {
			var err error
			var wrapped []error
			for _, msg := range testCase.wrappedMessages {
				wrapped = append(wrapped, errors.New(msg))
			}
			err = NewError(testCase.errorCode, testCase.errorMessage, wrapped...)
			assert.Equal(t, testCase.errorMessage, err.Error())
			perr, isPerr := err.(provider.Error)
			if assert.True(t, isPerr) {
				if testCase.errorCode == "" {
					assert.Equal(t, testCase.errorCode, perr.Fault.ReasonCode)
				} else {
					assert.Equal(t, testCase.errorCode, perr.Fault.ReasonCode)
				}
				assert.Equal(t, testCase.wrappedMessages, perr.Fault.Wrapped)
				assert.Nil(t, perr.Fault.Properties)
			}
		})
	}

	// With properties
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %v with properties", testCase.testName), func(t *testing.T) {
			var err error
			var wrapped []error
			for _, msg := range testCase.wrappedMessages {
				wrapped = append(wrapped, errors.New(msg))
			}
			err = NewErrorWithProperties(testCase.errorCode, testCase.errorMessage, map[string]string{"prop1": "val1", "prop2": "val2"}, wrapped...)
			assert.Equal(t, testCase.errorMessage, err.Error())
			perr, isPerr := err.(provider.Error)
			if assert.True(t, isPerr) {
				assert.Equal(t, testCase.errorCode, perr.Fault.ReasonCode)
				assert.Equal(t, testCase.wrappedMessages, perr.Fault.Wrapped)
				assert.Equal(t, map[string]string{"prop1": "val1", "prop2": "val2"}, perr.Fault.Properties)
			}
		})
	}

	// Don't explode on nil wrapped errors
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		NewError(reasoncode.ErrorUnclassified, "Parent",
			errors.New("1"),
			nil,
			NewError(reasoncode.ErrorUnclassified, "2",
				nil,
				errors.New("3"),
				NewError(reasoncode.ErrorUnclassified, "4",
					errors.New("5"),
					nil,
				),
			),
			NewError(reasoncode.ErrorUnclassified, "6",
				errors.New("7"),
			),
			errors.New("8"),
		).(provider.Error).Wrapped())
}

func TestNewError_ErrorDeepUnwrapString(t *testing.T) {
	assert.Equal(t, []string{},
		ErrorDeepUnwrapString(errors.New("generic")))

	assert.Equal(t, []string{},
		ErrorDeepUnwrapString(NewError("MyCode", "My message")))

	assert.Equal(t, []string{},
		ErrorDeepUnwrapString(NewError("MyCode", "My message", nil)))

	wrapped1 := errors.New("Wrapped 1")
	assert.Equal(t, []string{wrapped1.Error()},
		ErrorDeepUnwrapString(NewError("MyCode", "My message", wrapped1)))

	wrapped2 := errors.New("Wrapped 2")
	assert.Equal(t, []string{wrapped1.Error(), wrapped2.Error()},
		ErrorDeepUnwrapString(NewError("MyCode", "My message", wrapped1, wrapped2)))

	wrapped3 := NewError("MyCode", "Wrapped 3", wrapped1, nil)
	assert.Equal(t, []string{wrapped3.Error(), wrapped1.Error(), wrapped2.Error()},
		ErrorDeepUnwrapString(NewError("MyCode", "My message", wrapped3, nil, wrapped2)))
}

func TestErrorReasonCode(t *testing.T) {
	assert.Equal(t, reasoncode.ErrorUnclassified, ErrorReasonCode(errors.New("Test")))
	assert.Equal(t, reasoncode.ErrorUnclassified, ErrorReasonCode(provider.Error{}))
}

func TestErrorToFault(t *testing.T) {
	assert.Nil(t, ErrorToFault(nil))

	f := ErrorToFault(errors.New("test"))
	if assert.NotNil(t, f) {
		assert.Equal(t, reasoncode.ReasonCode(""), f.ReasonCode)
		assert.Equal(t, "test", f.Message)
	}

	f = ErrorToFault(NewError("MyCode", "My message"))
	if assert.NotNil(t, f) {
		assert.Equal(t, reasoncode.ReasonCode("MyCode"), f.ReasonCode)
		assert.Equal(t, "My message", f.Message)
	}
}

func TestFaultToError(t *testing.T) {
	assert.Nil(t, FaultToError(nil))

	e := FaultToError(&provider.Fault{
		ReasonCode: "MyCode",
		Message:    "My message",
		Wrapped:    []string{"wrapped"},
		Properties: map[string]string{"this": "that"},
	})
	if assert.Error(t, e) {
		assert.Equal(t, "My message", e.Error())
		if perr, isPerr := e.(provider.Error); assert.True(t, isPerr) {
			assert.Equal(t, reasoncode.ReasonCode("MyCode"), perr.Code())
			assert.Equal(t, []string{"wrapped"}, perr.Wrapped())
			assert.Equal(t, map[string]string{"this": "that"}, perr.Properties())
		}
	}
}

func TestSetResponseFault(t *testing.T) {
	testcases := []struct {
		name             string
		response         interface{}
		err              error
		expectedResponse interface{}
		expectedError    string
	}{
		{
			name:             "not_struct_ptr",
			response:         struct{}{},
			expectedResponse: struct{}{},
			expectedError:    "Value must be a pointer to a struct",
		}, {
			name:             "no_fault_field",
			response:         &struct{}{},
			expectedResponse: &struct{}{},
			expectedError:    "Value struct must have Fault provider.Fault field",
		}, {
			name:             "string_fault_field",
			response:         &struct{ Fault string }{},
			expectedResponse: &struct{ Fault string }{},
			expectedError:    "Value struct must have Fault provider.Fault field",
		}, {
			name:             "nil_fault",
			response:         &provider.FaultResponse{},
			expectedResponse: &provider.FaultResponse{},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {

			err := SetFaultResponse(testcase.err, testcase.response)

			assert.Equal(t, testcase.expectedResponse, testcase.response)

			if testcase.expectedError == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Equal(t, testcase.expectedError, err.Error())
			}
		})
	}
}

func TestZapError(t *testing.T) {
	assert.NotNil(t, ZapError(nil))
	assert.NotNil(t, ZapError(errors.New("Test")))
	assert.NotNil(t, ZapError(NewError("TEST", "Test")))
}
