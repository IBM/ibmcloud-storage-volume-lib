/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package fakes

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend"
)

// Session implements the backend.Session interface using a real softlayer-go Session
type Session struct {
	billingItemService            *BillingItemService
	billingOrderService           *BillingOrderService
	accountService                *AccountService
	networkStorageIscsiService    *NetworkStorageIscsiService
	productOrderService           *ProductOrderService
	productPackageService         *ProductPackageService
	resourceMetadataService       *ResourceMetadataService
	networkStorageService         *NetworkStorageService
	networkSubnetService          *NetworkSubnetService
	networkSubnetIpAddressService *NetworkSubnetIpAddressService
	locationService               *LocationService
}

func NewSession() *Session {
	return &Session{
		billingItemService:         &BillingItemService{},
		billingOrderService:        &BillingOrderService{},
		accountService:             &AccountService{},
		networkStorageIscsiService: &NetworkStorageIscsiService{},
		productOrderService:        &ProductOrderService{},
		productPackageService:      &ProductPackageService{},
		resourceMetadataService:    &ResourceMetadataService{},
		networkStorageService:      &NetworkStorageService{},
		locationService:            &LocationService{},
	}
}

// GetBillingItemService returns the BillingItemService from the session
func (s *Session) GetBillingItemService() backend.BillingItemService {
	return s.billingItemService
}

// GetBillingItemService returns the BillingItemService from the session
func (s *Session) GetBillingOrderService() backend.BillingOrderService {
	return s.billingOrderService
}

// GetAccountService returns the AccountService from the session
func (s *Session) GetAccountService() backend.AccountService {
	return s.accountService
}

func (s *Session) GetNetworkStorageIscsiService() backend.NetworkStorageIscsiService {
	return s.networkStorageIscsiService
}

func (s *Session) GetProductOrderService() backend.ProductOrderService {
	return s.productOrderService
}

func (s *Session) GetProductPackageService() backend.ProductPackageService {
	return s.productPackageService
}

func (s *Session) GetResourceMetadataService() backend.ResourceMetadataService {
	return s.resourceMetadataService
}
func (s *Session) GetNetworkStorageService() backend.NetworkStorageService {
	return s.networkStorageService
}

//GetNetworkSubnetService ...
func (s *Session) GetNetworkSubnetService() backend.NetworkSubnetService {
	return s.networkSubnetService
}

//GetNetworkSubnetIpAddressService ...
func (s *Session) GetNetworkSubnetIpAddressService() backend.NetworkSubnetIpAddressService {
	return s.networkSubnetIpAddressService
}

func (s *Session) GetLocationService() backend.LocationService {
	return s.locationService
}

// GetBillingItemService returns the BillingItemService from the session
func (s *Session) GetBillingItemServiceFake() *BillingItemService {
	return s.billingItemService
}

// GetBillingItemService returns the BillingItemService from the session
func (s *Session) GetBillingOrderServiceFake() *BillingOrderService {
	return s.billingOrderService
}

// GetAccountService returns the AccountService from the session
func (s *Session) GetAccountServiceFake() *AccountService {
	return s.accountService
}

func (s *Session) GetNetworkStorageIscsiServiceFake() *NetworkStorageIscsiService {
	return s.networkStorageIscsiService
}

func (s *Session) GetProductOrderServiceFake() *ProductOrderService {
	return s.productOrderService
}

func (s *Session) GetProductPackageServiceFake() *ProductPackageService {
	return s.productPackageService
}

func (s *Session) GetResourceMetadataServiceFake() *ResourceMetadataService {
	return s.resourceMetadataService
}
func (s *Session) GetNetworkStorageServiceFake() *NetworkStorageService {
	return s.networkStorageService
}

//GetNetworkSubnetServiceFake ...
func (s *Session) GetNetworkSubnetServiceFake() *NetworkSubnetService {
	return s.networkSubnetService
}

//GetNetworkSubnetIpAddressServiceFake ...
func (s *Session) GetNetworkSubnetIpAddressServiceFake() *NetworkSubnetIpAddressService {
	return s.networkSubnetIpAddressService
}

func (s *Session) GetLocationServiceFake() *LocationService {
	return s.locationService
}
