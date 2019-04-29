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
	"go.uber.org/zap"
	"strings"
	"time"
)

var logger *zap.Logger

func retry(retryfunc func() error) error {
	var err error
	MAX_RETRY_ATTEMPT := 5
	MAX_RETRY_GAP := 30
	RETRY_GAP := 5

	for i := 0; i < MAX_RETRY_ATTEMPT; i++ {
		if i > 0 {
			time.Sleep(time.Duration(RETRY_GAP) * time.Second)
		}
		err = retryfunc()
		if err != nil {
			//Skip retry for the below type of Errors
			if (strings.Contains(err.Error(), "unable to find network storage associated")) || (strings.Contains(err.Error(), "is Already Authorized for host")) {
				break
			}
			if i >= 1 {
				RETRY_GAP = 2 * RETRY_GAP
				if RETRY_GAP > MAX_RETRY_GAP {
					RETRY_GAP = MAX_RETRY_GAP
				}
			}
			if (i + 1) < MAX_RETRY_ATTEMPT {
				logger.Info("Error while executing the function. Re-attempting execution ..", zap.Int("attempt..", i+2), zap.Int("retry-gap", RETRY_GAP), zap.Int("max-retry-Attempts", MAX_RETRY_ATTEMPT), zap.Error(err))
			}
			continue
		}
		return err
	}
	return err
}
