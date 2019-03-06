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

// IPVersionType ...
type IPVersionType string

// String ...
func (i IPVersionType) String() string { return string(i) }

// IPVersion values
const (
	IPVersionIPV4 IPVersionType = "ipv4"
	IPVersionIPV6 IPVersionType = "ipv6"
	IPVersionBoth IPVersionType = "both"
)

// SubnetStatusType ...
type SubnetStatusType string

// String ...
func (s SubnetStatusType) String() string { return string(s) }

// SubnetStatus values
const (
	SubnetStatusAvailable SubnetStatusType = "available"
	SubnetStatusFailed    SubnetStatusType = "failed"
	SubnetStatusPending   SubnetStatusType = "pending"
)

// Subnet ...
type Subnet struct {
	Href string `json:"href,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	CRN string `json:"crn,omitempty"`

	AvailableIPV4AddressCount int64          `json:"available_ipv4_address_count,omitempty"`
	CreatedAt                 *time.Time     `json:"created_at,omitempty"`
	Generation                GenerationType `json:"generation,omitempty"`
	IPVersion                 IPVersionType  `json:"ip_version,omitempty"`
	IPV4CIDRBlock             string         `json:"ipv4_cidr_block,omitempty"`
	IPV6CIDRBlock             string         `json:"ipv6_cidr_block,omitempty"`
	//NetworkACL NetworkACL `json:"network_acl,omitempty"`
	PublicGateway         *PublicGateway   `json:"public_gateway,omitempty"`
	ResourceGroup         *ResourceGroup   `json:"resource_group,omitempty"`
	Status                SubnetStatusType `json:"status,omitempty"`
	Tags                  []string         `json:"tags,omitempty"`
	TotalIPV4AddressCount int64            `json:"total_ipv4_address_count,omitempty"`
	VPC                   *VPC             `json:"vpc,omitempty"`
	Zone                  *Zone            `json:"zone,omitempty"`
}

// SubnetList ...
type SubnetList struct {
	Subnets []*Subnet `json:"subnets,omitempty"`
	Limit   int       `json:"limit,omitempty"`
}
