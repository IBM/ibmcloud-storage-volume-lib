/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package client

// Params ...
type Params map[string]string

// Copy performs a shallow copy of a Params object
func (p Params) Copy() Params {
	params := Params{}
	for k, v := range p {
		params[k] = v
	}
	return params
}
