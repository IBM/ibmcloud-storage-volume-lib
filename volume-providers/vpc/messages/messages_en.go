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
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
)

// messagesEn ...
var messagesEn = map[string]util.Message{
	"AuthenticationFailed": util.Message{
		Code:        AuthenticationFailed,
		Description: "Failed to authenticate the user.",
		Type:        util.Unauthenticated,
		RC:          400,
		Action:      "Verify that you entered the correct IBM Cloud user name and password. If the error persists, the authentication service might be unavailable. Wait a few minutes and try again. ",
	},
	"ErrorRequiredFieldMissing": util.Message{
		Code:        "ErrorRequiredFieldMissing",
		Description: "[%s] is required to complete the operation.",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Review the error that is returned. Provide the missing information in your request and try again. ",
	},
	"FailedToPlaceOrder": util.Message{
		Code:        "FailedToPlaceOrder",
		Description: "Failed to create volume with the storage provider",
		Type:        util.ProvisioningFailed,
		RC:          500,
		Action:      "Review the error that is returned. If the volume creation service is currently unavailable, try to manually create the volume with the 'ibmcloud is volume-create' command.",
	},
	"FailedToDeleteVolume": util.Message{
		Code:        "FailedToDeleteVolume",
		Description: "The volume ID '%d' could not be deleted from your VPC.",
		Type:        util.DeletionFailed,
		RC:          500,
		Action:      "Verify that the volume ID exists. Run 'ibmcloud is volumes' to list available volumes in your account. If the ID is correct, try to delete the volume with the 'ibmcloud is volume-delete' command. ",
	},
	"FailedToUpdateVolume": util.Message{
		Code:        "FailedToUpdateVolume",
		Description: "The volume ID '%d' could not be updated",
		Type:        util.UpdateFailed,
		RC:          500,
		Action:      "Verify that the volume ID exists. Run 'ibmcloud is volumes' to list available volumes in your account.",
	},
	"FailedToDeleteSnapshot": util.Message{
		Code:        "FailedToDeleteSnapshot",
		Description: "Failed to delete '%d' snapshot ID",
		Type:        util.DeletionFailed,
		RC:          500,
		Action:      "Check whether the snapshot ID exists. You may need to verify by using 'ibmcloud is' cli",
	},
	"StorageFindFailedWithVolumeId": util.Message{
		Code:        "StorageFindFailedWithVolumeId",
		Description: "A volume with the specified volume ID '%s' could not be found.",
		Type:        util.RetrivalFailed,
		RC:          404,
		Action:      "Verify that the volume ID exists. Run 'ibmcloud is volumes' to list available volumes in your account.",
	},
	"StorageFindFailedWithVolumeName": util.Message{
		Code:        "StorageFindFailedWithVolumeName",
		Description: "A volume with the specified volume name '%s' does not exist.",
		Type:        util.RetrivalFailed,
		RC:          404,
		Action:      "Verify that the specified volume exists. Run 'ibmcloud is volumes' to list available volumes in your account.",
	},
	"StorageFindFailedWithSnapshotId": util.Message{
		Code:        "StorageFindFailedWithSnapshotId",
		Description: "No volume could be found for the specified snapshot ID '%s'. Description: %s",
		Type:        util.RetrivalFailed,
		RC:          400,
		Action:      "Please check the snapshot ID once, You many need to verify by using 'ibmcloud is' cli.",
	},
	"VolumeAttachFindFailed": util.Message{
		Code:        VolumeAttachFindFailed,
		Description: "No volume attachment could be found for the specified volume ID '%s' and instance ID '%s'.",
		Type:        util.VolumeAttachFindFailed,
		RC:          400,
		Action:      "Verify that a volume attachment for your instance exists. Run 'ibmcloud is in-vols INSTANCE_ID' to list active volume attachments for your instance ID. ",
	},
	"VolumeAttachFailed": util.Message{
		Code:        VolumeAttachFailed,
		Description: "The volume ID '%s' could not be attached to the instance ID '%s'.",
		Type:        util.AttachFailed,
		RC:          500,
		Action:      "Verify that the volume ID and instance ID exist. Run 'ibmcloud is volumes' to list available volumes, and 'ibmcloud is instances' to list available instances in your account. ",
	},
	"VolumeAttachTimedOut": util.Message{
		Code:        VolumeAttachTimedOut,
		Description: "The volume ID '%s' could not be attached to the instance ID '%s'",
		Type:        util.AttachFailed,
		RC:          500,
		Action:      "Verify that the volume ID and instance ID exist. Run 'ibmcloud is volumes' to list available volumes, and 'ibmcloud is instances' to list available instances in your account.",
	},
	"VolumeDetachFailed": util.Message{
		Code:        VolumeDetachFailed,
		Description: "The volumd ID '%s' could not be detached from the instance ID '%s'.",
		Type:        util.DetachFailed,
		RC:          500,
		Action:      "Verify that the specified instance ID has active volume attachments. Run 'ibmcloud is in-vols INSTANCE_ID' to list active volume attachments for your instance ID. ",
	},
	"VolumeDetachTimedOut": util.Message{
		Code:        VolumeDetachTimedOut,
		Description: "The volume ID '%s' could not be detached from the instance ID '%s'",
		Type:        util.DetachFailed,
		RC:          500,
		Action:      "Verify that the specified instance ID has active volume attachments. Run 'ibmcloud is in-vols INSTANCE_ID' to list active volume attachments for your instance ID.",
	},
	"InvalidVolumeID": util.Message{
		Code:        "InvalidVolumeID",
		Description: "The specified volume ID '%s' is not valid.",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Verify that the volume ID exists. Run 'ibmcloud is volumes' to list available volumes in your account.",
	},
	"InvalidVolumeName": util.Message{
		Code:        "InvalidVolumeName",
		Description: "The specified volume name '%s' is not valid. ",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Verify that the volume name exists. Run 'ibmcloud is volumes' to list available volumes in your account.",
	},
	"VolumeCapacityInvalid": util.Message{
		Code:        "VolumeCapacityInvalid",
		Description: "The specified volume capacity '%d' is not valid. ",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Verify the specified volume capacity. The volume capacity must be a positive number between 10 GB and 2000 GB. ",
	},
	"IopsInvalid": util.Message{
		Code:        "IopsInvalid",
		Description: "The specified volume IOPS '%s' is not valid for the selected volume profile. ",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Review available volume profiles and IOPS in the IBM Cloud Block Storage for VPC documentation https://cloud.ibm.com/docs/vpc-on-classic-block-storage?topic=vpc-on-classic-block-storage-block-storage-profiles.",
	},
	"VolumeProfileIopsInvalid": util.Message{
		Code:        "VolumeProfileIopsInvalid",
		Description: "The specified IOPS value is not valid for the selected volume profile. ",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Review available volume profiles and IOPS in the IBM Cloud Block Storage for VPC documentation https://cloud.ibm.com/docs/vpc-on-classic-block-storage?topic=vpc-on-classic-block-storage-block-storage-profiles.",
	},
	"EmptyResourceGroup": util.Message{
		Code:        "EmptyResourceGroup",
		Description: "Resource group information could not be found.",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Provide the name or ID of the resource group that you want to use for your volume. Run 'ibmcloud resource groups' to list the resource groups that you have access to. ",
	},
	"EmptyResourceGroupIDandName": util.Message{
		Code:        "EmptyResourceGroupIDandName",
		Description: "Resource group ID or name could not be found.",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Provide the name or ID of the resource group that you want to use for your volume. Run 'ibmcloud resource groups' to list the resource groups that you have access to.",
	},
	"SnapshotSpaceOrderFailed": util.Message{
		Code:        "SnapshotSpaceOrderFailed",
		Description: "Snapshot space order failed for the given volume ID",
		Type:        util.ProvisioningFailed,
		RC:          500,
		Action:      "Please check your input",
	},
	"VolumeNotInValidState": util.Message{
		Code:        "VolumeNotInValidState",
		Description: "Volume %s did not get valid (available) status within timeout period.",
		Type:        util.ProvisioningFailed,
		RC:          500,
		Action:      "Please check your input",
	},
	"VolumeDeletionInProgress": util.Message{
		Code:        "VolumeDeletionInProgress",
		Description: "Volume %s deletion in progress.",
		Type:        util.ProvisioningFailed,
		RC:          500,
		Action:      "Wait for volume deletion",
	},
	"ListVolumesFailed": util.Message{
		Code:        "ListVolumesFailed",
		Description: "Unable to fetch list of volumes.",
		Type:        util.RetrivalFailed,
		RC:          404,
		Action:      "Run 'ibmcloud is volumes' to list available volumes in your account.",
	},
	"InvalidListVolumesLimit": util.Message{
		Code:        "InvalidListVolumesLimit",
		Description: "The value specified in the limit parameter of the list volume call is not valid.",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Verify the limit parameter's value. The limit must be a positive number.",
	},
	"StartVolumeIDNotFound": util.Message{
		Code:        "StartVolumeIDNotFound",
		Description: "The volume ID specified in the start parameter of the list volume call '%s' could not be found.",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Please Verify that the start volume ID is correct and whether you have access to the volume ID.",
	},
}

// InitMessages ...
func InitMessages() map[string]util.Message {
	return messagesEn
}
