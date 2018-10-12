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

// AccountService is a wrapping interface for the softlayer-go API's AccountService
//go:generate counterfeiter -o fakes/account_service.go --fake-name AccountService . AccountService
type AccountService interface {
	Filter(string) AccountService
	Mask(string) AccountService
	ID(int) AccountService
	Limit(int) AccountService
	Offset(int) AccountService

	GetBlockDeviceTemplateGroups() ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	GetCurrentUser() (datatypes.User_Customer, error)
	GetHardware() ([]datatypes.Hardware, error)
	GetObject() (datatypes.Account, error)
	GetOrders() ([]datatypes.Billing_Order, error)
	GetSubnets() ([]datatypes.Network_Subnet, error)
	GetNetworkStorage() (resp []datatypes.Network_Storage, err error)
}

// AccountServiceSL is a softlayer implementation of the AccountService interface.
// All functions pass directly to the equivalent SL function
type AccountServiceSL struct {
	accountService services.Account
}

// ID pass-through for AccountService.Id
func (as *AccountServiceSL) ID(id int) AccountService {
	return &AccountServiceSL{accountService: as.accountService.Id(id)}
}

// Mask pass-through for AccountService.Mask
func (as *AccountServiceSL) Mask(mask string) AccountService {
	return &AccountServiceSL{accountService: as.accountService.Mask(mask)}
}

// Filter pass-through for AccountService.Filter
func (as *AccountServiceSL) Filter(filter string) AccountService {
	return &AccountServiceSL{accountService: as.accountService.Filter(filter)}
}

// Limit pass-through for AccountService.Limit
func (as *AccountServiceSL) Limit(limit int) AccountService {
	return &AccountServiceSL{accountService: as.accountService.Limit(limit)}
}

// Offset pass-through for AccountService.Offset
func (as *AccountServiceSL) Offset(offset int) AccountService {
	return &AccountServiceSL{accountService: as.accountService.Offset(offset)}
}

// GetBlockDeviceTemplateGroups pass-through for AccountService.GetBlockDeviceTemplateGroups
func (as *AccountServiceSL) GetBlockDeviceTemplateGroups() ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	return as.accountService.GetBlockDeviceTemplateGroups()
}

// GetCurrentUser pass-through for AccountService.GetCurrentUser
func (as *AccountServiceSL) GetCurrentUser() (datatypes.User_Customer, error) {
	var usrCustomer datatypes.User_Customer
	var usrError error
	usrError = retry(func() error {
		usrCustomer, usrError = as.accountService.GetCurrentUser()
		return usrError
	})
	return usrCustomer, usrError
	//return as.accountService.GetCurrentUser()
}

// GetHardware pass-through for AccountService.GetHardware
func (as *AccountServiceSL) GetHardware() ([]datatypes.Hardware, error) {
	var hrdWare []datatypes.Hardware
	var hrdError error
	hrdError = retry(func() error {
		hrdWare, hrdError = as.accountService.GetHardware()
		return hrdError
	})
	return hrdWare, hrdError
}

// GetObject pass-through for AccountService.GetObject
func (as *AccountServiceSL) GetObject() (datatypes.Account, error) {
	var actInfo datatypes.Account
	var actError error
	actError = retry(func() error {
		actInfo, actError = as.accountService.GetObject()
		return actError
	})
	return actInfo, actError
}

//GetOrders pass-through for AccountService.GetOrders
func (as *AccountServiceSL) GetOrders() ([]datatypes.Billing_Order, error) {
	var blOrder []datatypes.Billing_Order
	var blError error
	blError = retry(func() error {
		blOrder, blError = as.accountService.GetOrders()
		return blError
	})
	return blOrder, blError
}

// GetSubnets pass-through for AccountService.GetSubnets
func (as *AccountServiceSL) GetSubnets() ([]datatypes.Network_Subnet, error) {
	var ntwSubnet []datatypes.Network_Subnet
	var subError error
	subError = retry(func() error {
		ntwSubnet, subError = as.accountService.GetSubnets()
		return subError
	})
	return ntwSubnet, subError
}

func (as *AccountServiceSL) GetNetworkStorage() (resp []datatypes.Network_Storage, err error) {
	var ntwStorage []datatypes.Network_Storage
	var stgError error
	stgError = retry(func() error {
		ntwStorage, stgError = as.accountService.GetNetworkStorage()
		return stgError
	})
	return ntwStorage, stgError
}
