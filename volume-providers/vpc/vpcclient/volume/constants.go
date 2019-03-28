/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package volume

const (
	volumesPath   = "volumes"
	volumeIDParam = "volume-id"
	volumeIDPath  = volumesPath + "/{" + volumeIDParam + "}"

	snapshotsPath   = volumesPath + "/{" + volumeIDParam + "}" + "snapshots"
	snapshotIDParam = "snapshot-id"
	snapshotIDPath  = snapshotsPath + "/{" + snapshotIDParam + "}"

	volumeTagsPath    = volumesPath + "/{" + volumeIDParam + "}" + "tags"
	volumeTagParam    = "tag-name"
	volumeTagNamePath = volumeTagsPath + "/{" + volumeTagParam + "}"

	snapshotTagsPath    = snapshotIDPath + "/" + "tags"
	snapshotTagParam    = "tag-name"
	snapshotTagNamePath = snapshotTagsPath + "/{" + snapshotTagParam + "}"
)
