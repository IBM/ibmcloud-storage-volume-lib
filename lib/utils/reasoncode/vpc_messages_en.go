/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Bluemix Container Registry, 5737-D42
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets,  * irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/
package reasoncode

var messages_en = map[string]Message{

	"FailedToPlaceOrder": Message{
		Code:        "FailedToPlaceOrder",
		Description: "Failed to place storage order with the storage provider",
		Type:        "StorageOrderFailed",
		RC:          500,
		Action:      "Wait a few minutes, then try re-creating your PVC. If the problem persists, go to the IBM Cloud infrastructure (VPC) portal and open a support ticket.",
	},
	"CreateOrderTimeoutDueToPermissions": Message{
		Code: "CreateOrderTimeoutDueToPermissions",
		Description: "Storage with the order ID %d could not be created after retrying for %d seconds. Description: %s . " +
			"Contact your administrator to confirm that your account has the required permissions to add storage from Bluemix Infrastructure portfolio. " +
			"If your account does have the required permissions, try again. If the error persists, contact IBM Bluemix support.",
		Type:   "StorageProvisionTimeoutWithUserRestriction",
		RC:     408,
		Action: "Run `bx cs api-key-info` to see the owner of the API key that is used to order storage. Then, contact the account administrator to check that the required permissions exist to order storage. If infrastructure credentials were manually set via `bx cs credentials-set`, check the permissions of that user. If the required permissions exist, try again. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"StorageProvisionTimeout": Message{
		Code:        "StorageProvisionTimeout",
		Description: "Storage with the order ID %d could not be created after retrying for %d seconds. Description: %s ",
		Type:        "StorageProvisionTimeout",
		RC:          408,
		Action:      "Delete your PVC, then re-create it. If the problem persists, go to the IBM Cloud infrastructure (VPC) portal and open a support ticket.",
	},
	"FailedToDeleteVolume": Message{
		Code:        "FailedToDeleteVolume",
		Description: "Failed to delete the storage with storage id %d",
		Type:        "StorageDeletionFailed",
		RC:          500,
		Action:      "If you tried to delete storage on the last day of the billing cycle, delete failures are expected and no action is required. The delete operation is automatically retried after the billing cycle ends. If the problem persists after the billing cycle is over, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"InsufficientAuth": Message{
		Code:        "InsufficientAuth",
		Description: "Failed to find any valid softlayer credentials in configuration file",
		Type:        "InsufficientAuth",
		RC:          401,
		Action:      "Run `bx cs api-key-info` to see the owner of the API key that is used to order storage. Then, contact the account administrator to check that the required permissions exist to order storage. If infrastructure credentials were manually set via `bx cs credentials-set`, check the permissions of that user. If the required permissions exist, wait a few minutes, then delete the PVC and re-create it. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"StorageFindFailedWithVolumeId": Message{
		Code:        "StorageFindFailedWithVolumeId",
		Description: "Failed to find the storage with storage id %s. Description: %s",
		Type:        "StorageFindFailed",
		RC:          500,
		Action:      "Go to the IBM Cloud infrastructure (VPC) portal and verify that the storage exists.",
	},
	"VolumeDoesnotHaveIops": Message{
		Code:        "VolumeDoesnotHaveIops",
		Description: "Volume ID '%d' does not have IOPS",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please raise a ticket to IBM Infrasructure team",
	},
	"InvalidVolumeName": Message{
		Code:        "InvalidVolumeName",
		Description: "Volume name is not valid.",
		Type:        "InvalidVolumeName",
		RC:          500,
		Action:      "Volume name is not valid. Check whether the volume name is empty.",
	},
	"VolumeCapacityInvalid": Message{
		Code:        "VolumeCapacityInvalid",
		Description: "Volume capacity is not valid.",
		Type:        "VolumeCapacityInvalid",
		RC:          500,
		Action:      "Volume capacity is not valid. Check whether the volume capacity has allowed values.",
	},
	"IopsInvalid": Message{
		Code:        "IopsInvalid",
		Description: "Volume IOPs is not valid.",
		Type:        "IopsInvalid",
		RC:          500,
		Action:      "Volume IOPs is not valid. Check whether the volume IOPs has allowed values.",
	},
	"DuplicateVolumeName": Message{
		Code:        "DuplicateVolumeName",
		Description: "The volume name specified in the request already exists.",
		Type:        "DuplicateVolumeName",
		RC:          500,
		Action:      "Volume name already exists. Please give different volume name.",
	},
	"InvalidCapacityIopsProfile": Message{
		Code:        "InvalidCapacityIopsProfile",
		Description: "Capacity or IOPS value not allowed by profile.",
		Type:        "InvalidCapacityIopsProfile",
		RC:          500,
		Action:      "Failed to place storage order with the storage provider [Backend Error:The volume profile specified in the request cannot accept custom IOPS].",
	},
	"VolumeProfileIopsInvalid": Message{
		Code:        "VolumeProfileIopsInvalid",
		Description: "Capacity or IOPS value not allowed by profile. The volume profile specified in the request cannot accept custom IOPS.",
		Type:        "VolumeProfileIopsInvalid",
		RC:          500,
		Action:      "Capacity or IOPS value not allowed by profile. The volume profile specified in the request cannot accept custom IOPS.",
	},
	"SnapshotSpaceOrderFailed": Message{
		Code:        "SnapshotSpaceOrderFailed",
		Description: "Snapshot space order failed for the given volume ID",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please check your input",
	},
}
