#!/bin/sh
tailscaled --state=/var/lib/tailscale/tailscaled.state --socket=/var/run/tailscale/tailscaled.sock &
sleep 5
tailscale up --accept-dns=false --authkey=$TAILSCALE_AUTHKEY --hostname=adguard-home
/opt/AdGuardHome/AdGuardHome --host 0.0.0.0