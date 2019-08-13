#!/bin/bash
# Licensed Materials - Property of IBM
#
# (C) Copyright IBM Corp. 2017 All Rights Reserved
#
# US Government Users Restricted Rights - Use, duplicate or
# disclosure restricted by GSA ADP Schedule Contract with
# IBM Corp.
# encoding: utf-8

# Display essage to console and log file
function display {
    echo "$1"
    echo "$1" >>${LOG_FILE}
}

# Sending messages to slack channel
function send_messages_to_slack {
    DATA="\`\`\`"
    while IFS= read LINE; do DATA+="\n$LINE"; done < $LOG_FILE
    DATA+="\`\`\`"
    MESSAGES=$(echo -e $DATA)
    send_message_to_slack "$MESSAGES"
}

function send_messages_to_slack2 {
    while IFS= read LINE; do send_message_to_slack "$LINE"; done < $LOG_FILE
}

