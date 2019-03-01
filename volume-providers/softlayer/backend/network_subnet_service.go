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

// NetworkSubnetService is a wrapping interface for the softlayer-go API's NetworkSubnetService
//go:generate counterfeiter -o fakes/network_subnet_service.go --fake-name NetworkSubnetService . NetworkSubnetService
type NetworkSubnetService interface {
	Filter(string) NetworkSubnetService
	Mask(string) NetworkSubnetService
	ID(int) NetworkSubnetService

	GetObject() (resp datatypes.Network_Subnet, err error)
}

// NetworkSubnetServiceSL is a softlayer implementation of the networkSubnetService interface.
// All functions pass directly to the equivalent SL function
type NetworkSubnetServiceSL struct {
	networkSubnetService services.Network_Subnet
}

// ID pass-through for networkSubnetService.Id
func (ns *NetworkSubnetServiceSL) ID(id int) NetworkSubnetService {
	return &NetworkSubnetServiceSL{networkSubnetService: ns.networkSubnetService.Id(id)}
}

// Mask pass-through for NetworkSubnetService.Mask
func (ns *NetworkSubnetServiceSL) Mask(mask string) NetworkSubnetService {
	return &NetworkSubnetServiceSL{networkSubnetService: ns.networkSubnetService.Mask(mask)}
}

// Filter pass-through for NetworkSubnetService.Filter
func (ns *NetworkSubnetServiceSL) Filter(filter string) NetworkSubnetService {
	return &NetworkSubnetServiceSL{networkSubnetService: ns.networkSubnetService.Filter(filter)}
}

// GetObject pass-through for NetworkSubnetService.GetObject
func (ns *NetworkSubnetServiceSL) GetObject() (datatypes.Network_Subnet, error) {
	var subnet datatypes.Network_Subnet
	var ntError error
	ntError = retry(func() error {
		subnet, ntError = ns.networkSubnetService.GetObject()
		return ntError
	})
	return subnet, ntError
}
