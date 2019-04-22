/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

// Context represents the volume provider management API for individual account, user ID, etc.
//go:generate counterfeiter -o fakes/context.go --fake-name Context . Context
type Context interface {
	VolumeManager
	VolumeMountManager
	SnapshotManager
}

// Session is an Context that is notified when it is no longer required
//go:generate counterfeiter -o fakes/session.go --fake-name Session . Session
type Session interface {
	Context

	// GetProviderDisplayName returns the name of the provider that is being used
	GetProviderDisplayName() VolumeProvider

	// Close is called when the Session is nolonger required
	Close()
}
