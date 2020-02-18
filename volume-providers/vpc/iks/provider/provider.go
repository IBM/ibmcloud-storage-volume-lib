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
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/auth"
	userError "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/messages"
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

	// Update the iks api route to private route if cluster is private
	if conf.Bluemix.PrivateAPIRoute != "" {
		conf.Bluemix.APIEndpointURL = conf.Bluemix.PrivateAPIRoute
	}

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
	ccf, _ := iksp.vpcBlockProvider.ContextCredentialsFactory(nil)
	ctxLogger.Info("Its IKS dual session. Getttng IAM token for  VPC block session")
	vpcContextCredentials, err := ccf.ForIAMAccessToken(iksp.globalConfig.VPC.APIKey, ctxLogger)
	if err != nil {
		ctxLogger.Error("Error occured while generating IAM token for VPC", zap.Error(err))
		userErr := userError.GetUserError(string(userError.AuthenticationFailed), err)
		return nil, userErr
	}
	session, err := iksp.vpcBlockProvider.OpenSession(ctx, vpcContextCredentials, ctxLogger)
	if err != nil {
		ctxLogger.Error("Error occured while opening VPCSession", zap.Error(err))
		return nil, err
	}
	vpcSession, _ := session.(*vpcprovider.VPCSession)
	ctxLogger.Info("Opening IKS block session")

	//Create ContextCredentialsFactory
	ccf, err = iksp.ContextCredentialsFactory(nil)
	if err != nil {
		ctxLogger.Error("Error while creating the ContextCredentialsFactory", zap.Error(err))
		return nil, err
	}
	iksp.iksBlockProvider.ContextCF = ccf

	ctxLogger.Info("Its ISK dual session. Getttng IAM token for  IKS block session")
	iksContextCredentials, err := ccf.ForIAMAccessToken(iksp.globalConfig.Bluemix.IamAPIKey, ctxLogger)
	if err != nil {
		ctxLogger.Warn("Error occured while generating IAM token for IKS. But continue with VPC session alone. \n Volume Mount operation will fail but volume provisioning will work", zap.Error(err))
		session = &vpcprovider.VPCSession{} // Empty session to avoid Nil references.
	} else {
		session, err = iksp.iksBlockProvider.OpenSession(ctx, iksContextCredentials, ctxLogger)
		if err != nil {
			ctxLogger.Error("Error occured while opening IKSSession", zap.Error(err))
		}
	}

	iksSession, ok := session.(*vpcprovider.VPCSession)
	if ok && iksSession.Apiclient != nil {
		iksSession.APIClientVolAttachMgr = iksSession.Apiclient.IKSVolumeAttachService()
	}
	// Setup Dual Session that handles for VPC and IKS connections
	vpcIksSession := IksVpcSession{
		VPCSession: *vpcSession,
		IksSession: iksSession,
	}
	ctxLogger.Debug("IksVpcSession", zap.Reflect("IksVpcSession", vpcIksSession))
	return &vpcIksSession, nil
}

// ContextCredentialsFactory ...
func (iksp *IksVpcBlockProvider) ContextCredentialsFactory(zone *string) (local.ContextCredentialsFactory, error) {
	//  Datacenter hint not required by IKS provider implementation
	// VPC provider use differnt APIkey and Auth Endpoints
	authConfig := &config.BluemixConfig{
		IamURL:          iksp.globalConfig.Bluemix.IamURL,
		IamAPIKey:       iksp.globalConfig.Bluemix.IamAPIKey,
		IamClientID:     iksp.globalConfig.Bluemix.IamClientID,
		IamClientSecret: iksp.globalConfig.Bluemix.IamClientSecret,
		PrivateAPIRoute: iksp.globalConfig.Bluemix.PrivateAPIRoute, // Only for private cluster
		CSRFToken:       iksp.globalConfig.Bluemix.CSRFToken,       // required for private cluster
	}
	return auth.NewContextCredentialsFactory(authConfig, nil, iksp.globalConfig.VPC)
}
