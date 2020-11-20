/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Bluemix Container Registry, 5737-D42
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets,  * irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/
package messages

import (
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
)

var messages_en = map[string]util.Message{

	"E0001": {
		Code:        "E0001",
		Description: "Unable to locate datacenter with name '%s'",
		Type:        "DataCenterNotFound",
		RC:          404,
		Action:      "Contact the account administrator to verify that the owner of the API key that you use to order storage has the required permissions. If the required permissions exist, try again. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"E0002": {
		Code:        "E0002",
		Description: "Unable to find the exact ItemPriceIds(type|size|iops) for the specified storage",
		Type:        "ItemPriceIdNotFound",
		RC:          404,
		Action:      "Delete your PVC, then re-create it. If the problem persists, contact the account administrator to verify that the owner of the API key that you use to order storage has the required permissions. If the suggested actions do not solve the problem, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"E0003": {
		Code:        "E0003",
		Description: "Failed to place storage order with the storage provider",
		Type:        "StorageOrderFailed",
		RC:          500,
		Action:      "Wait a few minutes, then try re-creating your PVC. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"E0004": {
		Code: "E0004",
		Description: "Storage with the order ID %d could not be created after retrying for %d seconds. Description: %s . " +
			"Contact your administrator to confirm that your account has the required permissions to add storage from Bluemix Infrastructure portfolio. " +
			"If your account does have the required permissions, try again. If the error persists, contact IBM Bluemix support.",
		Type:   "StorageProvisionTimeoutWithUserRestriction",
		RC:     408,
		Action: "Run `bx cs api-key-info` to see the owner of the API key that is used to order storage. Then, contact the account administrator to check that the required permissions exist to order storage. If infrastructure credentials were manually set via `bx cs credentials-set`, check the permissions of that user. If the required permissions exist, try again. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"E0005": {
		Code:        "E0005",
		Description: "Storage with the order ID %d could not be created after retrying for %d seconds. Description: %s ",
		Type:        "StorageProvisionTimeout",
		RC:          408,
		Action:      "Delete your PVC, then re-create it. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"E0006": {
		Code:        "E0006",
		Description: "Failed to delete the storage with storage id %d",
		Type:        "StorageDeletionFailed",
		RC:          500,
		Action:      "If you tried to delete storage on the last day of the billing cycle, delete failures are expected and no action is required. The delete operation is automatically retried after the billing cycle ends. If the problem persists after the billing cycle is over, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"E0007": {
		Code:        "E0007",
		Description: "Failed to do subnet authorizations for the storage %d",
		Type:        "StorageSubnetAuthFailed",
		RC:          401,
		Action:      "Delete your PVC, then re-create it. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket. and post the subnetid which failed storage authorization.",
	},
	"E0008": {
		Code:        "E0008",
		Description: "Failed to do all subnets authorization for the storage %d",
		Type:        "StorageAllSubnetAuthFailed",
		RC:          401,
		Action:      "Delete your PVC, then re-create it. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket. and post the subnetid which failed storage authorization.",
	},
	"E0009": {
		Code:        "E0009",
		Description: "Failed to find any valid softlayer credentials in configuration file",
		Type:        "InsufficientAuth",
		RC:          401,
		Action:      "Run `bx cs api-key-info` to see the owner of the API key that is used to order storage. Then, contact the account administrator to check that the required permissions exist to order storage. If infrastructure credentials were manually set via `bx cs credentials-set`, check the permissions of that user. If the required permissions exist, wait a few minutes, then delete the PVC and re-create it. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"E0010": {
		Code:        "E0010",
		Description: "Failed to find the datacenter name in configuration file",
		Type:        "DataCenterNotSpecified",
		RC:          404,
		Action:      "If you specified a data center in your PVC, verify that the data center exists. Wait a few minutes, then delete the PVC and re-create it. If the problem persists, go to the IBM Cloud infrastructure (SoftLayer) portal and open a support ticket.",
	},
	"E0011": {
		Code:        "E0011",
		Description: "Failed to find the storage with storage id %d. Description: %s",
		Type:        "StorageFindFailed",
		RC:          500,
		Action:      "Go to the IBM Cloud infrastructure (SoftLayer) portal and verify that the storage exists.",
	},
	"E0012": {
		Code:        "E0012",
		Description: "Storage type is wrong or not provided , expected storage type is 'Endurance' or 'Performance'",
		Type:        "StorageTypeRequired",
		RC:          500,
		Action:      "If you use a custom storage class, verify that you defined `Endurance` or `Performance` as the type of storage that you want to provision. Delete the PVC and re-create it.",
	},
	"E0013": {
		Code:        "E0013",
		Description: "Snapshot space order failed because volume ID or snapshot size not provided",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please check your input",
	},
	"E0014": {
		Code:        "E0014",
		Description: "There is no billing information for volume ID %d. Please check if you have already deleted this volume",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please check volume looks its deleted",
	},
	"E0015": {
		Code:        "E0015",
		Description: "There is no category information for volume ID %d.",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please double check the volume if its not deleted else need to raise a ticket with IBM Infrasructure",
	},
	"E0016": {
		Code:        "E0016",
		Description: "VolumeID '%d' not ordered via storage_as_a_service or storage_service_enterprise category code",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please select any other volume which was ordered via storage_as_a_service or storage_service_enterprise category",
	},
	"E0017": {
		Code:        "E0017",
		Description: "Failed to get package details for '%s' category",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Need to raise a ticket to IBM Infrastructure team",
	},
	"E0018": {
		Code:        "E0018",
		Description: "Snapshot space cannot be ordered for this performance volume ID '%d' since it does not support Encryption at Rest.",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please select another volume which support Encryption at Rest",
	},
	"E0019": {
		Code:        "E0019",
		Description: "Could not create snapshot order as '%s' is not suppored volume type",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please select another volume",
	},
	"E0020": {
		Code:        "E0020",
		Description: "Failed to order snapshot space for volume ID '%d' of snapshot size '%d'",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please select other supported inputes",
	},
	"E0021": {
		Code:        "E0021",
		Description: "Please provide original volume details, like volume ID and snapshot ID",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0022": {
		Code:        "E0022",
		Description: "Please provide valid volume and snapshot IDs for creating volume from snapshot",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0023": {
		Code:        "E0023",
		Description: "Volume ID '%d' does not have Snapshot Capacity",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please creae snapshot space on the volume and then try again",
	},
	"E0024": {
		Code:        "E0024",
		Description: "Volume ID '%d' does not have location ID",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please raised a ticket to IBM Infrasructure team",
	},
	"E0025": {
		Code:        "E0025",
		Description: "Volume ID '%d' does not have IOPS",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please raise a ticket to IBM Infrasructure team",
	},
	"E0026": {
		Code:        "E0026",
		Description: "Volume ID '%d' does not support valid storage type i.e ENDURANCE or PERFORMANCE",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0027": {
		Code:        "E0027",
		Description: "Volume ID '%d' does not have capacity",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0028": {
		Code:        "E0028",
		Description: "Failed to order volume from original volume ID '%d' and snapshot ID '%d'",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0029": {
		Code:        "E0029",
		Description: "Failed to create snapshot for volume ID '%d'",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0030": {
		Code:        "E0030",
		Description: "Please provide valid snapshot ID",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0031": {
		Code:        "E0031",
		Description: "Failed to delete snapshot ID '%d'",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0032": {
		Code:        "E0032",
		Description: "Failed to get snapshot ID '%d' details from provider",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0033": {
		Code:        "E0033",
		Description: "Failed to get all snapshot from the IBM infrastructure a/c, Its a IBM infrastructure layer issue",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0034": {
		Code:        "E0034",
		Description: "Failed to get all snapshots from the IBM infrastructure a/c for volume ID '%d', Its a IBM infrastructure layer issue",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0035": {
		Code:        "E0035",
		Description: "Please provide valid volume ID, 0 or nil is not the valid ID",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0036": {
		Code:        "E0036",
		Description: "Not a valid volume size",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0037": {
		Code:        "E0037",
		Description: "Please provide iops value for performance storage type ordering",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0038": {
		Code:        "E0038",
		Description: "Please provide tier value for endurance storage type ordering",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0039": {
		Code:        "E0039",
		Description: "IBM lib operation time out",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Please provide correct inputs",
	},
	"E0040": {
		Code:        "E0040",
		Description: "Provisioning failed for volume order ID '%d'",
		Type:        "SnapshotSpaceOrderFailed",
		RC:          500,
		Action:      "Already deleted storage for this order, please double check on IBM infrastructure portal",
	},
}
