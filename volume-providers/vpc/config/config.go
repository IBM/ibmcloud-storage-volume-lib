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
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// GetStringList takes a string argument, removes whitespace and then splits the string content by ","
// returning the output as an array.
func GetStringList(val string) []string {
	val = strings.Replace(val, " ", "", -1)
	return strings.Split(val, ",")
}

func getEnv(key string) string {
	return os.Getenv(strings.ToUpper(key))
}

// GetGoPath inspects the environment for the GOPATH variable
func GetGoPath() string {
	if goPath := getEnv("GOPATH"); goPath != "" {
		return goPath
	}
	return ""
}

// LoadPrefixVarConfigs is for internal use by armada-cluster
func LoadPrefixVarConfigs(mappings string, template interface{}, offer func(string, interface{})) (err error) {
	for _, mapping := range strings.Split(mappings, " ") { //  e.g. "vmware:VMWARE gt:GT"
		if mapping != "" {
			p := strings.Split(mapping, ":")
			if len(p) != 2 {
				err = errors.New("Invalid prefix config spec: " + mapping)
				return
			}

			c := reflect.New(reflect.ValueOf(template).Type()).Interface()

			err = envconfig.Process(p[1], c)
			if err != nil {
				return
			}

			offer(p[0], c)
		}
	}
	return
}
