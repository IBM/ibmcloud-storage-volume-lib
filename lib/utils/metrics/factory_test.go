/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package metrics

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewFactory(t *testing.T) {

	var testcases = []struct {
		namespace       string
		metricsConfig   *Config
		expectedFactory Factory
		expectedError   string
	}{
		{
			namespace:     "my_default_namespace",
			metricsConfig: nil,
			expectedFactory: &factory{
				MaxAge:     60 * time.Minute,
				Namespace:  "my_default_namespace",
				gauges:     map[string]GaugeMetric{},
				summaries:  map[string]SummaryMetric{},
				histograms: map[string]HistogramMetric{},
			},
		},
		{
			namespace: "my_default_namespace",
			metricsConfig: &Config{
				MaxAge: "2m",
			},
			expectedFactory: &factory{
				MaxAge:     2 * time.Minute,
				Namespace:  "my_default_namespace",
				gauges:     map[string]GaugeMetric{},
				summaries:  map[string]SummaryMetric{},
				histograms: map[string]HistogramMetric{},
			},
		},
		{
			namespace: "my_default_namespace",
			metricsConfig: &Config{
				Namespace: "something else",
			},
			expectedFactory: &factory{
				MaxAge:     60 * time.Minute,
				Namespace:  "something else",
				gauges:     map[string]GaugeMetric{},
				summaries:  map[string]SummaryMetric{},
				histograms: map[string]HistogramMetric{},
			},
		},
		{
			namespace: "my_default_namespace",
			metricsConfig: &Config{
				MaxAge:    "6s",
				Namespace: "other",
			},
			expectedFactory: &factory{
				MaxAge:     6 * time.Second,
				Namespace:  "other",
				gauges:     map[string]GaugeMetric{},
				summaries:  map[string]SummaryMetric{},
				histograms: map[string]HistogramMetric{},
			},
		},
		{
			namespace: "my_default_namespace",
			metricsConfig: &Config{
				MaxAge: "nonsense",
			},
			expectedError: "time: invalid duration nonsense",
		},
		{
			expectedError: "Metrics namespace is not specified",
		},
	}

	for i, testcase := range testcases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			factory, actualErr := NewFactory(testcase.namespace, testcase.metricsConfig)

			assert.Equal(t, testcase.expectedFactory, factory)

			if testcase.expectedError == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Equal(t, testcase.expectedError, actualErr.Error())
			}
		})
	}
}
