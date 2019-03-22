/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package models

// Snapshot ...
type Snapshot struct {
	Href          string         `json:"href,omitempty"`
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	ResourceGroup *ResourceGroup `json:"resource_group,omitempty"`
	CRN           string         `json:"crn,omitempty"`
	Status        StatusType     `json:"status,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
}

// SnapshotList ...
type SnapshotList struct {
	Snapshots []*Snapshot `json:"snapshot,omitempty"`
}
