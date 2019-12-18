/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Kubernetes Service, 5737-D43
 * (C) Copyright IBM Corp. 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package metrics

import (
	"go.uber.org/zap"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// FunctionLabel is a name of CSI plugin operation for which
// we measure duration
type FunctionLabel string

const (
	pluginNamespace = "ibmcloud_storage_volume_lib"
)

var (
	functionDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: pluginNamespace,
			Name:      "function_duration_seconds",
			Help:      "Time taken by various operation of library",
			//Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1.0, 2.5, 5.0, 7.5, 10.0, 12.5, 15.0, 17.5, 20.0, 22.5, 25.0, 27.5, 30.0, 50.0, 75.0, 100.0, 1000.0},
		}, []string{"function"},
	)
	functionCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: pluginNamespace,
			Name:      "functions_total",
			Help:      "The number of library operation  completeted successfully.",
		}, []string{"function"},
	)

	errorsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: pluginNamespace,
			Name:      "errors_total",
			Help:      "The number of library operation  failed due to an error.",
		}, []string{"type"},
	)
)

// RegisterAll registers all metrics.
func RegisterAll() {
	prometheus.MustRegister(functionDuration)
	prometheus.MustRegister(functionCount)
	prometheus.MustRegister(errorsCount)
}

// UpdateDurationFromStart records the duration of the step identified by the
// label using start time
func UpdateDurationFromStart(logger *zap.Logger, label string, start time.Time) {
	duration := time.Now().Sub(start)
	logger.Info("Time to complete", zap.Float64(string(label), duration.Seconds()))
	UpdateDuration(label, duration)
}

// UpdateDuration records the duration of the step identified by the label
func UpdateDuration(label string, duration time.Duration) {
	functionDuration.WithLabelValues(label).Set(duration.Seconds())
}

// RegisterError records any errors for any lib operation.
func RegisterError(errType string, err error) {
	if err != nil {
		errType = err.Error() // TODO Get the error code
	}
	errorsCount.WithLabelValues(errType).Add(1.0)
}

// RegisterFunction records number of operation.
func RegisterFunction(label string) {
	functionCount.WithLabelValues(string(label)).Add(1.0)
}
