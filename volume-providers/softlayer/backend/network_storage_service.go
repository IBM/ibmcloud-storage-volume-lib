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
	//"github.com/softlayer/softlayer-go/sl"
)

// NetworkStorageService is a wrapping interface for the softlayer-go API's NetworkStorageService
//go:generate counterfeiter -o fakes/network_storage_service.go --fake-name NetworkStorageService . NetworkStorageService
type NetworkStorageService interface {
	Filter(string) NetworkStorageService
	Mask(string) NetworkStorageService
	ID(int) NetworkStorageService

	GetObject() (resp datatypes.Network_Storage, err error)
	CreateSnapshot(notes *string) (resp datatypes.Network_Storage, err error)
	GetSnapshots() (resp []datatypes.Network_Storage, err error)
	DeleteObject() (resp bool, err error)
	EditObject(templateObject *datatypes.Network_Storage) (resp bool, err error)
	AllowAccessFromSubnetList(subnetObjectTemplates []datatypes.Network_Subnet) (bool, error)
	AllowAccessFromIpAddressList(ipAddressObjectTemplates []datatypes.Network_Subnet_IpAddress) (bool, error)
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

func (ns *NetworkStorageServiceSL) EditObject(templateObject *datatypes.Network_Storage) (resp bool, err error) {
	var editStatus bool
	var editError error
	editError = retry(func() error {
		editStatus, editError = ns.networkStorageService.EditObject(templateObject)
		return editError
	})
	return editStatus, editError
}

func (ns *NetworkStorageServiceSL) AllowAccessFromSubnetList(subnetObjectTemplates []datatypes.Network_Subnet) (bool, error) {
	var bStatus bool
	var dtError error
	dtError = retry(func() error {
		bStatus, dtError = ns.networkStorageService.AllowAccessFromSubnetList(subnetObjectTemplates)
		return dtError
	})
	return bStatus, dtError
}

func (ns *NetworkStorageServiceSL) AllowAccessFromIpAddressList(ipAddressObjectTemplates []datatypes.Network_Subnet_IpAddress) (bool, error) {
	//Allow access to current host
	/*var allowedHost []datatypes.Network_Storage_Allowed_Host
	hostList := []datatypes.Container_Network_Storage_Host{}
	for _, ipAddress := range ipAddressObjectTemplates {
		host := datatypes.Container_Network_Storage_Host{Id: ipAddress.Id, ObjectType: sl.String("SoftLayer_Network_Subnet_IpAddress")}
		hostList = append(hostList, host)
	}*/

	var dtError error
	var dtStatus bool
	dtError = retry(func() error {
		//services.GetNetworkSubnetService(sess)
		dtStatus, dtError = ns.networkStorageService.AllowAccessFromIpAddressList(ipAddressObjectTemplates)
		return dtError
	})
	return dtStatus, dtError
}
