/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package iam

import (
	"errors"

	"github.com/dgrijalva/jwt-go"

	"go.uber.org/zap"
)

type accessTokenClaims struct {
	jwt.StandardClaims

	Account struct {
		Bss string `json:"bss"`
	} `json:"account"`
}

func (r *tokenExchangeService) GetIAMAccountIDFromAccessToken(accessToken AccessToken, logger *zap.Logger) (accountID string, err error) {

	// TODO - TEMPORARY CODE - VERIFY SIGNATURE HERE
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken.Token, &accessTokenClaims{})
	if err != nil {
		return
	}
	token.Valid = true
	// TODO - TEMPORARY CODE - DONT OVERRIDE VERIFICATION

	claims, haveClaims := token.Claims.(*accessTokenClaims)

	logger.Debug("Access token parsed", zap.Bool("haveClaims", haveClaims), zap.Bool("valid", token.Valid))

	if !token.Valid || !haveClaims {
		err = errors.New("Access token invalid")
		return
	}

	accountID = claims.Account.Bss
	logger.Debug("GetIAMAccountIDFromAccessToken", zap.Reflect("claims.Account.Bss", claims.Account.Bss))

	return
}
