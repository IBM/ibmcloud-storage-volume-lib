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
type ProductPackageService interface {
	Filter(string) ProductPackageService
	Mask(string) ProductPackageService
	ID(int) ProductPackageService
	//func (r Product_Package) GetAllObjects() (resp []datatypes.Product_Package, err error) {
	GetAllObjects() (resp []datatypes.Product_Package, err error)
	//func (r Product_Package) GetItemPrices() (resp []datatypes.Product_Item_Price, err error) {
	GetItemPrices() (resp []datatypes.Product_Item_Price, err error)
}

// ProductOrderServiceSL is a softlayer implementation of the ProductOrderService interface.
// All functions pass directly to the equivalent SL function
type ProductPackageServiceSL struct {
	productPackageService services.Product_Package
}

// ID pass-through for NetworkStorageService.Id
func (ps *ProductPackageServiceSL) ID(id int) ProductPackageService {
	return &ProductPackageServiceSL{productPackageService: ps.productPackageService.Id(id)}
}

// Mask pass-through for NetworkStorageIscsiService.Mask
func (ps *ProductPackageServiceSL) Mask(mask string) ProductPackageService {
	return &ProductPackageServiceSL{productPackageService: ps.productPackageService.Mask(mask)}
}

// Filter pass-through for NetworkStorageIscsiService.Filter
func (ps *ProductPackageServiceSL) Filter(filter string) ProductPackageService {
	return &ProductPackageServiceSL{productPackageService: ps.productPackageService.Filter(filter)}
}

func (ps *ProductPackageServiceSL) GetAllObjects() (resp []datatypes.Product_Package, err error) {
	var prdPackage []datatypes.Product_Package
	var prdError error
	prdError = retry(func() error {
		prdPackage, prdError = ps.productPackageService.GetAllObjects()
		return prdError
	})
	return prdPackage, prdError
}

func (ps *ProductPackageServiceSL) GetItemPrices() (resp []datatypes.Product_Item_Price, err error) {
	var prdItemPrice []datatypes.Product_Item_Price
	var prdError error
	prdError = retry(func() error {
		prdItemPrice, prdError = ps.productPackageService.GetItemPrices()
		return prdError
	})
	return prdItemPrice, prdError
}
