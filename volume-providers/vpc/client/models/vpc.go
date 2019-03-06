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

// VPC ...
type VPC struct {
	Href          string         `json:"href,omitempty"`
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	CRN           string         `json:"crn,omitempty"`
	ResourceGroup *ResourceGroup `json:"resource_group,omitempty"`
}

// VPCList ...
type VPCList struct {
	VPCs  []*VPC `json:"vpcs,omitempty"`
	Limit int    `json:"limit,omitempty"`
}
