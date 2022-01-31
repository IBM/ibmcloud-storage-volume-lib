#!/bin/bash
export instance_ip=$1
export attachement_id=$2
export fstype=$3
echo "ssh into vm"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i ssh-key root@$instance_ip "bash -s" <  ../scripts/fmt_mount.sh $attachement_id $fstype