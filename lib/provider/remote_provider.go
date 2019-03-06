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
	"context"
)

// RemoteProvider describes the contract that is implemented by an external IaaS provider implementation
//go:generate counterfeiter -o fakes/remote_provider.go --fake-name RemoteProvider . RemoteProvider
type RemoteProvider interface {

	// GetContext is called on each invocation to return the context in which the handler method is to be invoked.
	// The implementation can choose to verify the credentials and return an error if they are invalid.
	// Alternatively, the implementation can choose to defer credential verification until individual
	// methods of the context are called.
	GetContext(context.Context, ContextCredentials) (Context, error)
}
