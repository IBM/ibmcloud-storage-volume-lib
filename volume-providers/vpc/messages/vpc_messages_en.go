/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Bluemix Container Registry, 5737-D42
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets,  * irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package messages

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
)

// messagesEn ...
var messagesEn = map[string]reasoncode.Message{
	"FailedToPlaceOrder": reasoncode.Message{
		Code:        "FailedToPlaceOrder",
		Description: "Failed to place storage order with the storage provider",
		Type:        "ProvisioningFailed",
		RC:          500,
		Action:      "Wait a few minutes, then try re-creating your PVC. If the problem persists, go to the IBM Cloud infrastructure (VPC) portal and open a support ticket.",
	},
	"CreateOrderTimeoutDueToPermissions": reasoncode.Message{
		Code: "CreateOrderTimeoutDueToPermissions",
		Description: "Storage with the order ID %d could not be created after retrying for %d seconds. Description: %s . " +
			"Contact your administrator to confirm that your account has the required permissions to add storage from IBM Cloud portfolio. " +
			"If your account does have the required permissions, try again. If the error persists, contact IBM Bluemix support.",
		Type:   "ProvisioningFailed",
		RC:     408,
		Action: "Run `ibmcloud ks api-key-info` to see the owner of the API key that is used to order storage. Then, contact the account administrator to check that the required permissions exist to order storage. If infrastructure credentials were manually set via `ibmcloud ks credentials-set`, check the permissions of that user. If the required permissions exist, try again. If the problem persists, go to the IBM Cloud infrastructure portal and open a support ticket.",
	},
	"StorageProvisionTimeout": reasoncode.Message{
		Code:        "StorageProvisionTimeout",
		Description: "Storage with the order ID %d could not be created after retrying for %d seconds. Description: %s ",
		Type:        "ProvisioningFailed",
		RC:          408,
		Action:      "Delete your PVC, then re-create it. If the problem persists, go to the IBM Cloud infrastructure (VPC) portal and open a support ticket.",
	},
	"FailedToDeleteVolume": reasoncode.Message{
		Code:        "FailedToDeleteVolume",
		Description: "Failed to delete the storage with storage id %d",
		Type:        "DeletionFailed",
		RC:          500,
		Action:      "If you tried to delete storage on the last day of the billing cycle, delete failures are expected and no action is required. The delete operation is automatically retried after the billing cycle ends. If the problem persists after the billing cycle is over, go to the IBM Cloud infrastructure portal and open a support ticket.",
	},
	"FailedToDeleteSnapshot": reasoncode.Message{
		Code:        "FailedToDeleteSnapshot",
		Description: "Failed to delete the storage with snapshot id %d",
		Type:        "DeletionFailed",
		RC:          500,
		Action:      "Check whether the snapshot ID exists.",
	},
	"InsufficientAuth": reasoncode.Message{
		Code:        "InsufficientAuth",
		Description: "Failed to find any valid VPC credentials in configuration file",
		Type:        "InvalidAuthentication",
		RC:          401,
		Action:      "Run `ibmcloud ks api-key-info` to see the owner of the API key that is used to order storage. Then, contact the account administrator to check that the required permissions exist to order storage. If infrastructure credentials were manually set via `ibmcloud ks credentials-set`, check the permissions of that user. If the required permissions exist, wait a few minutes, then delete the PVC and re-create it. If the problem persists, go to the IBM Cloud infrastructure portal and open a support ticket.",
	},
	"StorageFindFailedWithVolumeId": reasoncode.Message{
		Code:        "StorageFindFailedWithVolumeId",
		Description: "Failed to find the storage with storage id %s. Description: %s",
		Type:        "RetrivalFailed",
		RC:          500,
		Action:      "Go to the IBM Cloud infrastructure (VPC) portal and verify that the storage exists.",
	},
	"StorageFindFailedWithSnapshotId": reasoncode.Message{
		Code:        "StorageFindFailedWithSnapshotId",
		Description: "Failed to find the storage with storage id %s. Description: %s",
		Type:        "RetrivalFailed",
		RC:          400,
		Action:      "Go to the IBM Cloud infrastructure (VPC) portal and verify that the storage exists.",
	},
	"VolumeDoesnotHaveIops": reasoncode.Message{
		Code:        "VolumeDoesnotHaveIops",
		Description: "Volume ID '%d' does not have IOPS",
		Type:        "ProvisioningFailed",
		RC:          500,
		Action:      "Please raise a ticket to IBM Infrasructure team",
	},
	"InvalidVolumeName": reasoncode.Message{
		Code:        "InvalidVolumeName",
		Description: "Volume name is not valid.",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Volume name is not valid. Check whether the volume name is empty.",
	},
	"VolumeCapacityInvalid": reasoncode.Message{
		Code:        "VolumeCapacityInvalid",
		Description: "Volume capacity is not valid.",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Volume capacity is not valid. Check whether the volume capacity has allowed values.",
	},
	"IopsInvalid": reasoncode.Message{
		Code:        "IopsInvalid",
		Description: "Volume IOPs is not valid.",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Volume IOPs is not valid. Check whether the volume IOPs has allowed values.",
	},
	"DuplicateVolumeName": reasoncode.Message{
		Code:        "DuplicateVolumeName",
		Description: "The volume name specified in the request already exists.",
		Type:        "DuplicateVolumeName",
		RC:          404,
		Action:      "Volume name already exists. Please give different volume name.",
	},
	"InvalidCapacityIopsProfile": reasoncode.Message{
		Code:        "InvalidCapacityIopsProfile",
		Description: "Capacity or IOPS value not allowed by profile.",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Failed to place storage order with the storage provider [Backend Error:The volume profile specified in the request cannot accept custom IOPS].",
	},
	"VolumeProfileIopsInvalid": reasoncode.Message{
		Code:        "VolumeProfileIopsInvalid",
		Description: "Capacity or IOPS value not allowed by profile. The volume profile specified in the request cannot accept custom IOPS.",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Capacity or IOPS value not allowed by profile. The volume profile specified in the request cannot accept custom IOPS.",
	},
	"SnapshotSpaceOrderFailed": reasoncode.Message{
		Code:        "SnapshotSpaceOrderFailed",
		Description: "Snapshot space order failed for the given volume ID",
		Type:        "ProvisioningFailed",
		RC:          500,
		Action:      "Please check your input",
	},
}

// InitMessages ...
func InitMessages() map[string]reasoncode.Message {
	return messagesEn
}
