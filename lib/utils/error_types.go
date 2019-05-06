/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package util

// These are the error types which all provider should categorize their errors
const (

	// ProvisioningFailed volume or snapshot provisioning failed
	ProvisioningFailed   = "ProvisioningFailed"

  // DeletionFailed ...
  DeletionFailed       = "DeletionFailed"

  // RetrivalFailed ...
  RetrivalFailed       = "RetrivalFailed"

  // InvalidRequest ...
  InvalidRequest       = "InvalidRequest"

  // EntityNotFound ...
  EntityNotFound       = "EntityNotFound"

  // PermissionDenied ...
  PermissionDenied     = "PermissionDenied"

  // Unauthenticated ...
  Unauthenticated      = "Unauthenticated"

	//
	ErrorTypeFailed			 = "ErrorTypeConversionFailed"
  )

func GetErrorType(err error) string {
	providerError, ok := err.(Message)
	if ok {
		return providerError.Type
	}
	return ErrorTypeFailed
}
