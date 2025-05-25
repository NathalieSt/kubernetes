#!/bin/bash
# PREREQUISITE: sudoer user
# run this file with sudo

apt update
apt install ufw

ufw allow in ssh

# allow netbird
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 33073/tcp
ufw allow 10000/tcp
ufw allow 33080/tcp
ufw allow 3478/udp
ufw allow 49152:65535/udp

ufw enable