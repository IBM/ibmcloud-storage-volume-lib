/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package registry

import (
	//"github.com/prometheus/client_golang/prometheus"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
)

// Providers is a registry interface for IaaS providers
//go:generate counterfeiter -o fakes/provider_registry.go --fake-name Providers . Providers
type Providers interface {
	Get(providerID string) (local.Provider, error)
	Register(providerID string, prov local.Provider)
}

var _ Providers = &ProviderRegistry{}

// ProviderRegistry is the core implementation of the Providers registry
type ProviderRegistry struct {
	providers map[string]local.Provider
}

// Get returns the identified Provider
func (pr *ProviderRegistry) Get(providerID string) (prov local.Provider, err error) {
	prov = pr.providers[providerID]
	if prov == nil {
		err = util.NewError("ErrorUnclassified", "Provider unknown: "+providerID)
	}
	return
}

// Register registers a given provider under the supplied key
func (pr *ProviderRegistry) Register(providerID string, p local.Provider) {
	if pr.providers == nil {
		pr.providers = map[string]local.Provider{}
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.InfoLevel)
	pr.providers[providerID] = p
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	), zap.AddCaller()).With(zap.String("name", "ibm-volume-lib/main")).With(zap.String("VolumeLib", "IKS-VOLUME-LIB"))

	defer logger.Sync()
	logger.Info("providerID", zap.Reflect("providers", pr.providers))
}
