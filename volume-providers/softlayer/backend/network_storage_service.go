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
type NetworkStorageService interface {
	Filter(string) NetworkStorageService
	Mask(string) NetworkStorageService
	ID(int) NetworkStorageService

	GetObject() (resp datatypes.Network_Storage, err error)
	CreateSnapshot(notes *string) (resp datatypes.Network_Storage, err error)
	GetSnapshots() (resp []datatypes.Network_Storage, err error)
	DeleteObject() (resp bool, err error)
}

// NetworkStorageServiceSL is a softlayer implementation of the NetworkStorageService interface.
// All functions pass directly to the equivalent SL function
type NetworkStorageServiceSL struct {
	networkStorageService services.Network_Storage
}

// ID pass-through for NetworkStorageService.Id
func (ns *NetworkStorageServiceSL) ID(id int) NetworkStorageService {
	return &NetworkStorageServiceSL{networkStorageService: ns.networkStorageService.Id(id)}
}

// Mask pass-through for NetworkStorageIscsiService.Mask
func (ns *NetworkStorageServiceSL) Mask(mask string) NetworkStorageService {
	return &NetworkStorageServiceSL{networkStorageService: ns.networkStorageService.Mask(mask)}
}

// Filter pass-through for NetworkStorageIscsiService.Filter
func (ns *NetworkStorageServiceSL) Filter(filter string) NetworkStorageService {
	return &NetworkStorageServiceSL{networkStorageService: ns.networkStorageService.Filter(filter)}
}

// GetObject pass-through for NetworkStorageIscsiService.GetObject
func (ns *NetworkStorageServiceSL) GetObject() (datatypes.Network_Storage, error) {
	var ntStorage datatypes.Network_Storage
	var ntError error
	ntError = retry(func() error {
		ntStorage, ntError = ns.networkStorageService.GetObject()
		return ntError
	})
	return ntStorage, ntError
}

func (ns *NetworkStorageServiceSL) CreateSnapshot(notes *string) (resp datatypes.Network_Storage, err error) {
	var ntStorage datatypes.Network_Storage
	var ntError error
	ntError = retry(func() error {
		ntStorage, ntError = ns.networkStorageService.CreateSnapshot(notes)
		return ntError
	})
	return ntStorage, ntError
}

func (ns *NetworkStorageServiceSL) GetSnapshots() (resp []datatypes.Network_Storage, err error) {
	var ntStorage []datatypes.Network_Storage
	var ntError error
	ntError = retry(func() error {
		ntStorage, ntError = ns.networkStorageService.GetSnapshots()
		return ntError
	})
	return ntStorage, ntError
}

func (ns *NetworkStorageServiceSL) DeleteObject() (resp bool, err error) {
	var bStatus bool
	var dtError error
	dtError = retry(func() error {
		bStatus, dtError = ns.networkStorageService.DeleteObject()
		return dtError
	})
	return bStatus, dtError
}
