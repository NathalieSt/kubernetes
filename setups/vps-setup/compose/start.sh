#!/bin/bash

# export env variables
export $(cat ./crowdsec.env | xargs)

export $(cat ./crowdsec-ui/netbird-agent.env | xargs)

export $(cat ./keycloak/keycloak.env | xargs)

export $(cat ./netbird/dashboard.env | xargs)
export $(cat ./netbird/relay.env | xargs)

export $(cat ./zdb.env | xargs)

export $(cat ./zitadel.env | xargs)

sudo docker compose up -d

# wait for some time to ensure everything has started
sleep(500)

# add machine to crowdsec so that the UI can access it
sudo docker compose exec crowdsec cscli machines add $CROWDSEC_USER --password $CROWDSEC_PASSWORD -f /dev/null

# install and enable firewall bouncer 
sudo apt install crowdsec-firewall-bouncer
sudo systemctl enable crowdsec-firewall-bouncer.service

sudo docker compose exec cscli bouncers add host-firewall-bouncer
