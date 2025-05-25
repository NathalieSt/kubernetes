#!/bin/bash
# PREREQUISITE: sudoer user
# run this file with sudo

sudo apt update
sudo apt install ufw

sudo ufw allow in ssh

# allow netbird
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 33073/tcp
sudo ufw allow 10000/tcp
sudo ufw allow 33080/tcp
sudo ufw allow 3478/udp
sudo ufw allow 49152:65535/udp

sudo ufw enable