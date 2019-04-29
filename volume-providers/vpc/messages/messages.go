/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Bluemix Container Registry, 5737-D42
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets,  * irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package messages

import (
	"fmt"
	util "github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
)

// MessagesEn ...
var MessagesEn map[string]util.Message

// GetUserErr ...
func GetUserErr(code string, err error, args ...interface{}) error {
	//Incase of no error message, dont construct the Error Object
	if err == nil {
		return nil
	}
	userMsg := GetUserMsg(code, args...)
	userMsg.BackendError = err.Error()
	return userMsg
}

// GetUserMsg ...
func GetUserMsg(code string, args ...interface{}) util.Message {
	userMsg := MessagesEn[code]
	if len(args) > 0 {
		userMsg.Description = fmt.Sprintf(userMsg.Description, args...)
	}
	return userMsg
}

// GetUserError ...
func GetUserError(code string, err error, args ...interface{}) error {
	userMsg := GetUserMsg(code, args...)

	if err != nil {
		userMsg.BackendError = err.Error()
	}
	return userMsg
}
