/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

// Package metrics helps define and record metrics using Prometheus
package metrics

import (
	"errors"
	"time"
)

// DefaultMaxAge is 60 minutes
const DefaultMaxAge = 60 * time.Minute

// Config contains configuration for NewFactory()
type Config struct {
	// MaxAge is expressed as a time.Duration, e.g. "60m"
	MaxAge string `envconfig:"METRICS_MAX_AGE"`

	// Namespace is the namespace for all the metrics within this microservice
	Namespace string `envconfig:"METRICS_NAMESPACE"`
}

//go:generate counterfeiter -o fakes/factory.go --fake-name Factory . Factory

// Factory creates metric instances
type Factory interface {
	// GetSummaryMetric prepares or retrieves a singleton summary vector metric.
	// When getting previously prepared metric, the help and labels are unchanged.
	GetSummaryMetric(properties MetricProperties) (SummaryMetric, error)
	// GetGaugeMetric prepares or retrieves a singleton gauge metric.
	// When getting previously prepared metric, the help and labels are unchanged.
	GetGaugeMetric(properties MetricProperties) (GaugeMetric, error)
	// GetHistogramMetric prepares or retrieves a singleton histogram metric using the
	// buckets specified by an array of strictly increasing float values.
	// When getting previously prepared metric, the buckets, help and labels are unchanged.
	GetHistogramMetric(properties MetricProperties, buckets []float64) (HistogramMetric, error)
}

type factory struct {
	MaxAge    time.Duration
	Namespace string

	gauges     map[string]GaugeMetric
	summaries  map[string]SummaryMetric
	histograms map[string]HistogramMetric
}

var _ Factory = &factory{}

// NewFactory creates a new metrics Factory using the specified Prometheus namespace.
// The Config argument is optional; if supplied then the specified namespace
// is overridden if Config.Namespace is non-empty.
func NewFactory(namespace string, conf *Config) (result Factory, err error) {
	maxAge := DefaultMaxAge

	if conf != nil {

		if conf.MaxAge != "" {
			maxAge, err = time.ParseDuration(conf.MaxAge)
			if err != nil {
				return
			}
		}

		if conf.Namespace != "" {
			namespace = conf.Namespace
		}

	}

	if namespace == "" {
		err = errors.New("Metrics namespace is not specified")
		return
	}

	result = &factory{
		Namespace:  namespace,
		MaxAge:     maxAge,
		gauges:     map[string]GaugeMetric{},
		summaries:  map[string]SummaryMetric{},
		histograms: map[string]HistogramMetric{},
	}

	return
}
