#!/bin/bash
netbird up --management-url https://netbird.nathalie-stiefsohn.eu --setup-key <insert setup key here> --extra-dns-labels "webserver"

apt update
apt install golang-go

apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/xcaddy/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-xcaddy-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/xcaddy/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-xcaddy.list
apt update
apt install xcaddy


xcaddy build \
    --with github.com/libdns/hetzner/v2 \
    --with github.com/caddy-dns/hetzner/v2 

./caddy run --config=caddy.json