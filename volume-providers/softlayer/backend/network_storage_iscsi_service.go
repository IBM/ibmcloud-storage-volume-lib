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

// NetworkStorageIscsiService is a wrapping interface for the softlayer-go API's NetworkStorageIscsiService
//go:generate counterfeiter -o fakes/network_storage_iscsi_service.go --fake-name NetworkStorageIscsiService . NetworkStorageIscsiService
type NetworkStorageIscsiService interface {
	Filter(string) NetworkStorageIscsiService
	Mask(string) NetworkStorageIscsiService
	ID(int) NetworkStorageIscsiService

	//func (r Network_Storage_Iscsi) CreateSnapshot(notes *string) (resp datatypes.Network_Storage, err error)
	CreateSnapshot(notes *string) (resp datatypes.Network_Storage, err error)
	//func (r Network_Storage_Iscsi) GetObject() (resp datatypes.Network_Storage_Iscsi, err error)
	GetObject() (resp datatypes.Network_Storage_Iscsi, err error)
	//func (r Network_Storage_Iscsi) GetSnapshots() (resp []datatypes.Network_Storage, err error) {
	GetSnapshots() (resp []datatypes.Network_Storage, err error)
	//func (r Network_Storage_Iscsi) GetSnapshotsForVolume() (resp []datatypes.Network_Storage, err error) {
	GetSnapshotsForVolume() (resp []datatypes.Network_Storage, err error)
}

// NetworkStorageServiceSL is a softlayer implementation of the NetworkStorageService interface.
// All functions pass directly to the equivalent SL function
type NetworkStorageIscsiServiceSL struct {
	networkStorageIscsiService services.Network_Storage_Iscsi
}

// ID pass-through for NetworkStorageService.Id
func (ns *NetworkStorageIscsiServiceSL) ID(id int) NetworkStorageIscsiService {
	return &NetworkStorageIscsiServiceSL{networkStorageIscsiService: ns.networkStorageIscsiService.Id(id)}
}

// Mask pass-through for NetworkStorageIscsiService.Mask
func (ns *NetworkStorageIscsiServiceSL) Mask(mask string) NetworkStorageIscsiService {
	return &NetworkStorageIscsiServiceSL{networkStorageIscsiService: ns.networkStorageIscsiService.Mask(mask)}
}

// Filter pass-through for NetworkStorageIscsiService.Filter
func (ns *NetworkStorageIscsiServiceSL) Filter(filter string) NetworkStorageIscsiService {
	return &NetworkStorageIscsiServiceSL{networkStorageIscsiService: ns.networkStorageIscsiService.Filter(filter)}
}

// GetObject pass-through for NetworkStorageIscsiService.GetObject
func (ns *NetworkStorageIscsiServiceSL) GetObject() (datatypes.Network_Storage_Iscsi, error) {
	var ntwSorage datatypes.Network_Storage_Iscsi
	var ntwError error
	ntwError = retry(func() error {
		ntwSorage, ntwError = ns.networkStorageIscsiService.GetObject()
		return ntwError
	})
	return ntwSorage, ntwError
}

//GetOrders pass-through for NetworkStorageIscsiService.CreateSnapshot
func (ns *NetworkStorageIscsiServiceSL) CreateSnapshot(notes *string) (datatypes.Network_Storage, error) {
	var ntwSorage datatypes.Network_Storage
	var ntwError error
	ntwError = retry(func() error {
		ntwSorage, ntwError = ns.networkStorageIscsiService.CreateSnapshot(notes)
		return ntwError
	})
	return ntwSorage, ntwError
}

func (ns *NetworkStorageIscsiServiceSL) GetSnapshots() (resp []datatypes.Network_Storage, err error) {
	var ntwSorage []datatypes.Network_Storage
	var ntwError error
	ntwError = retry(func() error {
		ntwSorage, ntwError = ns.networkStorageIscsiService.GetSnapshots()
		return ntwError
	})
	return ntwSorage, ntwError
}

func (ns *NetworkStorageIscsiServiceSL) GetSnapshotsForVolume() (resp []datatypes.Network_Storage, err error) {
	var ntwSorage []datatypes.Network_Storage
	var ntwError error
	ntwError = retry(func() error {
		ntwSorage, ntwError = ns.networkStorageIscsiService.GetSnapshotsForVolume()
		return ntwError
	})
	return ntwSorage, ntwError
}
