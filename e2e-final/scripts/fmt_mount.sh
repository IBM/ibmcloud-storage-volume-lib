#!/bin/bash
attachment_id=$1
fstype=$2
attach_dir=${attachment_id:0:20}
echo $attach_dir
dir=$(dir -l /dev/disk/by-id | grep virtio-$attach_dir)
device=$(echo $dir | tr '../../' '\n' | tail -1)
devicepath=/dev/$device
echo "devicepath is $devicepath"
echo "Creating partition"
echo -e "o\nn\np\n1\n\n\nw" | fdisk $devicepath
echo "formatting"
p=1
mkfs.$fstype $devicepath$p
echo "create mount dir"
mkdir /mount-dir
mount -t $fstype $devicepath$p /mount-dir
touch /mount-dir/temp
echo "hello how are you . I am fine, helo world" > /mount-dir/temp
sudo dd if=/dev/urandom of=/mount-dir/test1.bin bs=1G count=1 iflag=fullblock oflag=direct
sudo dd if=/dev/urandom of=/mount-dir/test2.bin bs=1G count=1 iflag=fullblock oflag=direct
sudo dd if=/dev/urandom of=/mount-dir/test3.bin bs=1G count=1 iflag=fullblock oflag=direct
echo "checking for used size"
used_size=$(echo $(df -h | grep /mount-dir) | awk '{print $3}')