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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseBucketsConfiguration(t *testing.T) {
	testcases := []struct {
		name            string
		bucketsStr      string
		expectedBuckets []float64
		expectedError   bool
	}{
		{
			name:            "empty_array",
			bucketsStr:      "",
			expectedBuckets: []float64{},
		},
		{
			name:          "invalid_buckets_string",
			bucketsStr:    "this_is_bollocks",
			expectedError: true,
		},
		{
			name:            "valid_buckets",
			bucketsStr:      "0.1,0.2,0.3",
			expectedBuckets: []float64{0.1, 0.2, 0.3},
		},
		{
			name:            "valid_buckets_but_need_sorting",
			bucketsStr:      "0.2,0.1,0.3",
			expectedBuckets: []float64{0.1, 0.2, 0.3},
		},
		{
			name:            "valid_buckets_with_spaces",
			bucketsStr:      "0.1 , 0.2, 0.3",
			expectedBuckets: []float64{0.1, 0.2, 0.3},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {

			actualBuckets, err := ParseBucketsConfiguration(testcase.bucketsStr)

			if testcase.expectedError {
				assert.NotNil(t, err)
				assert.Empty(t, actualBuckets)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, testcase.expectedBuckets, actualBuckets)
			}

		})
	}
}
