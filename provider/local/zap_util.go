/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package local

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapError formats provider error messages in a useful way for logging,
// and performs the standard zap.Error on non provider errors
func ZapError(err error) zapcore.Field {
	return zap.Error(err)
}
