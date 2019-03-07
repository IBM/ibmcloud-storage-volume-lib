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

// -- General provider API
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

// -- Provider operations problems --

const (

	// ErrorInvalidDurableID indicates the format of a supplied DurableID is not understood by the provider
	// (Caller can treat this as a fatal failure)
	ErrorInvalidDurableID = ReasonCode("ErrorInvalidDurableID") // *** REPLACES ErrorInvalidOrderID ***

	// ErrorInvalidInstanceID indicates the format of a supplied InstanceID is not understood by the provider
	// (Caller can treat this as a fatal failure)
	ErrorInvalidInstanceID = ReasonCode("ErrorInvalidInstanceID") // *** REPLACES ErrorInvalidMachineID ***

	// ErrorInvalidNetworkID indicates the format of a supplied network ID is not understood by the provider
	// (Caller can treat this as a fatal failure)
	ErrorInvalidNetworkID = ReasonCode("ErrorInvalidNetworkID") // *** REPLACES ErrorInvalidPublicVLANID AND ErrorInvalidPrivateVLANID ***

	// ErrorInstanceNotFound indicates the machine specified by a valid InstanceID cannot be found
	// e.g. the machine doesn't exist or unconfirmed problem with IaaS user permissions
	// (Caller can retry later, but not indefinitely)
	ErrorInstanceNotFound = ReasonCode("ErrorInstanceNotFound") // *** REPLACES ErrorMachineNotFound ***

	// ErrorProvisionStatusNotFound indicates the provision status specified by a valid DurableID cannot be found
	// e.g. doesn't exist or unconfirmed problem with IaaS user permissions
	// (Caller can retry later, but not indefinitely)
	ErrorProvisionStatusNotFound = ReasonCode("ErrorProvisionStatusNotFound") // *** REPLACES ErrorOrderNotFound ***

	// ErrorVPCNotFound indicates the VPC specified by a valid ID cannot be found
	// e.g. doesn't exist or unconfirmed problem with IaaS user permissions
	// (Caller can retry later, but not indefinitely)
	ErrorVPCNotFound = ReasonCode("ErrorVPCNotFound")

	// ErrorSubnetNotFound indicates the Subnet specified by a valid ID cannot be found
	// e.g. doesn't exist or unconfirmed problem with IaaS user permissions
	// (Caller can retry later, but not indefinitely)
	ErrorSubnetNotFound = ReasonCode("ErrorSubnetNotFound")

	// ErrorNotReadyForOperation indicates the resource is not yet in a state where this operation can be performed
	// (Caller can continue to retry indefinitely)
	ErrorNotReadyForOperation = ReasonCode("ErrorNotReadyForOperation")

	// ErrorInvalidMachineConfig indicates the specific worker machine configuration is generally invalid or incomplete
	// e.g. a specified network ID does not exist
	// (Caller can treat this as a fatal failure)
	ErrorInvalidMachineConfig = ReasonCode("ErrorInvalidMachineConfig")

	// ErrorUnreconciledMachineConfig indicates the worker machine configuration could not be reconciled
	// due to a general configuration problem that will affect other workers of that Flavor or Zone
	// e.g, missing/invalid mandatory configuration in Flavor/Zone, or missing/ambiguous prices in SoftLayer package
	// ** CAUSES AN IMMEDIATE ALERT **
	// (Caller can retry later, but not indefinitely)
	ErrorUnreconciledMachineConfig = ReasonCode("ErrorUnreconciledMachineConfig")

	// ErrorUnreconciledOSImage indicates the OS image GUID could not be found in the IaaS image/template registry
	// ** CAUSES AN IMMEDIATE ALERT **
	// (Caller can retry later, but not indefinitely)
	ErrorUnreconciledOSImage = ReasonCode("ErrorUnreconciledOSImage")

	// ErrorInsufficientResources indicates there are insufficient IaaS resources to fulfill the provisioning request
	// (Caller can treat this as a fatal failure)
	ErrorInsufficientResources = ReasonCode("ErrorInsufficientResources")

	// ErrorInvalidOrderLocation indicates that provisioning order cannot be satisfied in the requested zone
	// TODO Deprecate into ErrorUnreconciledMachineConfig?
	// ** CAUSES AN IMMEDIATE ALERT **
	// (Caller can retry later, but not indefinitely)
	ErrorInvalidOrderLocation = ReasonCode("ErrorInvalidOrderLocation")
)

// -- Internal armada-cluster reason codes --

const (

	// ErrorInvalidProvisionedConfiguration indicates that provisioned worker differs from that specified in the model
	// e.g. connected to incorrect networks
	// TODO Deprecate or define elsewhere?
	ErrorInvalidProvisionedConfiguration = ReasonCode("ErrorInvalidProvisionedConfiguration")

	// ErrorInsufficientAuthentication indicates the supplied credentials are incomplete
	// TODO Deprecate or define elsewhere?
	ErrorInsufficientAuthentication = ReasonCode("ErrorInsufficientAuthentication")
)
