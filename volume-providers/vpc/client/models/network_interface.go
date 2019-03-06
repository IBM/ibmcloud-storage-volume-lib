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

// PortSpeedType ...
type PortSpeedType int

// String ...
func (i PortSpeedType) String() string { return string(i) }

// Port speed values
const (
	PortSpeed10   PortSpeedType = 10
	PortSpeed100  PortSpeedType = 100
	PortSpeed1000 PortSpeedType = 1000
)

// NetworkInterface ...
type NetworkInterface struct {
	Href string `json:"href,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	PortSpeed          PortSpeedType `json:"port_speed,omitempty"`
	PrimaryIPV4Address string        `json:"primary_ipv4_address,omitempty"`
	Subnet             *Subnet       `json:"subnet,omitempty"`
}
