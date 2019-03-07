/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/provider"
)

func TestGetMemorySizeInGB(t *testing.T) {

	testCases := []struct {
		value string

		expectedValue int
		expectedError string
	}{
		{
			value:         "1024MB",
			expectedValue: 1,
		}, {
			value:         "16GB",
			expectedValue: 16,
		}, {
			value:         "1000TRIBBLES",
			expectedError: "Invalid format for config value: 1000TRIBBLES",
		}, {
			value:         "1000",
			expectedError: "Invalid format for config value: 1000",
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.value, func(t *testing.T) {

			val, err := GetMemorySizeInGB(testCase.value)

			if testCase.expectedError != "" {
				assert.Equal(t, 0, val)
				assert.EqualError(t, err, testCase.expectedError)
			} else {
				assert.Equal(t, testCase.expectedValue, val)
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetNetworkPortSpeedInMbps(t *testing.T) {

	testCases := []struct {
		value string

		expectedValue int
		expectedError string
	}{
		{
			value:         "1000MBPS",
			expectedValue: 1000,
		}, {
			value:         "1000GBPS",
			expectedValue: 1000000,
		}, {
			value:         "1000TRIBBLES",
			expectedError: "Invalid format for config value: 1000TRIBBLES",
		}, {
			value:         "1000",
			expectedError: "Invalid format for config value: 1000",
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.value, func(t *testing.T) {

			val, err := GetNetworkPortSpeedInMbps(testCase.value)

			if testCase.expectedError != "" {
				assert.Equal(t, 0, val)
				assert.EqualError(t, err, testCase.expectedError)
			} else {
				assert.Equal(t, testCase.expectedValue, val)
				assert.Nil(t, err)
			}
		})
	}
}

func TestProcessDiskConfig(t *testing.T) {

	testcases := []struct {
		name string

		value string

		expectedConfigs int
		expectErr       bool
	}{
		{
			name:            "OK 1 Drive",
			value:           "100GB_SATA",
			expectedConfigs: 1,
		},
		{
			name:            "OK 2 Drives",
			value:           "100GB_SATA;   1000GB_SSD",
			expectedConfigs: 2,
		},
		{
			name:            "OK 3 Drives",
			value:           "100GB_SATA;   1000GB_SSD;2x_56GB_RAID1",
			expectedConfigs: 3,
		},

		// RAID 10 validation
		{
			name: "InalidNoCountRAID10",

			value:     "2000GB_RAID10",
			expectErr: true,
		},
		{
			name:      "InvalidCountRAID10NotDivisibleBy2",
			value:     "5x_4000GB_RAID10",
			expectErr: true,
		},
		{
			name:      "InvalidCountRAID10LessThan4",
			value:     "2x_4000GB_RAID10",
			expectErr: true,
		},
		{
			name:      "InValidCountRAID1NotDivisibleBy2",
			value:     "3x_4000GB_RAID1",
			expectErr: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			dcs, err := ProcessDiskConfig(testcase.value, []string{"SSD", "SATA", "RAID1", "RAID10"})

			if testcase.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testcase.expectedConfigs, len(dcs))
			}
		})
	}
}

func Test_splitQuantityAndUnits(t *testing.T) {

	testCases := []struct {
		desc string

		value             string
		permittedSuffixes []string
		suffixOptional    bool

		expectedValue  int
		expectedSuffix string
		expectedError  bool
	}{
		{
			desc: "ValidNoSuffix",

			value:             "123",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    true,

			expectedValue:  123,
			expectedSuffix: "",
			expectedError:  false,
		}, {
			desc: "ValidSuffix1",

			value:             "123GB",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    true,

			expectedValue:  123,
			expectedSuffix: "GB",
			expectedError:  false,
		}, {
			desc: "ValidSuffix2",

			value:             "123gB",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    true,

			expectedValue:  123,
			expectedSuffix: "gB",
			expectedError:  false,
		}, {
			desc: "ValidSuffix3",

			value:             "123Gb",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    true,

			expectedValue:  123,
			expectedSuffix: "Gb",
			expectedError:  false,
		}, {
			desc: "ValidSuffix4",

			value:             "123gb",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    true,

			expectedValue:  123,
			expectedSuffix: "gb",
			expectedError:  false,
		}, {
			desc: "ValidSuffixRequired",

			value:             "123GB",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    false,

			expectedValue:  123,
			expectedSuffix: "GB",
			expectedError:  false,
		}, {
			desc: "ValidMultipleSuffixAllowed1",

			value:             "123GB",
			permittedSuffixes: []string{"GB", "MB"},
			suffixOptional:    true,

			expectedValue:  123,
			expectedSuffix: "GB",
			expectedError:  false,
		}, {
			desc: "ValidMultipleSuffixAllowed2",

			value:             "123GB",
			permittedSuffixes: []string{"GB", "MB"},
			suffixOptional:    false,

			expectedValue:  123,
			expectedSuffix: "GB",
			expectedError:  false,
		}, {
			desc: "InvalidNoSuffix",

			value:             "123",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    false,

			expectedError: true,
		}, {
			desc: "InvalidNoInt",

			value:             "GB",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    false,

			expectedError: true,
		}, {
			desc: "InvalidExtraCharsBefore",

			value:             "AA123GB",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    false,

			expectedError: true,
		}, {
			desc: "InvalidExtraCharsAfter",

			value:             "123GBAA",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    false,

			expectedError: true,
		}, {
			desc: "InvalidDecimal",

			value:             "123.123GB",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    false,

			expectedError: true,
		}, {
			desc: "InvalidDecimalNoSuffix",

			value:             "123.123",
			permittedSuffixes: []string{"GB"},
			suffixOptional:    false,

			expectedError: true,
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.desc, func(t *testing.T) {

			val, suff, err := splitQuantityAndUnits(testCase.value, testCase.permittedSuffixes, testCase.suffixOptional)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				if assert.Nil(t, err) {

					assert.Equal(t, testCase.expectedValue, val)
					assert.Equal(t, testCase.expectedSuffix, suff)
				}
			}
		})
	}
}

func Test_splitStorageConfiguration(t *testing.T) {

	testCases := []struct {
		desc string

		value string

		expectedCount    int
		expectedCapacity int
		expectedDiskType string
		expectedError    bool
	}{
		{
			desc: "ValidNoCountSATA",

			value:            "1000GB_SATA",
			expectedCount:    1,
			expectedCapacity: 1000,
			expectedDiskType: "SATA",
			expectedError:    false,
		},
		{
			desc: "ValidNoCountSATA2",

			value:            "2000GB_SATA",
			expectedCount:    1,
			expectedCapacity: 2000,
			expectedDiskType: "SATA",
			expectedError:    false,
		},
		{
			desc: "ValidNoCountSSD",

			value:            "2000GB_SSD",
			expectedCount:    1,
			expectedCapacity: 2000,
			expectedDiskType: "SSD",
			expectedError:    false,
		},
		{
			desc: "InvalidDiskType",

			value:         "2000GB_RAID",
			expectedError: true,
		},
		{
			desc: "WithCount",

			value:            "4x_4000GB_RAID10",
			expectedCount:    4,
			expectedCapacity: 4000,
			expectedDiskType: "RAID10",
			expectedError:    false,
		},
		{
			desc: "InvalidCount",

			value:         "ax_4000GB_RAID10",
			expectedError: true,
		},
		{
			desc: "InvalidUnits",

			value:         "2x 4000TB RAID10",
			expectedError: true,
		},
		{
			desc: "ValidMBConvertedToGB",

			value:            "1000MB_SATA",
			expectedCount:    1,
			expectedCapacity: 1,
			expectedDiskType: "SATA",
			expectedError:    false,
		},
		{
			desc: "ValidTBConvertedToGB",

			value:            "1TB_SATA",
			expectedCount:    1,
			expectedCapacity: 1000,
			expectedDiskType: "SATA",
			expectedError:    false,
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.desc, func(t *testing.T) {

			dc, err := splitStorageConfiguration(testCase.value, []string{"SSD", "SATA", "RAID1", "RAID10"})

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				if assert.Nil(t, err) {

					assert.Equal(t, testCase.expectedCount, dc.Count)
					assert.Equal(t, testCase.expectedCapacity, dc.Size)
					assert.Equal(t, testCase.expectedDiskType, dc.DiskType)
				}
			}
		})
	}
}

