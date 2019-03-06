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
	"errors"
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	prometheusfakes "github.ibm.com/narkarum/ibmcloud-storage-volume-lib/volume-providers/vpc/metrics/fakes/prometheus"
)

func Test_GetHistogramMetric(t *testing.T) {

	var testcases = []struct {
		metricsProperties MetricProperties
		buckets           []float64
		expectedBuckets   []float64
		expectedError     error
	}{
		{
			metricsProperties: MetricProperties{},
			expectedError:     errors.New("No name provided for metric"),
		},
		{
			metricsProperties: MetricProperties{
				Name:      "abc",
				Subsystem: "def",
				Help:      "ghi",
				Labels:    []string{"j", "k", "l"},
			},
			buckets:         []float64{0.1, 0.2},
			expectedBuckets: []float64{0.1, 0.2},
		},
		{
			metricsProperties: MetricProperties{
				Name:      "abc",
				Subsystem: "def",
				Help:      "ghi",
				Labels:    []string{"j", "k", "l"},
			},
			buckets:         []float64{},
			expectedBuckets: prometheus.DefBuckets,
		},
	}

	for i, testcase := range testcases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			factory := factory{
				Namespace:  "a",
				histograms: map[string]HistogramMetric{},
			}

			hm, actualErr := factory.GetHistogramMetric(testcase.metricsProperties, testcase.buckets)

			assert.Equal(t, testcase.expectedError, actualErr)

			if actualErr == nil {
				hm2, actualErr := factory.GetHistogramMetric(testcase.metricsProperties, testcase.buckets)

				assert.True(t, hm == hm2)
				assert.NoError(t, actualErr)
			}

		})
	}
}

func Test_HistogramMetric_RegisterMetric(t *testing.T) {
	registerer := prometheusfakes.Registerer{}

	factory := factory{
		Namespace:  "a",
		histograms: map[string]HistogramMetric{},
	}

	buckets := []float64{0.1, 0.5, 0.9, 0.99}

	metric, err := factory.GetHistogramMetric(MetricProperties{
		Name:      "abc",
		Subsystem: "def",
		Help:      "ghi",
		Labels:    []string{"j", "k", "l"},
	}, buckets)
	assert.NoError(t, err)
	assert.NotNil(t, metric.(*histogramMetric).HistogramVec)

	metric.RegisterMetric(&registerer)
	metric.RegisterMetric(&registerer)
	metric.RegisterMetric(&registerer)

	if assert.Equal(t, 1, registerer.MustRegisterCallCount()) {
		assert.Equal(t, []prometheus.Collector{metric.(*histogramMetric).HistogramVec},
			registerer.MustRegisterArgsForCall(0))
	}
}
