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

	AuthenticationFailed: util.Message{
		Code:        AuthenticationFailed,
		Description: "Failed to authenticate the user",
		Type:        util.Unauthenticated,
		RC:          400,
		Action:      "Either authentication service is not working properly OR user credentials are not correct. You many need to verify by using 'ibmcloud iam' cli.",
	},
	"ErrorRequiredFieldMissing": util.Message{
		Code:        "ErrorRequiredFieldMissing",
		Description: "[%s] is required to complete the operation",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Please check the [BackendError] error which is returned",
	},
	"FailedToPlaceOrder": util.Message{
		Code:        "FailedToPlaceOrder",
		Description: "Failed to create volume with the storage provider",
		Type:        util.ProvisioningFailed,
		RC:          500,
		Action:      "Please check the [BackendError] error which is returned, You may need to try by using 'ibmcloud is volume-create' cli.",
	},
	"FailedToDeleteVolume": util.Message{
		Code:        "FailedToDeleteVolume",
		Description: "Failed to delete '%d' volume from VPC",
		Type:        util.DeletionFailed,
		RC:          500,
		Action:      "Please check volume ID, You may need to try by using 'ibmcloud is volume-delete' or check if volume exists 'ibmcloud is volume VOLUME_ID' cli",
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
		Description: "Failed to find '%s' volume ID.",
		Type:        util.RetrivalFailed,
		RC:          404,
		Action:      "Verify that the volume ID exists. Run 'ibmcloud is volumes' to list available volumes in your account.",
	},
	"StorageFindFailedWithVolumeName": util.Message{
		Code:        "StorageFindFailedWithVolumeName",
		Description: "Failed to find '%s' volume name.",
		Type:        util.RetrivalFailed,
		RC:          404,
		Action:      "Verify that the specified volume exists. Run 'ibmcloud is volumes' to list available volumes in your account.",
	},
	"StorageFindFailedWithSnapshotId": util.Message{
		Code:        "StorageFindFailedWithSnapshotId",
		Description: "Failed to find the volume by using '%s' snapshot ID. Description: %s",
		Type:        util.RetrivalFailed,
		RC:          400,
		Action:      "Please check the snapshot ID once, You many need to verify by using 'ibmcloud is' cli.",
	},
	VolumeAttachFindFailed: util.Message{
		Code:        VolumeAttachFindFailed,
		Description: "Failed to find the volume attachment by using the volume ID :'%s' and the instance ID : '%s'.",
		Type:        "RetrivalFailed",
		RC:          400,
		Action:      "Verify that the volume attachment exists. Run 'ibmcloud is in-vols INSTANCE_ID' to list volume attachments. ",
	},
	VolumeAttachFailed: util.Message{
		Code:        VolumeAttachFailed,
		Description: "Failed to attach volume :'%s' to  instance : '%s'.",
		Type:        "AttachFailed",
		RC:          500,
		Action:      "Verify that the volume ID and instance ID exist. Run 'ibmcloud is volumes' to list available volumes, and 'ibmcloud is instances' to list available instances in your account. ",
	},
	VolumeAttachTimedOut: util.Message{
		Code:        VolumeAttachTimedOut,
		Description: "Volume attach timed out. Failed to attach volume :'%s' to  instance : '%s' in %s",
		Type:        "AttachFailed",
		RC:          500,
		Action:      "Verify that the volume ID and instance ID exist. Run 'ibmcloud is volumes' to list available volumes, and 'ibmcloud is instances' to list available instances in your account.",
	},
	VolumeDetachFailed: util.Message{
		Code:        VolumeDetachFailed,
		Description: "Failed to detach volume :'%s' from  instance : '%s'.",
		Type:        "DetachFailed",
		RC:          500,
		Action:      "Please check the volume attachment once, You many need to verify by using 'ibmcloud is in-vol' cli.",
	},
	VolumeDetachTimedOut: util.Message{
		Code:        VolumeDetachTimedOut,
		Description: "Volume detach timed out. Failed to detach volume :'%s' from  instance : '%s' in %s",
		Type:        "DetachFailed",
		RC:          500,
		Action:      "Please check the volume and instance details once, You many need to verify by using 'ibmcloud is in-vol' cli.",
	},
	"InvalidVolumeID": util.Message{
		Code:        "InvalidVolumeID",
		Description: "'%s' volume ID is not valid. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_id_invalid",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Volume ID is not valid. Please review provided URL for valid volume ID.",
	},
	"InvalidVolumeName": util.Message{
		Code:        "InvalidVolumeName",
		Description: "'%s' volume name is not valid. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#validation_invalid_name",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Volume name is not valid. Please review provided URL for valid volume name.",
	},
	"VolumeCapacityInvalid": util.Message{
		Code:        "VolumeCapacityInvalid",
		Description: "'%d' volume capacity is not valid. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_capacity_zero_or_negative",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Volume capacity is not valid. Please review provided URL for valid volume capacity",
	},
	"IopsInvalid": util.Message{
		Code:        "IopsInvalid",
		Description: "'%s' volume IOPs  is not valid.Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_iops_zero_or_negative",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "Volume IOPs is not valid. Please review provided URL for valid Iops",
	},
	"VolumeProfileIopsInvalid": util.Message{
		Code:        "VolumeProfileIopsInvalid",
		Description: "IOPS value not allowed by profile. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_profile_iops_invalid",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "IOPS value not allowed by profile. Please review provided URL for profile and iops",
	},
	"EmptyResourceGroup": util.Message{
		Code:        "EmptyResourceGroup",
		Description: "Resource group information not provided",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "VPC volume is associated with resource group. Please provide either resource group ID or Name to create volume",
	},
	"EmptyResourceGroupIDandName": util.Message{
		Code:        "EmptyResourceGroupIDandName",
		Description: "Resource group ID and name are empty",
		Type:        util.InvalidRequest,
		RC:          400,
		Action:      "VPC volume is associated with resource group. Please provide either resource group ID or Name to create volume",
	},
	"SnapshotSpaceOrderFailed": util.Message{
		Code:        "SnapshotSpaceOrderFailed",
		Description: "Snapshot space order failed for the given volume ID",
		Type:        util.ProvisioningFailed,
		RC:          500,
		Action:      "Please check your input",
	},
}

// InitMessages ...
func InitMessages() map[string]util.Message {
	return messagesEn
}
