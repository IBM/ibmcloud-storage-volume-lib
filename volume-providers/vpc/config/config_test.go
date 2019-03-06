/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetStringList(t *testing.T) {
	assert.Equal(t, []string{""}, GetStringList(""))
	assert.Equal(t, []string{"hello"}, GetStringList("hello"))
	assert.Equal(t, []string{"hello", "how", "are", "you?"}, GetStringList("   hello, how     , are, you?   "))
}

func TestGetGoPath(t *testing.T) {
	originalValue := os.Getenv("GOPATH")

	os.Unsetenv("GOPATH")
	actual := GetGoPath()
	assert.Equal(t, "", actual)
	expected := "/tmp"
	os.Setenv("GOPATH", expected)
	actual = GetGoPath()
	assert.Equal(t, expected, actual)
	os.Unsetenv("GOPATH")

	// Reset back to original value
	if originalValue == "" {
		os.Unsetenv("GOPATH")
	} else {
		os.Setenv("GOPATH", originalValue)
	}
}

func Test_LoadPrefixVarConfigs(t *testing.T) {

	type NestedConfig struct {
		NestedValue string `envconfig:"NESTED_VAR"`
	}

	type TestConfig struct {
		NestedConfig
		StringValue string `envconfig:"STRING_VAR"`
		TrueValue   bool   `envconfig:"TRUE_VAR"`
		FalseValue  bool   `envconfig:"FALSE_VAR"`
	}

	testcases := []struct {
		name           string
		mappings       string
		envVars        map[string]string
		expectedResult map[string]TestConfig
		expectedError  string
	}{{
		name:           "empty",
		expectedResult: map[string]TestConfig{},
	}, {
		name:     "single_no_content",
		mappings: "fred:FRED",
		expectedResult: map[string]TestConfig{
			"fred": TestConfig{},
		},
	}, {
		name:     "single_with_content",
		mappings: "fred:FRED",
		envVars: map[string]string{
			"FRED_STRING_VAR": "value1",
		},
		expectedResult: map[string]TestConfig{
			"fred": TestConfig{StringValue: "value1"},
		},
	}, {
		name:     "nested",
		mappings: "fred:FRED",
		envVars: map[string]string{
			"FRED_STRING_VAR": "My string value",
			"FRED_TRUE_VAR":   "true",
			"FRED_FALSE_VAR":  "false",
			"FRED_NESTED_VAR": "My nested value",
		},
		expectedResult: map[string]TestConfig{
			"fred": TestConfig{
				NestedConfig: NestedConfig{NestedValue: "My nested value"},
				StringValue:  "My string value",
				TrueValue:    true,
			},
		},
	}, {
		name:     "multi_with_content",
		mappings: "fred:FRED pete:PETE",
		envVars: map[string]string{
			"FRED_STRING_VAR": "value1",
			"PETE_STRING_VAR": "value2",
		},
		expectedResult: map[string]TestConfig{
			"fred": TestConfig{StringValue: "value1"},
			"pete": TestConfig{StringValue: "value2"},
		},
	}, {
		name:           "error",
		mappings:       "fred:FRED:PETE",
		expectedResult: map[string]TestConfig{},
		expectedError:  "Invalid prefix config spec: fred:FRED:PETE",
	}}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {

			for key, val := range testcase.envVars {
				os.Setenv(key, val)
			}

			defer func() {
				for key := range testcase.envVars {
					os.Unsetenv(key)
				}
			}()

			result := map[string]TestConfig{}

			err := LoadPrefixVarConfigs(testcase.mappings, TestConfig{}, func(name string, value interface{}) {
				result[name] = *value.(*TestConfig)
			})

			assert.Equal(t, testcase.expectedResult, result)

			if testcase.expectedError == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Equal(t, testcase.expectedError, err.Error())
			}
		})
	}
}
