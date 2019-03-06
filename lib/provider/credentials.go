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

// AuthType ...
type AuthType string

const (
	// IaaSAPIKey is an IaaS-native user ID and API key
	IaaSAPIKey = AuthType("IAAS_API_KEY")

	// IAMAPIKey is an IAM account ID and API key
	IAMAPIKey = AuthType("IAM_API_KEY")

	// IAMAccessToken indicates the credential is an IAM access token
	IAMAccessToken = AuthType("IAM_ACCESS_TOKEN")
)

// ContextCredentials represents user credentials (e.g. API key) for volume operations from IaaS provider
type ContextCredentials struct {
	AuthType       AuthType
	DefaultAccount bool
	Region         string
	IAMAccountID   string
	UserID         string `json:"-"` // Do not trace
	Credential     string `json:"-"` // Do not trace
}
