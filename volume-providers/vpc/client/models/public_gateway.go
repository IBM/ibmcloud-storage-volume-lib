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

import "time"

// PublicGatewayStatusType ...
type PublicGatewayStatusType string

// String ...
func (s PublicGatewayStatusType) String() string { return string(s) }

// PublicGatewayStatus values
const (
	PublicGatewayStatusAvailable PublicGatewayStatusType = "available"
	PublicGatewayStatusFailed    PublicGatewayStatusType = "failed"
	PublicGatewayStatusPending   PublicGatewayStatusType = "pending"
)

// PublicGateway ...
type PublicGateway struct {
	Href string `json:"href,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	CRN  string `json:"crn,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	//FloatingIP *FloatingIP            `json:"floating_ip,omitempty"`
	ResourceGroup *ResourceGroup          `json:"resource_group,omitempty"`
	Status        PublicGatewayStatusType `json:"status,omitempty"`
	Tags          []string                `json:"tags,omitempty"`
	VPC           *VPC                    `json:"vpc,omitempty"`
	Zone          *Zone                   `json:"zone,omitempty"`
}
