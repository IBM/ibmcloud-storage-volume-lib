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
	"fmt"
	"strconv"
	"strings"
	"time"
)

// maxRetryAttempt ...
var maxRetryAttempt = 5

// maxRetryGap ...
var maxRetryGap = 30

// retryGap ...
var retryGap = 5

// retry ...
func retry(retryfunc func() error) error {
	var err error
	for i := 0; i < maxRetryAttempt; i++ {
		if i > 0 {
			time.Sleep(time.Duration(retryGap) * time.Second)
		}
		err = retryfunc()
		if err != nil {
			//Skip retry for the below type of Errors
			if (strings.Contains(err.Error(), "unable to find network storage associated")) || (strings.Contains(err.Error(), "is Already Authorized for host")) {
				break
			}
			if i >= 1 {
				retryGap = 2 * retryGap
				if retryGap > maxRetryGap {
					retryGap = maxRetryGap
				}
			}
			if (i + 1) < maxRetryAttempt {
				fmt.Printf("\nReattenmpting execution func: %#v, attempt =%d,  max attepmt = %d ,error %#v", retryfunc, i+2, maxRetryAttempt, err) // TODO: need to use logger
				//c.logger.Info("Error while executing the function. Re-attempting execution ..", zap.Int("attempt..", i+2), zap.Int("retry-gap", retryGap), zap.Int("max-retry-Attempts", maxRetryAttempt), zap.Error(err))
			}
			continue
		}
		return err
	}
	return err
}

// ToInt ...
func ToInt(valueInInt string) int {
	value, err := strconv.Atoi(valueInInt)
	if err != nil {
		return 0
	}
	return value
}

// ToInt64 ...
func ToInt64(valueInInt string) int64 {
	value, err := strconv.ParseInt(valueInInt, 10, 64)
	if err != nil {
		return 0
	}
	return value
}
