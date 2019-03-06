/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package models

// CPU ...
type CPU struct {
	Architecture string `json:"Architecture,omitempty"`
	Cores        int64  `json:"cores,omitempty"`
	Frequency    int64  `json:"frequency,omitempty"`
}
