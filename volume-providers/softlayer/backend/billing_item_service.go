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
	"github.com/softlayer/softlayer-go/services"
)

// BillingItemService is a wrapping interface for the softlayer-go API's BillingItemService
//go:generate counterfeiter -o fakes/billing_item_service.go --fake-name BillingItemService . BillingItemService
type BillingItemService interface {
	Filter(filter string) BillingItemService
	Mask(string) BillingItemService
	ID(int) BillingItemService
	CancelItem(cancelImmediately *bool, cancelAssociatedBillingItems *bool, reason *string, customerNote *string) (resp bool, err error)
	CancelService() (resp bool, err error)
}

// BillingItemServiceSL is a softlayer implementation of the BillingItemService interface.
// All functions pass directly to the equivalent SL function
type BillingItemServiceSL struct {
	billingItemService services.Billing_Item
}

// Filter pass-through for BillingItemService.Filter
func (bis *BillingItemServiceSL) Filter(filter string) BillingItemService {
	return &BillingItemServiceSL{billingItemService: bis.billingItemService.Filter(filter)}
}

// Mask pass-through for BillingItemService.Mask
func (bis *BillingItemServiceSL) Mask(mask string) BillingItemService {
	return &BillingItemServiceSL{billingItemService: bis.billingItemService.Mask(mask)}
}

// ID pass-through for BillingItemService.ID
func (bis *BillingItemServiceSL) ID(id int) BillingItemService {
	return &BillingItemServiceSL{billingItemService: bis.billingItemService.Id(id)}
}

// CancelItem cancels the billing item
func (bis *BillingItemServiceSL) CancelItem(cancelImmediately *bool, cancelAssociatedBillingItems *bool, reason *string, customerNote *string) (bool, error) {
	var bStatus bool
	var ciError error
	ciError = retry(func() error {
		bStatus, ciError = bis.billingItemService.CancelItem(cancelImmediately, cancelAssociatedBillingItems, reason, customerNote)
		return ciError
	})
	return bStatus, ciError
}

func (bis *BillingItemServiceSL) CancelService() (resp bool, err error) {
	var bStatus bool
	var ciError error
	ciError = retry(func() error {
		bStatus, ciError = bis.billingItemService.CancelService()
		return ciError
	})
	return bStatus, ciError
}
