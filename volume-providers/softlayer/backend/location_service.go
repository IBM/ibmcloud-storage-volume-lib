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

// LocationService is a wrapping interface for the softlayer-go API's LocationService
//go:generate counterfeiter -o fakes/location_service.go --fake-name LocationService . LocationService
type LocationService interface {
	Filter(string) LocationService
	Mask(string) LocationService
	ID(int) LocationService
	//func (r Location) GetDatacenters() (resp []datatypes.Location, err error) {
	GetDatacenters() (resp []datatypes.Location, err error)
}

// LocationServiceSL is a softlayer implementation of the LocationService interface.
// All functions pass directly to the equivalent SL function
type LocationServiceSL struct {
	locationService services.Location
}

// ID pass-through for NetworkStorageService.Id
func (ls *LocationServiceSL) ID(id int) LocationService {
	return &LocationServiceSL{locationService: ls.locationService.Id(id)}
}

// Mask pass-through for NetworkStorageIscsiService.Mask
func (ls *LocationServiceSL) Mask(mask string) LocationService {
	return &LocationServiceSL{locationService: ls.locationService.Mask(mask)}
}

// Filter pass-through for NetworkStorageIscsiService.Filter
func (ls *LocationServiceSL) Filter(filter string) LocationService {
	return &LocationServiceSL{locationService: ls.locationService.Filter(filter)}
}

func (ls *LocationServiceSL) GetDatacenters() (resp []datatypes.Location, err error) {
	var locaions []datatypes.Location
	var lError error
	lError = retry(func() error {
		locaions, lError = ls.locationService.GetDatacenters()
		return lError
	})
	return locaions, lError
}
