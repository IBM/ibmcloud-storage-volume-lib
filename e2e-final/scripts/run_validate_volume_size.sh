#!/bin/bash
export instance_ip=$1
export attachement_id=$2
export fstype=$3
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i ssh-key root@$instance_ip "bash -s" < ../scripts/validate_data.sh $attachement_id $fstype