/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package models

const (
	// APIVersion is the target RIaaS API spec version
	APIVersion = "2019-07-02"

	// APIGeneration ...
	APIGeneration = 1

	// UserAgent identifies IKS to the RIaaS API
	UserAgent = "IBM-Kubernetes-Service"

	// GTypeClassic ...
	GTypeClassic = "gc"

	// GTypeClassicDevicePrefix ...
	GTypeClassicDevicePrefix = "/dev/"

	// GTypeG2 ...
	GTypeG2 = "g2"

	// GTypeG2DevicePrefix ...
	GTypeG2DevicePrefix = "/dev/disk/by-id/virtio-"

	// VolumeAttached ...
	VolumeAttached = "attached"
)
