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

const (
	// AuthTypeHeader is used by the remote provider protocol
	AuthTypeHeader = "Auth-Type"

	// ContextIDHeader is used by the remote provider protocol
	ContextIDHeader = "Context-Id"

	// DefaultAccountHeader is used by the remote provider protocol
	DefaultAccountHeader = "Default-Account"

	// IAMAccountHeader is used by the remote provider protocol
	IAMAccountHeader = "Iam-Account"

	// IAMAPIKeyHeader is used by the remote provider protocol
	IAMAPIKeyHeader = "Iam-Api-Key"

	// IAMAccessTokenHeader is used by the remote provider protocol
	IAMAccessTokenHeader = "Iam-Access-Token"

	// IaaSAPIKeyHeader is used by the remote provider protocol
	IaaSAPIKeyHeader = "Iaas-Api-Key"

	// IaaSAPIUserHeader is used by the remote provider protocol
	IaaSAPIUserHeader = "Iaas-Api-User"

	// ContextIDLabel matches the logging label used by armada-api
	ContextIDLabel = "req-id"

	// IKSResourceTag can be used by provider implementations that support tagging to tag resources provisioned by IKS
	IKSResourceTag = "ibm-kubernetes-service"
)
