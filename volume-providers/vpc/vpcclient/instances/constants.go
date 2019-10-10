/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package instances

const (
	instanceIDParam                = "instance-id"
	clusterIDParam                 = "cluster-id"
	volumeIDParam                  = "volume-id"
	attachmentIDParam              = "id"
	instanceIDPath                 = "/{" + instanceIDParam + "}"
	volumeAttachmentPath           = "volume_attachments"
	instanceIDvolumeAttachmentPath = instanceIDPath + "/" + volumeAttachmentPath
	instanceIDattachmentIDPath     = instanceIDvolumeAttachmentPath + "/{" + attachmentIDParam + "}"

	// VpcPathPrefix  VPC URL path prefix
	VpcPathPrefix = "v1/instances"

	// IksPathPrefix  IKS URL path prefix
	IksPathPrefix = "v2/storage/"

	// IksClusterQueryKey ...
	IksClusterQueryKey = "cluster"

	// IksWorkerQueryKey ...
	IksWorkerQueryKey = "worker"

	// IksVolumeQueryKey ...
	IksVolumeQueryKey = "volumeID"

	// IksVolumeAttachmentIDQueryKey ...
	IksVolumeAttachmentIDQueryKey = "volumeAttachmentID"
)