func Test_ProcessGPUConfig(t *testing.T) {

	testcases := []struct {
		name          string
		configuration string
		expectedKeys  []string
		expectedErr   error
	}{
		{
			name:          "success_1_type_1_gpu",
			configuration: "1x_K80",
			expectedKeys:  []string{"K80"},
		},
		{
			name:          "success_1_type_2_gpu",
			configuration: "2x_K80",
			expectedKeys:  []string{"K80", "K80"},
		},
		{
			name:          "success_2_types_1_gpu",
			configuration: "1x_K80;1x_V100",
			expectedKeys:  []string{"K80", "V100"},
		},
		{
			name:          "success_2_types_mixed_gpu",
			configuration: "1x_K80;2x_V100",
			expectedKeys:  []string{"K80", "V100", "V100"},
		},
		{
			name:          "success_3_types_1_gpu",
			configuration: "1x_K80;1x_V100;1x_P100",
			expectedKeys:  []string{"K80", "V100", "P100"},
		},
		{
			name:          "failure_bad_format",
			configuration: "x_K80:2xV100",
			expectedErr:   provider.Error{Fault: provider.Fault{Message: "Invalid format for GPU configuration: x_K80:2xV100", ReasonCode: "ErrorUnreconciledMachineConfig"}},
		},
		{
			name:        "failure_no_configuration",
			expectedErr: provider.Error{Fault: provider.Fault{Message: "Invalid format for GPU configuration: ", ReasonCode: "ErrorUnreconciledMachineConfig"}},
		},
		{
			name:          "failure_unknown_type",
			configuration: "1x_V80",
			expectedErr:   provider.Error{Fault: provider.Fault{Message: "Invalid format for GPU configuration: 1x_V80", ReasonCode: "ErrorUnreconciledMachineConfig"}},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			actualKeys, err := ProcessGPUConfig(testcase.configuration, []string{"K80", "V100", "P100"})
			if testcase.expectedErr != nil {
				assert.Equal(t, testcase.expectedErr, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, testcase.expectedKeys, actualKeys)
		})
	}
}
