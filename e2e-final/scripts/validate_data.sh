#!/bin/bash
attachment_id=$1
fstype=$2
attach_dir=${attachment_id:0:20}
echo $attach_dir
dir=$(dir -l /dev/disk/by-id | grep virtio-$attach_dir)
device=$(echo $dir | tr '../../' '\n' | tail -1)
devicepath=/dev/$device
echo "devicepath is $devicepath"
echo "create mount dir"
mkdir /restore-dir
mount -t $fstype $devicepath$p /restore-dir
actual_size=$(echo $(df -h | grep /mount-dir) | awk '{print $3}')
echo "checking for used size"
START_TIME=$(date +%s)
until [ $(echo $(df -h | grep /restore-dir) | awk '{print $3}') != actual_size ]
do 
    echo "$(echo $(df -h | grep /restore-dir) | awk '{print $3}') of data restored"
done
END_TIME=$(date +%s)
echo "It took $(($END_TIME - $START_TIME)) seconds to restore the whole data"