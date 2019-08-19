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


function ibmcloud_login {
    echo 'Logging Into IbmCloud Container Service'
    ibmcloud --version
    ibmcloud plugin list
    ibmcloud login  -r $TEST_REGION -a $IC_API_ENDPOINT -u $IC_USERNAME -p $IC_LOGIN_PASSWORD -c $IC_ACCOUNT -o $IC_ORG -s $IC_SPACE
    ibmcloud ks init --host $IC_HOST_EP
}

function check_instance_state {
   attempts=0
   instance_id=$1
   while true; do
      attempts=$((attempts+1))
      instance_status=$(ibmcloud is in $instance_id |grep Status | tr -s " " | awk -F ' ' '{print $2}')
      if [   "$instance_status" = "running" ]; then
         echo "$instance_id is ready."
         break
      fi
      if [[ $attempts -gt 30 ]]; then
         echo "$instance_id is not ready."
         ibmcloud is in $instance_id
         exit 1
      fi
      echo "$instance_id state == $instance_status Sleeping 10 seconds"
      sleep 10
  done
}
