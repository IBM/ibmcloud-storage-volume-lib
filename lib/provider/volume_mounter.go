/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/
package provider

import ()

const (
	SUCCESS      = "Success"
	FAILURE      = "Failure"
	NOTSUPPORTED = "Not supported"
)

type Plugin interface {

	//Init method is to initialize the volume, it is a no op right now
	Init() VolumeResponse

	//Attach method attaches a volume/ fileset to a pod
	Attach(attachRequest VolumeAttachRequest) VolumeResponse

	//Detach detaches the volume/ fileset from the pod
	Detach(detachRequest VolumeDetachRequest) VolumeResponse

	//Mount method allows to mount the volume/fileset to a given location for a pod
	Mount(mountRequest VolumeMountRequest) VolumeResponse

	//Unmount methods unmounts the volume/ fileset from the pod
	Unmount(unmountRequest VolumeUnmountRequest) VolumeResponse

	//Gets the volume name from the pod
	GetVolumeName(request map[string]string) VolumeResponse

	//Wait for the volume to be attached on the node
	WaitForAttach(devicePath string, opts map[string]string) VolumeResponse

	//Wait for the volume to be detached on the node
	WaitForDetach(devicePath string) VolumeResponse

	//Checks if the volume is attached to the node
	IsAttached(request map[string]string, nodeName string) VolumeResponse

	//Mounts the device to a global path which individual pods can then bind mount
	MountDevice(deviceMountPath string, devicePath string, opts map[string]string) VolumeResponse

	//Unmounts the global mount for the device. This is called once all bind mounts have been unmounted
	UnmountDevice(deviceMountPath string) VolumeResponse
}

type VolumeResponse struct {
	// Status should be either "Success", "Failure" or "Not supported".
	Status string `json:"status"`
	// Reason for success or failure.
	Message string `json:"message,omitempty"`
	// Path to the device attached. This field is valid only for attach calls.
	// ie: /dev/sdx
	DevicePath string `json:"device,omitempty"`
	// Cluster wide unique name of the volume. This can be name of the
	// persistent volume.
	VolumeName string `json:"volumeName,omitempty"`
	// Is the volume is attached on the node
	Attached bool `json:"attached,omitempty"`
	// Returns capabilities of the driver.
	// By default we assume all the capabilities are supported.
	// If the plugin does not support a capability, it can return false for that capability.
	Capabilities map[string]bool `json:"capabilities,omitempty"`
}

type VolumeMountRequest struct {
	MountDir   string            `json:"mountDir"`
	DevicePath string            `json:"devicePath"`
	Opts       map[string]string `json:"opts"`
}

type VolumeUnmountRequest struct {
	MountDir string `json:"mountDir"`
}

type VolumeDetachRequest struct {
	PvOrVolumeName string `json:"pvOrVolumeName"`
	HostName       string `json:"hostName"`
}

type VolumeAttachRequest struct {
	Opts     map[string]string `json:"opts"`
	HostName string            `json:"hostName"`
}
