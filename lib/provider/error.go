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

// Error implements the error interface for a Fault.
// Most easily constructed using util.NewError() or util.NewErrorWithProperties()
type Error struct {
	// Fault ...
	Fault Fault
}

// Fault encodes a fault condition.
// Does not implement the error interface so that cannot be accidentally
// misassigned to error variables when returned in a function response.
type Fault struct {
	// Message is the fault message (required)
	Message string `json:"msg"`

	// ReasonCode is fault reason code (required)  //TODO: will have better reasoncode mechanism
	ReasonCode string `json:"code"`

	// WrappedErrors contains wrapped error messages (if applicable)
	Wrapped []string `json:"wrapped,omitempty"`

	// Properties contains diagnostic properties (if applicable)
	Properties map[string]string `json:"properties,omitempty"`
}

// FaultResponse is an optional Fault
type FaultResponse struct {
	Fault *Fault `json:"fault,omitempty"`
}

var _ error = Error{}

// Error satisfies the error contract
func (err Error) Error() string {
	return err.Fault.Message
}

// Code satisfies the legacy provider.Error interface
func (err Error) Code() string {
	if err.Fault.ReasonCode == "" {
		return ""
	}
	return err.Fault.ReasonCode
}

// Wrapped mirrors the legacy provider.Error interface
func (err Error) Wrapped() []string {
	return err.Fault.Wrapped
}

// Properties satisfies the legacy provider.Error interface
func (err Error) Properties() map[string]string {
	return err.Fault.Properties
}
