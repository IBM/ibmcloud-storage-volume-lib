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
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"

	"github.com/IBM/ibmcloud-storage-volume-lib/config"
	"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"

	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/auth"
)

// tokenGenerator ...
type tokenGenerator struct {
	config *config.VPCProviderConfig

	tokenKID        string
	tokenTTL        time.Duration
	tokenBeforeTime time.Duration

	privateKey *rsa.PrivateKey // Secret. Do not export
}

// readConfig ...
func (tg *tokenGenerator) readConfig(logger zap.Logger) (err error) {
	logger.Info("Entering readConfig")
	defer func() {
		logger.Info("Exiting readConfig", zap.Duration("tokenTTL", tg.tokenTTL), zap.Duration("tokenBeforeTime", tg.tokenBeforeTime), zap.String("tokenKID", tg.tokenKID), local.ZapError(err))
	}()

	if tg.privateKey != nil {
		return
	}

	path := filepath.Join(config.GetEtcPath(), tg.tokenKID)

	pem, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("Error reading PEM", local.ZapError(err))
		return
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		logger.Error("Error parsing PEM", local.ZapError(err))
		return
	}

	tg.privateKey = privateKey

	return
}

// buildToken ...
func (tg *tokenGenerator) buildToken(contextCredentials provider.ContextCredentials, ts time.Time, logger zap.Logger) (token *jwt.Token, err error) {
	logger.Info("Entering getJWTToken", zap.Reflect("contextCredentials", contextCredentials))
	defer func() {
		logger.Info("Exiting getJWTToken", zap.Reflect("token", token), local.ZapError(err))
	}()

	err = tg.readConfig(logger)
	if err != nil {
		return
	}

	claims := jwt.MapClaims{
		"iss": "armada",
		"exp": ts.Add(tg.tokenTTL).Unix(),
		"nbf": ts.Add(tg.tokenBeforeTime).Unix(),
		"iat": ts.Unix(),
	}

	switch {
	case contextCredentials.UserID == "":
		errStr := "User ID is not configured"
		logger.Error(errStr)
		err = errors.New(errStr)
		return

	case contextCredentials.AuthType == auth.IMSToken:
		claims["ims_user_id"] = contextCredentials.UserID

	default:
		claims["ims_username"] = contextCredentials.UserID

	}

	token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = tg.tokenKID

	return
}

// getServiceToken ...
func (tg *tokenGenerator) getServiceToken(contextCredentials provider.ContextCredentials, logger zap.Logger) (signedToken *string, err error) {
	token, err := tg.buildToken(contextCredentials, time.Now(), logger)
	if err != nil {
		return
	}

	signedString, err := token.SignedString(tg.privateKey)
	if err != nil {
		return
	}

	signedToken = &signedString

	return
}
