/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package payload

import (
	"encoding/json"
	"io"
)

// JSONConsumer ...
type JSONConsumer struct{ receiver interface{} }

// NewJSONConsumer ...
func NewJSONConsumer(receiver interface{}) *JSONConsumer {
	return &JSONConsumer{receiver: receiver}
}

// Consume ...
func (c *JSONConsumer) Consume(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(c.receiver)
}

// Receiver ...
func (c *JSONConsumer) Receiver() interface{} {
	return c.receiver
}
