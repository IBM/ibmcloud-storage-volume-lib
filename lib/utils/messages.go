/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Bluemix Container Registry, 5737-D42
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets,  * irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package util

import (
	"fmt"
)

// Message Wrapper Message/Error Class
type Message struct {
	Code         string
	Type         string
	RequestID    string
	Description  string
	BackendError string
	RC           int
	Action       string
}

// Error Implement the Error() interface method
func (msg Message) Error() string {
	return msg.Info()
}

// Info ...
func (msg Message) Info() string {
	return fmt.Sprintf("{Code:%s, Type:%s, Description:%s, BackendError:%s, RC:%d}", msg.Code, msg.Type, msg.Description, msg.BackendError, msg.RC)
}
