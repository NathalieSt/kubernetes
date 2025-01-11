# Use AdGuard Home as base image
FROM adguard/adguardhome:latest

# Install required packages for Tailscale
RUN apk update && apk add tailscale

# Create startup script
RUN echo '#!/bin/sh\n\
tailscaled --state=/var/lib/tailscale/tailscaled.state --socket=/var/run/tailscale/tailscaled.sock &\n\
sleep 5\n\
tailscale up --accept-dns=false --authkey=$TAILSCALE_AUTHKEY --hostname=adguard-home\n\
/opt/adguardhome/AdGuardHome --host 0.0.0.0' > /start.sh \
    && chmod +x /start.sh

# Expose AdGuard Home ports
EXPOSE 53/tcp 53/udp 80/tcp 443/tcp 853/tcp 3000/tcp

ENTRYPOINT ["/start.sh"]