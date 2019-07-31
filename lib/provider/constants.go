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

//RequestIDType ... Context won't allow keys as "string" type
type RequestIDType string

const (
	// RequestID ...
	RequestID RequestIDType = RequestIDType("request-id")
)
