/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"context"
	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	vpcprovider "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/provider"
	"go.uber.org/zap"
)

//IksVpcBlockProvider  handles both IKS and  RIAAS sessions
type IksVpcBlockProvider struct {
	vpcprovider.VPCBlockProvider
	vpcBlockProvider *vpcprovider.VPCBlockProvider // Holds VPC provider. Requires to avoid recursive calls
	iksBlockProvider *vpcprovider.VPCBlockProvider // Holds IKS provider
	globalConfig     *config.Config
}

var _ local.Provider = &IksVpcBlockProvider{}

//NewProvider handles both IKS and  RIAAS sessions
func NewProvider(conf *config.Config, logger *zap.Logger) (local.Provider, error) {
	//Setup vpc provider
	provider, _ := vpcprovider.NewProvider(conf, logger)
	vpcBlockProvider, _ := provider.(*vpcprovider.VPCBlockProvider)
	// Setup IKS provider
	provider, _ = vpcprovider.NewProvider(conf, logger)
	iksBlockProvider, _ := provider.(*vpcprovider.VPCBlockProvider)
	//Overrider Base URL
	iksBlockProvider.APIConfig.BaseURL = conf.Bluemix.APIEndpointURL
	// Setup IKS-VPC dual provider
	iksVpcBlockProvider := &IksVpcBlockProvider{
		VPCBlockProvider: *vpcBlockProvider,
		vpcBlockProvider: vpcBlockProvider,
		iksBlockProvider: iksBlockProvider,
		globalConfig:     conf,
	}

	//vpcBlockProvider.ApiConfig.BaseURL = conf.Bluemix.APIEndpointURL
	return iksVpcBlockProvider, nil

}

// OpenSession opens a session on the provider
func (iksp *IksVpcBlockProvider) OpenSession(ctx context.Context, contextCredentials provider.ContextCredentials, ctxLogger *zap.Logger) (provider.Session, error) {
	ctxLogger.Info("Entering IksVpcBlockProvider.OpenSession")

	defer func() {
		ctxLogger.Debug("Exiting IksVpcBlockProvider.OpenSession")
	}()
	ctxLogger.Info("Opening VPC block session")
	session, err := iksp.vpcBlockProvider.OpenSession(ctx, contextCredentials, ctxLogger)
	if err != nil {
		ctxLogger.Error("Error occured while opening VPCSession", zap.Error(err))
		return nil, err
	}
	vpcSession, _ := session.(*vpcprovider.VPCSession)
	ctxLogger.Info("Opening IKS block session")
	session, err = iksp.iksBlockProvider.OpenSession(ctx, contextCredentials, ctxLogger)
	if err != nil {
		ctxLogger.Error("Error occured while opening IKSSession", zap.Error(err))
		return nil, err
	}
	iksSession, _ := session.(*vpcprovider.VPCSession)
	// Setup Dual Session that handles for VPC and IKS connections
	vpcIksSession := IksVpcSession{
		VPCSession: *vpcSession,
		IksSession: iksSession,
	}
	ctxLogger.Debug("IksVpcSession", zap.Reflect("IksVpcSession", vpcIksSession))
	return &vpcIksSession, nil
}
