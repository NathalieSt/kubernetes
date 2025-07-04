#!/bin/bash

# create mount directory
mkdir /mnt/external_ssd

# mount external ssd
echo "/dev/sda /mnt/external_ssd    ext4    defaults        0       0" | tee -a /etc/fstab

# test mount
mount -a

# install nfs server
dietpi-software install 109

# export mounted ssd via nfs
echo "/mnt/external_ssd *(rw,async,no_root_squash,crossmnt,no_subtree_check)" | tee -a /etc/exports.d/dietpi.exports