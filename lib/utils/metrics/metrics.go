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
	"fmt"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"

	"github.com/prometheus/client_golang/prometheus"
)

// OK is a pseudo-ReasonCode value to indicate success
const OK = "OK"

// ReasonCodeToMetricsLabel maps a ReasonCode to a string, defaulting to "OK"
// for the null ReasonCode value
func ReasonCodeToMetricsLabel(code reasoncode.ReasonCode) string {
	if code == "" {
		return OK
	}
	return string(code)
}

// Metric represents a metric instance
type Metric interface {

	// RegisterMetric is used to register this metric with Prometheus
	RegisterMetric(registerer prometheus.Registerer)
}

// MetricLabels represents a the labels for an individual observation
type MetricLabels interface {

	// GetMappedLabels returns the labels map for this observation
	GetMappedLabels() map[string]string
}

// MetricProperties defines a new metric
type MetricProperties struct {

	// Name is the Prometheus metric name
	Name string

	// Subsystem is the Prometheus metric subsystem
	Subsystem string

	// Help is the Prometheus metric help text
	Help string

	// Lables defines the Prometheus metric labels
	Labels []string
}

func (mp MetricProperties) key() string {
	return fmt.Sprintf("%v:%v", mp.Subsystem, mp.Name)
}

func (mp MetricProperties) validate() error {
	if mp.Name == "" {
		return errors.New("No name provided for metric")
	}

	if mp.Subsystem == "" {
		return errors.New("No subsystem provided for metric")
	}

	if mp.Help == "" {
		return errors.New("No help provided for metric")
	}

	if mp.Labels == nil {
		return errors.New("No labels provided for metric")
	}
	return nil
}
