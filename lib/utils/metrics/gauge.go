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
	"github.com/prometheus/client_golang/prometheus"
)

//go:generate counterfeiter -o fakes/gauge_metric.go --fake-name GaugeMetric . GaugeMetric

// GaugeMetric is a gauge metric
type GaugeMetric interface {
	Metric

	// Add increments the gauge by i
	Add(i float64, labels MetricLabels)

	// Sub decrements the gauge by i
	Sub(i float64, labels MetricLabels)
}

type gaugeMetric struct {
	isRegistered bool
	*prometheus.GaugeVec
}

// GetGaugeMetric implements Factory
func (f *factory) GetGaugeMetric(properties MetricProperties) (gm GaugeMetric, err error) {
	if err = properties.validate(); err != nil {
		return
	}

	key := properties.key()
	gm = f.gauges[key]

	if gm == nil {
		gm = &gaugeMetric{
			GaugeVec: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name:      properties.Name,
					Subsystem: properties.Subsystem,
					Namespace: f.Namespace,
					Help:      properties.Help,
				},
				properties.Labels),
		}

		f.gauges[key] = gm
	}

	return
}

// RegisterMetric implements Metric
func (gm *gaugeMetric) RegisterMetric(registerer prometheus.Registerer) {
	if !gm.isRegistered {
		registerer.MustRegister(gm.GaugeVec)
	}
	gm.isRegistered = true
}

// Add implements GaugeMetric
func (gm *gaugeMetric) Add(i float64, labels MetricLabels) {
	gm.GaugeVec.With(labels.GetMappedLabels()).Add(i)
}

// Add implements GaugeMetric
func (gm *gaugeMetric) Sub(i float64, labels MetricLabels) {
	gm.GaugeVec.With(labels.GetMappedLabels()).Sub(i)
}
