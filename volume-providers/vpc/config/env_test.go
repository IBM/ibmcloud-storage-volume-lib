/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultEnvironmentClock(t *testing.T) {
	clock := DefaultEnvironmentClock

	result := clock.NowSecondUTC()
	now := time.Now()
	assert.NotEmpty(t, result)
	assert.True(t, result.Before(now))
	assert.WithinDuration(t, now, result, time.Second)
	assert.Equal(t, 0, result.Nanosecond())
	_, offset := result.Zone()
	assert.Equal(t, 0, offset)

	a := clock.NowUnixNano()
	time.Sleep(time.Nanosecond * 100)
	b := clock.NowUnixNano()
	assert.NotEqual(t, a, b)

	assert.NotEmpty(t, clock.Visitor())
}
