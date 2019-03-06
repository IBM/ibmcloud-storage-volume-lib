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

	"github.com/prometheus/client_golang/prometheus"
	prometheusfakes "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/metrics/fakes/prometheus"

	"github.com/stretchr/testify/assert"
)

func Test_GetSummaryMetric(t *testing.T) {

	var testcases = []struct {
		metricsProperties MetricProperties
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
		},
	}

	for i, testcase := range testcases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			factory := factory{
				Namespace: "a",
				MaxAge:    DefaultMaxAge,
				summaries: map[string]SummaryMetric{},
			}

			sm, actualErr := factory.GetSummaryMetric(testcase.metricsProperties)

			assert.Equal(t, testcase.expectedError, actualErr)

			if actualErr == nil {
				sm2, actualErr := factory.GetSummaryMetric(testcase.metricsProperties)

				assert.True(t, sm == sm2)
				assert.NoError(t, actualErr)
			}

		})
	}
}

func Test_SummaryMetric_RegisterMetric(t *testing.T) {
	registerer := prometheusfakes.Registerer{}

	factory := factory{
		Namespace: "a",
		MaxAge:    DefaultMaxAge,
		summaries: map[string]SummaryMetric{},
	}

	metric, err := factory.GetSummaryMetric(MetricProperties{
		Name:      "abc",
		Subsystem: "def",
		Help:      "ghi",
		Labels:    []string{"j", "k", "l"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, metric.(*summaryMetric).SummaryVec)

	metric.RegisterMetric(&registerer)
	metric.RegisterMetric(&registerer)
	metric.RegisterMetric(&registerer)

	if assert.Equal(t, 1, registerer.MustRegisterCallCount()) {
		assert.Equal(t, []prometheus.Collector{metric.(*summaryMetric).SummaryVec},
			registerer.MustRegisterArgsForCall(0))
	}
}
