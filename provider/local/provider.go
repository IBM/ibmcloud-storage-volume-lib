/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package local

import (
	"context"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"

	"go.uber.org/zap"
)

// Provider describes the contract that is implemented by an internal provider implementation
//go:generate counterfeiter -o fakes/provider.go --fake-name Provider . Provider
type Provider interface {
	// OpenSession begins and initialises a new provider session.
	// The implementation can choose to verify the credentials and return an error if they are invalid.
	// Alternatively, the implementation can choose to defer credential verification until individual
	// methods of the context are called.
	OpenSession(context.Context, provider.ContextCredentials, zap.Logger) (provider.Session, error)

	// Returns a configured ContextCredentialsFactory for this provider
	ContextCredentialsFactory(datacenter *string) (ContextCredentialsFactory, error)
}

// ContextCredentialsFactory is a factory which can generate ContextCredentials instances
//go:generate counterfeiter -o fakes/context_credentials_factory.go --fake-name ContextCredentialsFactory . ContextCredentialsFactory
type ContextCredentialsFactory interface {
	// ForIaaSAPIKey returns a config using an explicit API key for an IaaS user account
	ForIaaSAPIKey(iamAccountID, iaasUserID, iaasAPIKey string, logger zap.Logger) (provider.ContextCredentials, error)

	// ForIAMAPIKey returns a config derived from an IAM API key (if applicable)
	ForIAMAPIKey(iamAccountID, iamAPIKey string, logger zap.Logger) (provider.ContextCredentials, error)
}
