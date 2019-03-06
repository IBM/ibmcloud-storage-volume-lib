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
	"github.com/prometheus/client_golang/prometheus"
)

//go:generate counterfeiter -o fakes/histogram_metric.go --fake-name HistogramMetric . HistogramMetric

// HistogramMetric is a histogram metric
type HistogramMetric interface {
	Metric

	// ObserveValue observes the metric with the specified labels with given value
	ObserveValue(labels MetricLabels, value float64)
}

type histogramMetric struct {
	isRegistered bool
	*prometheus.HistogramVec
}

var _ HistogramMetric = &histogramMetric{}

func (f *factory) GetHistogramMetric(properties MetricProperties, buckets []float64) (hm HistogramMetric, err error) {
	if err = properties.validate(); err != nil {
		return
	}

	key := properties.key()
	hm = f.histograms[key]

	if hm == nil {

		// Better to have default buckets than none...
		if len(buckets) == 0 {
			buckets = prometheus.DefBuckets
		}

		hm = &histogramMetric{
			HistogramVec: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:      properties.Name,
					Subsystem: properties.Subsystem,
					Namespace: f.Namespace,
					Help:      properties.Help,
					Buckets:   buckets,
				},
				properties.Labels),
		}

		f.histograms[key] = hm
	}

	return
}

// RegisterMetric implements Metric
func (hm *histogramMetric) RegisterMetric(registerer prometheus.Registerer) {
	if !hm.isRegistered {
		registerer.MustRegister(hm.HistogramVec)
	}
	hm.isRegistered = true
}

func (hm *histogramMetric) ObserveValue(labels MetricLabels, value float64) {
	hm.HistogramVec.With(labels.GetMappedLabels()).Observe(value)
}
