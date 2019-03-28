/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/stretchr/testify/assert"
)

// SetupServer ...
func SetupServer(t *testing.T) (m *http.ServeMux, c client.Client, teardown func()) {

	m = http.NewServeMux()
	s := httptest.NewServer(m)

	log := new(bytes.Buffer)

	c = client.New(s.URL, http.DefaultClient).WithDebug(log).WithAuthToken("auth-token")

	teardown = func() {
		s.Close()
		CheckTestFail(t, log)

	}

	return
}

// SetupMuxResponse ...
func SetupMuxResponse(t *testing.T, m *http.ServeMux, path string, expectedMethod string, expectedContent *string, statusCode int, body string, verify func(t *testing.T, r *http.Request)) {

	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, expectedMethod, r.Method)

		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer auth-token", authHeader)

		acceptHeader := r.Header.Get("Accept")
		assert.Equal(t, "application/json", acceptHeader)

		if expectedContent != nil {
			b, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, *expectedContent, string(b))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if body != "" {
			fmt.Fprint(w, body)
		}

		if verify != nil {
			verify(t, r)
		}
	})
}

// CheckTestFail ...
func CheckTestFail(t *testing.T, buf *bytes.Buffer) {

	if t.Failed() {
		t.Log(buf)
	}
}

// Sptr ...
func Sptr(s string) *string {
	return &s
}
