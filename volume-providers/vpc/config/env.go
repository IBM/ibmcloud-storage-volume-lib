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
	"os"
	"time"
)

// EnvironmentClock exposes envrionmental and temporal context information
// such that it can be overriden and fixed in unit tests
//go:generate counterfeiter -o fakes/env.go --fake-name EnvironmentClock . EnvironmentClock
type EnvironmentClock interface {
	// NowSecondUTC returns current time, rounded down to nearest second (in UTC)
	NowSecondUTC() time.Time

	// NowUnixNano returns the current time in nano-seconds
	NowUnixNano() int64

	// Sleep for the duration, but not unit test
	Sleep(d time.Duration)

	// GitCommit returns the GIT commit SHA, if known
	GitCommit() string

	// BuildNumber returns the build number, if known
	BuildNumber() string

	// Visitor returns the "visitor" name, based on the pod name and build number
	Visitor() string
}

type envClock struct{}

// NowSecondUTC is for internal use by armada-cluster
func (envClock) NowSecondUTC() time.Time {
	// Round down to the nearest second, to make sure we don't
	// carry any more precision than is preserved in model.TimeFormat
	// If we don't then we risk sending different values to
	// the metering service!
	return time.Now().Truncate(time.Second).UTC()
}

// NowUnixNano is for internal use by armada-cluster
func (envClock) NowUnixNano() int64 {
	return time.Now().UnixNano()
}

// Sleep is for internal use by armada-cluster
func (envClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// GitCommit is for internal use by armada-cluster
func (envClock) GitCommit() string {
	return os.Getenv("GIT_COMMIT_SHA")
}

// BuildNumber is for internal use by armada-cluster
func (envClock) BuildNumber() string {
	return os.Getenv("BUILD_NUMBER")
}

// Visitor is for internal use by armada-cluster
func (c envClock) Visitor() string {
	pod, err := os.Hostname()
	if err != nil {
		pod = "unknown"
	}
	build := c.BuildNumber()
	if build != "" {
		return pod + ":" + build
	}
	return pod
}

// DefaultEnvironmentClock is for internal use by armada-cluster
var DefaultEnvironmentClock EnvironmentClock = envClock{}
