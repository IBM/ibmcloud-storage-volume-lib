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
	"github.com/softlayer/softlayer-go/sl"
)

// ProductOrderService is a wrapping interface for the softlayer-go API's ProductOrderService
//go:generate counterfeiter -o fakes/product_order_service.go --fake-name ProductOrderService . ProductOrderService
type ProductOrderService interface {
	Filter(string) ProductOrderService
	Mask(string) ProductOrderService
	ID(int) ProductOrderService
	//func (r Product_Order) PlaceOrder(orderData interface{}, saveAsQuote *bool) (resp datatypes.Container_Product_Order_Receipt, err error) {
	PlaceOrder(orderData interface{}, saveAsQuote *bool) (resp datatypes.Container_Product_Order_Receipt, err error)
}

// ProductOrderServiceSL is a softlayer implementation of the ProductOrderService interface.
// All functions pass directly to the equivalent SL function
type ProductOrderServiceSL struct {
	productOrderService services.Product_Order
}

// ID pass-through for NetworkStorageService.Id
func (ps *ProductOrderServiceSL) ID(id int) ProductOrderService {
	return &ProductOrderServiceSL{productOrderService: ps.productOrderService.Id(id)}
}

// Mask pass-through for NetworkStorageIscsiService.Mask
func (ps *ProductOrderServiceSL) Mask(mask string) ProductOrderService {
	return &ProductOrderServiceSL{productOrderService: ps.productOrderService.Mask(mask)}
}

// Filter pass-through for NetworkStorageIscsiService.Filter
func (ps *ProductOrderServiceSL) Filter(filter string) ProductOrderService {
	return &ProductOrderServiceSL{productOrderService: ps.productOrderService.Filter(filter)}
}

func (ps *ProductOrderServiceSL) PlaceOrder(orderData interface{}, saveAsQuote *bool) (resp datatypes.Container_Product_Order_Receipt, err error) {
	var prdOrderReceipt datatypes.Container_Product_Order_Receipt
	var prdError error
	prdError = retry(func() error {
		prdOrderReceipt, prdError = ps.productOrderService.PlaceOrder(orderData, sl.Bool(false))
		return prdError
	})
	return prdOrderReceipt, prdError
}
