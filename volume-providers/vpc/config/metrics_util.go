/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package impl

import (
	"sort"
	"strconv"
	"strings"
)

// ParseBucketsConfiguration takes a comma separated string of float values,
// converts it to a sorted slice of float64
// for internal use by armada-cluster only
func ParseBucketsConfiguration(bucketsStr string) ([]float64, error) {
	if bucketsStr != "" {
		bs := strings.Split(bucketsStr, ",")

		var buckets []float64
		for _, b := range bs {
			b = strings.TrimSpace(b)
			fl, err := strconv.ParseFloat(b, 10)
			if err != nil {
				return []float64{}, err
			}
			buckets = append(buckets, fl)
		}
		sort.Float64s(buckets)
		return buckets, nil
	}
	return []float64{}, nil
}
