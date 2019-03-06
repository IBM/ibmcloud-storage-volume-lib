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
	"net/http"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"go.uber.org/zap"
)

// SessionSL implements the backend.Session interface using a real softlayer-go Session
type SessionSL struct {
	session *session.Session
}

// GetBillingItemService returns the BillingItemService from the session
func (theSession *SessionSL) GetBillingItemService() BillingItemService {
	return &BillingItemServiceSL{billingItemService: services.GetBillingItemService(theSession.session)}
}

// GetBillingOrderService returns the BillingOrderService from the session
func (theSession *SessionSL) GetBillingOrderService() BillingOrderService {
	return &BillingOrderServiceSL{billingOrderService: services.GetBillingOrderService(theSession.session)}
}

// GetAccountService returns the AccountService from the session
func (theSession *SessionSL) GetAccountService() AccountService {
	return &AccountServiceSL{accountService: services.GetAccountService(theSession.session)}
}

func (theSession *SessionSL) GetNetworkStorageIscsiService() NetworkStorageIscsiService {
	return &NetworkStorageIscsiServiceSL{networkStorageIscsiService: services.GetNetworkStorageIscsiService(theSession.session)}
}

func (theSession *SessionSL) GetProductOrderService() ProductOrderService {
	return &ProductOrderServiceSL{productOrderService: services.GetProductOrderService(theSession.session)}
}

func (theSession *SessionSL) GetProductPackageService() ProductPackageService {
	return &ProductPackageServiceSL{productPackageService: services.GetProductPackageService(theSession.session)}
}

func (theSession *SessionSL) GetNetworkStorageService() NetworkStorageService {
	return &NetworkStorageServiceSL{networkStorageService: services.GetNetworkStorageService(theSession.session)}
}

func (theSession *SessionSL) GetResourceMetadataService() ResourceMetadataService {
	return &ResourceMetadataServiceSL{resourceMetadataService: services.GetResourceMetadataService(theSession.session)}
}

func (theSession *SessionSL) GetLocationService() LocationService {
	return &LocationServiceSL{locationService: services.GetLocationService(theSession.session)}
}

// NewSoftLayerSession creates a Session backed using a real softlayer-go session
func NewSoftLayerSession(url string, conf provider.ContextCredentials, httpClient *http.Client, debug bool, logger *zap.Logger) Session {
	sess := session.New(conf.UserID, conf.Credential, url)
	sess.Debug = debug
	sess.HTTPClient = httpClient
	return &SessionSL{session: sess}
}
