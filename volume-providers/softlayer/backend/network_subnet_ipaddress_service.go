/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package backend

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
)

// NetworkSubnetIpAddressService is a wrapping interface for the softlayer-go API's NetworkSubnetIpAddressService
//go:generate counterfeiter -o fakes/network_subnet_service.go --fake-name NetworkSubnetIpAddressService . NetworkSubnetIpAddressService
type NetworkSubnetIpAddressService interface {
	Filter(string) NetworkSubnetIpAddressService
	Mask(string) NetworkSubnetIpAddressService
	ID(int) NetworkSubnetIpAddressService

	GetObject() (resp datatypes.Network_Subnet_IpAddress, err error)
	GetByIPAddress(ipAddress *string) (datatypes.Network_Subnet_IpAddress, error)
}

// NetworkSubnetIpAddressServiceSL is a softlayer implementation of the networkSubnetIpAddressService interface.
// All functions pass directly to the equivalent SL function
type NetworkSubnetIpAddressServiceSL struct {
	networkSubnetIpAddressService services.Network_Subnet_IpAddress
}

// ID pass-through for networkSubnetIpAddressService.Id
func (ns *NetworkSubnetIpAddressServiceSL) ID(id int) NetworkSubnetIpAddressService {
	return &NetworkSubnetIpAddressServiceSL{networkSubnetIpAddressService: ns.networkSubnetIpAddressService.Id(id)}
}

// Mask pass-through for NetworkSubnetIpAddressService.Mask
func (ns *NetworkSubnetIpAddressServiceSL) Mask(mask string) NetworkSubnetIpAddressService {
	return &NetworkSubnetIpAddressServiceSL{networkSubnetIpAddressService: ns.networkSubnetIpAddressService.Mask(mask)}
}

// Filter pass-through for NetworkSubnetIpAddressService.Filter
func (ns *NetworkSubnetIpAddressServiceSL) Filter(filter string) NetworkSubnetIpAddressService {
	return &NetworkSubnetIpAddressServiceSL{networkSubnetIpAddressService: ns.networkSubnetIpAddressService.Filter(filter)}
}

// GetObject pass-through for NetworkSubnetIpAddressService.GetObject
func (ns *NetworkSubnetIpAddressServiceSL) GetObject() (datatypes.Network_Subnet_IpAddress, error) {
	var subnetIP datatypes.Network_Subnet_IpAddress
	var ntError error
	ntError = retry(func() error {
		subnetIP, ntError = ns.networkSubnetIpAddressService.GetObject()
		return ntError
	})
	return subnetIP, ntError
}

// GetByIPAddress returns Network_Subnet_IpAddress from IP address string
func (ns *NetworkSubnetIpAddressServiceSL) GetByIPAddress(ipAddress *string) (datatypes.Network_Subnet_IpAddress, error) {

	var subnetIP datatypes.Network_Subnet_IpAddress
	var nsErr error
	nsErr = retry(func() error {
		//services.GetNetworkSubnetService(sess)
		subnetIP, nsErr = ns.networkSubnetIpAddressService.GetByIpAddress(ipAddress)
		return nsErr
	})
	return subnetIP, nsErr
}
