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
	//"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
)

// AccountService is a wrapping interface for the softlayer-go API's AccountService
//go:generate counterfeiter -o fakes/account_service.go --fake-name AccountService . AccountService
type ResourceMetadataService interface {
	Filter(string) ResourceMetadataService
	Mask(string) ResourceMetadataService
	ID(int) ResourceMetadataService

	//func (r Resource_Metadata) GetDatacenterId() (resp int, err error) {
	GetDatacenterId() (resp int, err error)
}

// ResourceMetadataServiceSL is a softlayer implementation of the NetworkStorageService interface.
// All functiors pass directly to the equivalent SL function
type ResourceMetadataServiceSL struct {
	resourceMetadataService services.Resource_Metadata
}

// ID pass-through for NetworkStorageService.Id
func (rs *ResourceMetadataServiceSL) ID(id int) ResourceMetadataService {
	return &ResourceMetadataServiceSL{resourceMetadataService: rs.resourceMetadataService.Id(id)}
}

// Mask pass-through for ResourceMetadataService.Mask
func (rs *ResourceMetadataServiceSL) Mask(mask string) ResourceMetadataService {
	return &ResourceMetadataServiceSL{resourceMetadataService: rs.resourceMetadataService.Mask(mask)}
}

// Filter pass-through for ResourceMetadataService.Filter
func (rs *ResourceMetadataServiceSL) Filter(filter string) ResourceMetadataService {
	return &ResourceMetadataServiceSL{resourceMetadataService: rs.resourceMetadataService.Filter(filter)}
}

func (rs *ResourceMetadataServiceSL) GetDatacenterId() (int, error) {
	var dcID int
	var dcError error
	dcError = retry(func() error {
		dcID, dcError = rs.resourceMetadataService.GetDatacenterId()
		return dcError
	})
	return dcID, dcError
}
