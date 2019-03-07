/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package provider

import (
	"strconv"

	"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/metrics"
	"go.uber.org/zap"
)

// VPCMetrics ...
type VPCMetrics struct {
	ResponseTime metrics.SummaryMetric
}

// Labels used for the VPC metrics
const operationLabel = "operation"            //TODO: extract and share with other providers
const reasonCodeLabel = "reason_code"         //TODO: extract and share with other providers
const authTypeLabel = "auth_type"             //TODO: extract and share with other providers
const defaultAccountLabel = "default_account" //TODO: extract and share with other providers

// ResponseTimeMetricLabels is an ordered list of all the of metrics for the ResponseTime metric
var ResponseTimeMetricLabels = []string{operationLabel, reasonCodeLabel, authTypeLabel, defaultAccountLabel} //TODO: extract and share with other providers

// ResponseTimeMetricInstance contains named fields for each of the ResponseTime metric labels
type ResponseTimeMetricInstance struct { //TODO: extract and share with other providers
	Operation      string
	ReasonCode     string
	AuthType       provider.AuthType
	DefaultAccount bool
}

// GetMappedLabels populates a map with the labels and their associated values
func (instance *ResponseTimeMetricInstance) GetMappedLabels() map[string]string { //TODO: extract and share with other providers
	return map[string]string{
		operationLabel:      instance.Operation,
		reasonCodeLabel:     instance.ReasonCode,
		authTypeLabel:       string(instance.AuthType),
		defaultAccountLabel: strconv.FormatBool(instance.DefaultAccount),
	}
}

func buildVPCMetrics(subsystem string, factory metrics.Factory, logger *zap.Logger) (vpcMetrics *VPCMetrics, err error) {

	responseTimeMetric, err := factory.GetSummaryMetric(metrics.MetricProperties{
		Name:      "response_time",
		Subsystem: subsystem,
		Help:      "Remote provider response time metric for vpc-provider",
		Labels:    ResponseTimeMetricLabels,
	})
	if err != nil {
		logger.Error("Unable to create response_time metric", zap.Error(err))
		return nil, err
	}

	return &VPCMetrics{
		ResponseTime: responseTimeMetric,
	}, nil
}
