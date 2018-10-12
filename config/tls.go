/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package config

import (
	"crypto/tls"
	"net/http"
	"time"
)

// GeneralCAHttpClient returns an http.Client configured for general use
func GeneralCAHttpClient() (*http.Client, error) {

	httpClient := &http.Client{

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12, // Require TLS 1.2 or higher
			},
		},

		// softlayer.go has been overriding http.DefaultClient and forcing 120s
		// timeout on us, so we'll continue to force it on ourselves in case
		// we've accidentally become acustomed to it.
		Timeout: time.Second * 120,
	}

	return httpClient, nil
}

// GeneralCAHttpClientWithTimeout returns an http.Client configured for general use
func GeneralCAHttpClientWithTimeout(timeout time.Duration) (*http.Client, error) {

	httpClient := &http.Client{

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12, // Require TLS 1.2 or higher
			},
		},

		Timeout: timeout,
	}

	return httpClient, nil
}
