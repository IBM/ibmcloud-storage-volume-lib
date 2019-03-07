/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Bluemix Container Registry, 5737-D42
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets,  * irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/
package reasoncode

import (
	"fmt"
)

// Wrapper Message/Error Class
type Message struct {
	Code        string
	Description string
	Type        string
	RC          int
	Action      string
}

//Implement the Error() interface method
func (msg Message) Error() string {

	return msg.Info()
}

func (msg Message) Info() string {

	return fmt.Sprintf("{Code:%s, Description:%s, Type:%s, RC:%d}", msg.Code, msg.Description, msg.Type, msg.RC)
}

func GetUserErr(code string, err error, args ...interface{}) error {
	//Incase of no error message, dont construct the Error Object
	if err == nil {
		return nil
	}
	userMsg := GetUserMsg(code, args...)
	userMsg.Description = userMsg.Description + " [Backend Error:" + err.Error() + "]"
	return userMsg
}

func GetUserMsg(code string, args ...interface{}) Message {
	userMsg := messages_en[code]
	if len(args) > 0 {
		userMsg.Description = fmt.Sprintf(userMsg.Description, args...)
	}
	return userMsg
}

func GetUserError(code string, err error, args ...interface{}) error {
	userMsg := GetUserMsg(code, args...)

	if err != nil {
		userMsg.Description = userMsg.Description + " [Backend Error:" + err.Error() + "]"
	}
	return userMsg
}
