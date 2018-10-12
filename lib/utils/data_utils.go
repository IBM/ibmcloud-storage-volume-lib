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

// SafeStringValue returns the referenced string value, treating nil as equivalent to "".
// It is intended as a type-safe and nil-safe test for empty values in data fields of
func SafeStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
