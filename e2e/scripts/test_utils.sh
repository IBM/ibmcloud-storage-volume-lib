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
      if [[ $attempts -gt 60 ]]; then
         echo "$instance_id is not ready."
         ibmcloud is in $instance_id
         exit 1
      fi
      echo "$instance_id state == $instance_status Sleeping 60 seconds"
      sleep 60
  done
}

function cleanup {
  attempts=0
  echo "Detaching instance attachments"
  ibmcloud is instances |grep "e2e-common-lib" | awk -F ' ' '{print $1}' |
  while IFS= read -r instanceID
  do
    ibmcloud is in-vols $instanceID | grep -v "boot" | grep "data" | grep "e2e-" | awk -F ' ' '{print $1}' |
    while IFS= read -r attachmentID
    do
      echo "Detaching $attachmentID"
      ibmcloud is instance-volume-attachment-detach $instanceID $attachmentID -f
      while true; do
        attempts=$((attempts+1))
        attachmentStatus=$(ibmcloud is in-vol $instanceID $attachmentID | grep Status | awk -F ' ' '{print $2}')
        if [ -z "$attachmentStatus" ]; then
           echo "$attachmentStatus is detached."
           break
        fi
        if [ $attempts -gt 60 ]; then
           echo "$attachmentStatus is still in attached state."
           ibmcloud is in-vol $instanceID $attachmentID
           exit 1
        fi
        echo "$attachmentID state == $attachmentStatus Sleeping 60 seconds"
        sleep 30
      done
    done
    ibmcloud is instance-delete $instanceID -f
  done

  attempts=0
  echo "Deleting volumes"
  ibmcloud is vols | grep -v "boot" | grep "e2e-" | awk -F ' ' '{print $1}' |
  while IFS= read -r volumeID
  do
    echo "Deleting $volumeID"
    ibmcloud is vold $volumeID -f
    while true; do
      attempts=$((attempts+1))
      volumeStatus=$(ibmcloud is vol $volumeID | grep Status | awk -F ' ' '{print $2}')
      if [ -z "$volumeStatus" ]; then
         echo "$volumeStatus is deleted."
         break
      fi
      if [ $attempts -gt 60 ]; then
         echo "$volumeStatus is still available."
         ibmcloud is vol $volumeID
         exit 1
      fi
      echo "$volumeID state == $volumeStatus Sleeping 60 seconds"
      sleep 30
    done
  done
}
