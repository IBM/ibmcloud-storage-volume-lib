/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package reasoncode

// ReasonCode ...
type ReasonCode string

// -- General error codes --

const (

	// ErrorUnclassified indicates a generic unclassified error
	ErrorUnclassified = ReasonCode("ErrorUnclassified")

	// ErrorPanic indicates recovery from a panic
	ErrorPanic = ReasonCode("ErrorPanic")

	// ErrorTemporaryConnectionProblem indicates an *AMBIGUOUS RESPONSE* due to IaaS API timeout or reset
	// (Caller can continue to retry indefinitely)
	ErrorTemporaryConnectionProblem = ReasonCode("ErrorTemporaryConnectionProblem")

	// ErrorRateLimitExceeded indicates IaaS API rate limit has been exceeded
	// (Caller can continue to retry indefinitely)
	ErrorRateLimitExceeded = ReasonCode("ErrorRateLimitExceeded")
)

// -- General provider API (RPC) errors ---

const (

	// ErrorBadRequest indicates a generic bad request to the Provider API
	// (Caller can treat this as a fatal failure)
	ErrorBadRequest = ReasonCode("ErrorBadRequest")

	// ErrorUnsupportedAuthType indicates the requested Auth-Type is not supported
	// (Caller can treat this as a fatal failure)
	ErrorUnsupportedAuthType = ReasonCode("ErrorUnsupportedAuthType")

	// ErrorUnsupportedMethod indicates the requested Provider API method is not supported
	// (Caller can treat this as a fatal failure)
	ErrorUnsupportedMethod = ReasonCode("ErrorUnsupportedMethod")
)

// -- Authentication and authorization problems --

const (

	// ErrorUnknownProvider indicates the named provider is not known
	ErrorUnknownProvider = ReasonCode("ErrorUnknownProvider")

	// ErrorUnauthorised indicates an IaaS authorisation error
	ErrorUnauthorised = ReasonCode("ErrorUnauthorised")

	// ErrorFailedTokenExchange indicates an IAM token exchange problem
	ErrorFailedTokenExchange = ReasonCode("ErrorFailedTokenExchange")

	// ErrorProviderAccountTemporarilyLocked indicates the IaaS account as it has been temporarily locked
	ErrorProviderAccountTemporarilyLocked = ReasonCode("ErrorProviderAccountTemporarilyLocked")

	// ErrorInsufficientPermissions indicates an operation failed due to a confirmed problem with IaaS user permissions
	// (Caller can retry later, but not indefinitely)
	ErrorInsufficientPermissions = ReasonCode("ErrorInsufficientPermissions")
)
