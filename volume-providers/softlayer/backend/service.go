/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package backend

// Service is the parent of all services, providing common methods such as Mask,
// Id, Filter etc
type Service interface {
	Filter(string) Service
	ID(int) Service
	Mask(string) Service
}
