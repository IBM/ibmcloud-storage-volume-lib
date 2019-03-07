/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package metrics

import (
	"errors"
	"strconv"
	"testing"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

//go:generate counterfeiter -o fakes/prometheus/registerer.go --fake-name Registerer . registerer
type registerer interface {
	prometheus.Registerer
}

func Test_key_validate(t *testing.T) {
	var testcases = []struct {
		metricsProperties MetricProperties
		expectedKey       string
		expectedError     error
	}{
		{
			metricsProperties: MetricProperties{},
			expectedKey:       ":",
			expectedError:     errors.New("No name provided for metric"),
		},
		{
			metricsProperties: MetricProperties{
				Name: "a",
			},
			expectedKey:   ":a",
			expectedError: errors.New("No subsystem provided for metric"),
		},
		{
			metricsProperties: MetricProperties{
				Name:      "a",
				Subsystem: "b",
			},
			expectedKey:   "b:a",
			expectedError: errors.New("No help provided for metric"),
		},
		{
			metricsProperties: MetricProperties{
				Name:      "a",
				Subsystem: "b",
				Help:      "c",
			},
			expectedKey:   "b:a",
			expectedError: errors.New("No labels provided for metric"),
		},
		{
			metricsProperties: MetricProperties{
				Name:      "a",
				Subsystem: "b",
				Help:      "c",
				Labels:    []string{},
			},
			expectedKey: "b:a",
		},
		{
			metricsProperties: MetricProperties{
				Name:      "a",
				Subsystem: "b",
				Help:      "c",
				Labels:    []string{"1", "2", "3"},
			},
			expectedKey: "b:a",
		},
	}

	for i, testcase := range testcases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			assert.Equal(t, testcase.expectedKey, testcase.metricsProperties.key())

			actualErr := testcase.metricsProperties.validate()

			assert.Equal(t, testcase.expectedError, actualErr)
		})
	}
}

func Test_ReasonCodeToMetricsLabel(t *testing.T) {
	assert.Equal(t, OK, ReasonCodeToMetricsLabel(""))
	assert.Equal(t, "Fred", ReasonCodeToMetricsLabel(reasoncode.ReasonCode("Fred")))
}
