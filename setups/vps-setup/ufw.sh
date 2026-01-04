#!/bin/bash
# PREREQUISITE: sudoer user
# run this file with sudo

sudo apt update
sudo apt install ufw

sudo ufw allow in ssh

# just to make absolutely sure no one can access crowdsec api or ui
sudo ufw deny 8080/tcp
sudo ufw deny 3000/tcp

# allow netbird
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 3478/udp