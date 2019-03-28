/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package auth

import (
	"go.uber.org/zap"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
)

// ForIaaSAPIKey ...
func (ccf *ContextCredentialsFactory) ForIaaSAPIKey(iamAccountID, userid, apikey string, logger *zap.Logger) (provider.ContextCredentials, error) {
	return provider.ContextCredentials{
		AuthType:     provider.IaaSAPIKey,
		IAMAccountID: iamAccountID,
		UserID:       userid,
		Credential:   apikey,
	}, nil
}
