#!/bin/bash

sudo apt-get update && sudo apt-get install fail2ban

sudo systemctl start fail2ban

sudo systemctl enable fail2ban

# PERMISISION DENIED: just make it manually
#sudo touch /etc/fail2ban/jail.d/fort.conf
#
#sudo cat > /etc/fail2ban/jail.d/fort.conf <<EOF
#
#[DEFAULT]
#
#bantime = 5d
#
#findtime = 2d
#
#ignoreip = 127.0.0.1/8 192.168.0.0/16
#
#maxretry = 2
#
#banaction = ufw
#banaction_allports = ufw
#
#EOF

sudo systemctl restart fail2ban