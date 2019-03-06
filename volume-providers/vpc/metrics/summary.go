/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//go:generate counterfeiter -o fakes/summary_metric.go --fake-name SummaryMetric . SummaryMetric

// SummaryMetric is a summary vector metric
type SummaryMetric interface {
	Metric

	// ObserveSinceStart observes the metric with the specified labels with the time since startTime
	ObserveSinceStart(labels MetricLabels, startTime time.Time)

	// ObserveBetween observes the metric with the specified labels with a start and end time
	ObserveBetween(labels MetricLabels, startTime time.Time, endTime time.Time)

	// ObserveValue observes the metric with the specified labels with given value
	ObserveValue(labels MetricLabels, value float64)
}

type summaryMetric struct {
	isRegistered bool
	*prometheus.SummaryVec
}

var _ SummaryMetric = &summaryMetric{}

// GetSummaryMetric implements Factory
func (f *factory) GetSummaryMetric(properties MetricProperties) (sm SummaryMetric, err error) {
	if err = properties.validate(); err != nil {
		return
	}

	key := properties.key()
	sm = f.summaries[key]

	if sm == nil {
		sm = &summaryMetric{
			SummaryVec: prometheus.NewSummaryVec(
				prometheus.SummaryOpts{
					Name:      properties.Name,
					Subsystem: properties.Subsystem,
					Namespace: f.Namespace,
					Help:      properties.Help,
					MaxAge:    f.MaxAge,
				},
				properties.Labels),
		}

		f.summaries[key] = sm
	}

	return
}

// RegisterMetric implements Metric
func (sm *summaryMetric) RegisterMetric(registerer prometheus.Registerer) {
	if !sm.isRegistered {
		registerer.MustRegister(sm.SummaryVec)
	}
	sm.isRegistered = true
}

// ObserveSinceStart implements SummaryMetric
func (sm *summaryMetric) ObserveSinceStart(labels MetricLabels, startTime time.Time) {
	sm.SummaryVec.With(labels.GetMappedLabels()).Observe(time.Since(startTime).Seconds())
}

// ObserveBetween implements SummaryMetric
func (sm *summaryMetric) ObserveBetween(labels MetricLabels, startTime time.Time, endTime time.Time) {
	sm.SummaryVec.With(labels.GetMappedLabels()).Observe(endTime.Sub(startTime).Seconds())
}

// ObserveValue implements SummaryMetric
func (sm *summaryMetric) ObserveValue(labels MetricLabels, value float64) {
	sm.SummaryVec.With(labels.GetMappedLabels()).Observe(value)
}
