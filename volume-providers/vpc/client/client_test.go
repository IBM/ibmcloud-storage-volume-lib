/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package client_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/client/models"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/client/riaas/test"
)

var getOperation = &client.Operation{
	Name:        "GetOperation",
	Method:      "GET",
	PathPattern: "/resource",
}
var postOperation = &client.Operation{
	Name:        "PostOperation",
	Method:      "POST",
	PathPattern: "/resource",
}

func TestClient(t *testing.T) {

	var (
		request   *client.Request
		errResult models.Error
	)

	testcases := []struct {
		name string

		operation *client.Operation

		modifyRequest func()

		requestBody *string

		responseBody string
		responseCode int

		expectErr string
		verify    func(t *testing.T)
		muxVerify func(*testing.T, *http.Request)
	}{
		{
			name:      "creates invokable requests from static operations (GET)",
			operation: getOperation,
		}, {
			name:      "creates invokable requests from static operations (POST)",
			operation: postOperation,
		}, {
			name:      "encodes query parameters",
			operation: getOperation,
			modifyRequest: func() {
				request = request.AddQueryValue("name", "value1").AddQueryValue("name", "value2").AddQueryValue("another", "value3")
			},
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"name": []string{"value1", "value2"}, "another": []string{"value3"}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name:      "encodes multipart form data",
			operation: postOperation,
			modifyRequest: func() {
				request.MultipartFileBody("file", strings.NewReader("file-contents"))
			},
			responseBody: "{}",
			muxVerify: func(t *testing.T, r *http.Request) {

				ct := r.Header.Get("content-type")
				assert.True(t, strings.HasPrefix(ct, "multipart/form-data"))

				err := r.ParseMultipartForm(2 << 10)
				assert.NoError(t, err)

				file, header, err := r.FormFile("file")
				if assert.NoError(t, err) {
					assert.Equal(t, "image", header.Filename)

					bytes, err := ioutil.ReadAll(file)
					assert.NoError(t, err)

					assert.Equal(t, "file-contents", string(bytes))
				}
			},
		}, {
			name:         "single error",
			operation:    getOperation,
			responseBody: "{\"errors\":[{\"message\":\"testerr\"}]}",
			responseCode: http.StatusNotAcceptable,
			expectErr:    "testerr",
			verify: func(t *testing.T) {
				assert.Equal(t, 1, len(errResult.Errors))
			},
		}, {
			name:         "multiple errors",
			operation:    getOperation,
			responseBody: "{\"errors\":[{\"message\":\"testerr\"},{\"message\":\"another\"}]}",
			responseCode: http.StatusNotAcceptable,
			expectErr:    "testerr",
			verify: func(t *testing.T) {
				assert.Equal(t, 2, len(errResult.Errors))
				assert.Equal(t, "another", errResult.Errors[1].Message)
			},
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {

			mux, riaas, teardown := test.SetupServer(t)
			defer teardown()

			if testcase.responseCode == 0 {
				testcase.responseCode = http.StatusOK
			}

			test.SetupMuxResponse(t, mux, "/resource", testcase.operation.Method, testcase.requestBody, testcase.responseCode, testcase.responseBody, testcase.muxVerify)

			request = riaas.NewRequest(testcase.operation)

			if testcase.modifyRequest != nil {
				testcase.modifyRequest()
			}

			request.JSONError(&errResult)

			resp, err := request.Invoke()

			assert.Equal(t, testcase.responseCode, resp.StatusCode)

			if testcase.expectErr != "" && assert.Error(t, err) {
				assert.Equal(t, testcase.expectErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			if testcase.verify != nil {
				testcase.verify(t)
			}
		})
	}
}

func TestDebugMode(t *testing.T) {

	var (
		riaas   client.Client
		request *client.Request
		log     *bytes.Buffer
	)

	testcases := []struct {
		name string

		operation *client.Operation

		setup func()

		verify func(t *testing.T)
	}{
		{
			name:      "records the request method and resource",
			operation: getOperation,
			verify: func(t *testing.T) {
				assert.Contains(t, log.String(), "REQUEST:")
				assert.Contains(t, log.String(), "GET /resource HTTP/1.1")
			},
		}, {
			name:      "records the request body",
			operation: postOperation,
			setup: func() {
				body := map[string]string{"name": "value"}
				request = request.JSONBody(body)
			},
			verify: func(t *testing.T) {
				assert.Contains(t, log.String(), "\n"+`{"name":"value"}`+"\n")
			},
		}, {
			name:      "records the response code",
			operation: getOperation,
			verify: func(t *testing.T) {
				assert.Contains(t, log.String(), "RESPONSE:")
				assert.Contains(t, log.String(), "HTTP/1.1 200 OK")
			},
		}, {
			name:      "records the response body",
			operation: getOperation,
			verify: func(t *testing.T) {
				assert.Contains(t, log.String(), "testBody")
			},
		}, {
			name:      "redacts the Authorizration header value",
			operation: getOperation,
			verify: func(t *testing.T) {
				assert.Contains(t, log.String(), "Authorization: [REDACTED]")
			},
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {

			mux := http.NewServeMux()
			s := httptest.NewServer(mux)

			log = &bytes.Buffer{}

			riaas = client.New(s.URL, http.DefaultClient).WithDebug(log).WithAuthToken("auth-token")

			defer s.Close()

			mux.HandleFunc("/resource", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "testBody")
			})

			request = riaas.NewRequest(testcase.operation)

			if testcase.setup != nil {
				testcase.setup()
			}

			_, err := request.Invoke()

			assert.NoError(t, err)

			testcase.verify(t)
		})
	}
}

func TestOperationURLProcessing(t *testing.T) {

	testcases := []struct {
		name        string
		baseURL     string
		operation   *client.Operation
		expectedURL string
	}{
		{
			"absolute path",
			"http://127.0.0.1/v2",
			&client.Operation{PathPattern: "/absolute/path"},
			"http://127.0.0.1/absolute/path",
		}, {
			"relative path base does not end with slash",
			"http://127.0.0.1/v2",
			&client.Operation{PathPattern: "relative/path"},
			"http://127.0.0.1/v2/relative/path",
		}, {
			"relative path when base ends with slash",
			"http://127.0.0.1/v2/",
			&client.Operation{PathPattern: "relative/path"},
			"http://127.0.0.1/v2/relative/path",
		}, {
			"relative path parent",
			"http://127.0.0.1/v2",
			&client.Operation{PathPattern: "../path"},
			"http://127.0.0.1/path",
		}, {
			"relative path with .. beyond root",
			"http://127.0.0.1/v2",
			&client.Operation{PathPattern: "../../../../path"},
			"http://127.0.0.1/path",
		}, {
			"broken base URL",
			"://127.0.0.1/v2",
			&client.Operation{PathPattern: "/path"},
			"",
		}, {
			"broken relative path",
			"http://127.0.0.1/v2",
			&client.Operation{PathPattern: "://example.com"},
			"",
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {
			c := client.New(testcase.baseURL, http.DefaultClient)
			actualURL := c.NewRequest(testcase.operation).URL()
			assert.Equal(t, testcase.expectedURL, actualURL)
		})
	}
}
