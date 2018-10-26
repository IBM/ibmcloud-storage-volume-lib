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

// BillingOrderService is a wrapping interface for the softlayer-go API's BillingOrderService
//go:generate counterfeiter -o fakes/billing_order_service.go --fake-name BillingOrderService . BillingOrderService
type BillingOrderService interface {
	Filter(filter string) BillingOrderService
	Mask(string) BillingOrderService
	ID(int) BillingOrderService
	GetObject() (datatypes.Billing_Order, error)
	GetAllObjects() ([]datatypes.Billing_Order, error)
}

// BillingOrderServiceSL is a softlayer implementation of the BillingOrderService interface.
// All functions pass directly to the equivalent SL function
type BillingOrderServiceSL struct {
	billingOrderService services.Billing_Order
}

// Filter pass-through for BillingOrderService.Filter
func (bos *BillingOrderServiceSL) Filter(filter string) BillingOrderService {
	return &BillingOrderServiceSL{billingOrderService: bos.billingOrderService.Filter(filter)}
}

// Mask pass-through for BillingOrderService.Mask
func (bos *BillingOrderServiceSL) Mask(mask string) BillingOrderService {
	return &BillingOrderServiceSL{billingOrderService: bos.billingOrderService.Mask(mask)}
}

// ID pass-through for BillingOrderService.ID
func (bos *BillingOrderServiceSL) ID(id int) BillingOrderService {
	return &BillingOrderServiceSL{billingOrderService: bos.billingOrderService.Id(id)}
}

// GetObject returns the billing order
func (bos *BillingOrderServiceSL) GetObject() (datatypes.Billing_Order, error) {
	var blOrder datatypes.Billing_Order
	var blError error
	blError = retry(func() error {
		blOrder, blError = bos.billingOrderService.GetObject()
		return blError
	})
	return blOrder, blError
}

// GetAllObjects returns all billing orders
func (bos *BillingOrderServiceSL) GetAllObjects() ([]datatypes.Billing_Order, error) {
	var blOrder []datatypes.Billing_Order
	var blError error
	blError = retry(func() error {
		blOrder, blError = bos.billingOrderService.GetAllObjects()
		return blError
	})
	return blOrder, blError
}
