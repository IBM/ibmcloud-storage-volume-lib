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
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"reflect"
)

// NewError returns an error that is implemented by provider.Error.
// If optional wrapped errors are a provider.Error, this preserves all child wrapped
// errors in depth-first order.
func NewError(code reasoncode.ReasonCode, msg string, wrapped ...error) error {
	return NewErrorWithProperties(code, msg, nil, wrapped...)
}

// NewErrorWithProperties returns an error that is implemented provider.Error and
// which is decorated with diagnostic properties.
// If optional wrapped errors are a provider.Error, this preserves all child wrapped
// errors in depth-first order.
func NewErrorWithProperties(code reasoncode.ReasonCode, msg string, properties map[string]string, wrapped ...error) error {
	if code == "" {
		code = "" // TODO: ErrorUnclassified
	}
	var werrs []string
	for _, w := range wrapped {
		if w != nil {
			werrs = append(werrs, w.Error())
			if p, isPerr := w.(provider.Error); isPerr {
				for _, u := range p.Wrapped() {
					werrs = append(werrs, u)
				}
			}
		}
	}
	return provider.Error{
		Fault: provider.Fault{
			ReasonCode: code,
			Message:    msg,
			Properties: properties,
			Wrapped:    werrs,
		},
	}
}

// ErrorDeepUnwrapString returns the full list of unwrapped error strings
// Returns empty slice if not a provider.Error
func ErrorDeepUnwrapString(err error) []string {
	if perr, isPerr := err.(provider.Error); isPerr && perr.Wrapped() != nil {
		return perr.Wrapped()
	}
	return []string{}
}

// ErrorReasonCode returns reason code if a provider.Error, else ErrorUnclassified
func ErrorReasonCode(err error) reasoncode.ReasonCode {
	if pErr, isPerr := err.(provider.Error); isPerr {
		if code := pErr.Code(); code != "" {
			return code
		}
	}
	return reasoncode.ErrorUnclassified
}

// ErrorToFault returns or builds a Fault pointer for an error (e.g. for a response object)
// Returns nil if no error,
func ErrorToFault(err error) *provider.Fault {
	if err == nil {
		return nil
	}
	if pErr, isPerr := err.(provider.Error); isPerr {
		return &pErr.Fault
	}
	return &provider.Fault{
		ReasonCode: "", // TODO: ErrorUnclassified,
		Message:    err.Error(),
	}
}

// FaultToError builds a Error from a Fault pointer (e.g. from a response object)
// Returns nil error if no Fault.
func FaultToError(fault *provider.Fault) error {
	if fault == nil {
		return nil
	}
	return provider.Error{Fault: *fault}
}

// SetResponseFault sets the Fault field of any response struct
func SetResponseFault(fault error, response interface{}) error {
	value := reflect.ValueOf(response)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return errors.New("Value must be a pointer to a struct")
	}
	field := value.Elem().FieldByName("Fault")
	if field.Kind() != reflect.Ptr {
		return errors.New("Value struct must have Fault provider.Fault field")
	}
	field.Set(reflect.ValueOf(ErrorToFault(fault)))
	return nil
}

// ZapError returns a zapcore.Field for an error that includes the metadata
// associated with a provider.Error. If the error is not a provider.Error then
// the standard zap.Error is used.
func ZapError(err error) zapcore.Field {
	if perr, isPerr := err.(provider.Error); isPerr {
		// Use zap.Relfect() to format all fields of struct
		// zap.Any() would select standard error formatting
		return zap.Reflect("error", perr)
	}

	return zap.Error(err)
}
