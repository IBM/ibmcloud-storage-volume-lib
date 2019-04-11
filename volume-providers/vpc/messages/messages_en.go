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

// messagesEn ...
var messagesEn = map[string]Message{
	"FailedToPlaceOrder": Message{
		Code:        "FailedToPlaceOrder",
		Description: "Failed to create volume with the storage provider",
		Type:        "ProvisioningFailed",
		RC:          500,
		Action:      "Please check the [Backend] error which is returned, You may need to try by using 'ibmcloud is volume-create' cli.",
	},
	"FailedToDeleteVolume": Message{
		Code:        "FailedToDeleteVolume",
		Description: "Failed to delete '%d' volume from VPC",
		Type:        "DeletionFailed",
		RC:          500,
		Action:      "Please check volume ID, You may need to try by using 'ibmcloud is volume-delete' or check if volume exists 'ibmcloud is volume <VolumeID>' cli",
	},
	"FailedToDeleteSnapshot": Message{
		Code:        "FailedToDeleteSnapshot",
		Description: "Failed to delete '%d' snapshot ID",
		Type:        "DeletionFailed",
		RC:          500,
		Action:      "Check whether the snapshot ID exists. You may need to verify by using 'ibmcloud is' cli",
	},
	"StorageFindFailedWithVolumeId": Message{
		Code:        "StorageFindFailedWithVolumeId",
		Description: "Failed to find '%s' volume ID.",
		Type:        "RetrivalFailed",
		RC:          404,
		Action:      "Please check volume ID if this is correct, You may need to verify by using 'ibmcloud is volume <VolumeID>' cli.",
	},
	"StorageFindFailedWithSnapshotId": Message{
		Code:        "StorageFindFailedWithSnapshotId",
		Description: "Failed to find the volume by using '%s' snapshot ID. Description: %s",
		Type:        "RetrivalFailed",
		RC:          400,
		Action:      "Please check the snapshot ID once, You many need to verify by using 'ibmcloud is' cli.",
	},
	"InvalidVolumeID": Message{
		Code:        "InvalidVolumeID",
		Description: "'%s' volume ID is not valid. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_id_invalid",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Volume ID is not valid. Please review provided URL for valid volume ID.",
	},
	"InvalidVolumeName": Message{
		Code:        "InvalidVolumeName",
		Description: "'%s' volume name is not valid. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#validation_invalid_name",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Volume name is not valid. Please check review provided URL for valid volume name.",
	},
	"VolumeCapacityInvalid": Message{
		Code:        "VolumeCapacityInvalid",
		Description: "'%d' volume capacity is not valid. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_capacity_zero_or_negative",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Volume capacity is not valid. Please review provided URL for valid volume capacity",
	},
	"IopsInvalid": Message{
		Code:        "IopsInvalid",
		Description: "'%s' volume IOPs  is not valid.Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_iops_zero_or_negative",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "Volume IOPs is not valid. Please review provided URL for valid Iops",
	},
	"VolumeProfileIopsInvalid": Message{
		Code:        "VolumeProfileIopsInvalid",
		Description: "IOPS value not allowed by profile. Please check https://cloud.ibm.com/docs/infrastructure/vpc?topic=vpc-rias-error-messages#volume_profile_iops_invalid",
		Type:        "InvalidRequest",
		RC:          400,
		Action:      "IOPS value not allowed by profile. Please review provided URL for profile and iops",
	},
	"SnapshotSpaceOrderFailed": Message{
		Code:        "SnapshotSpaceOrderFailed",
		Description: "Snapshot space order failed for the given volume ID",
		Type:        "ProvisioningFailed",
		RC:          500,
		Action:      "Please check your input",
	},
}

// InitMessages ...
func InitMessages() map[string]Message {
	return messagesEn
}
