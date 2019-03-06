/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package impl

// SlackMessage represents a report that can be marshalled and sent through a slack webhook
type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment is a struct that can be added to a SlackMessage with additional information
type Attachment struct {
	Text             string   `json:"text"`
	Colour           string   `json:"color"`
	Fields           []Field  `json:"fields"`
	FormattedEntries []string `json:"mrkdwn_in"`
}

// Field is a struct that can be used to add individual items of data to a SlackMessage Attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}
