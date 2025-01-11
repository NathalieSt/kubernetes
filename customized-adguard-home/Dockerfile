FROM alpine:latest

# Install required packages
RUN apk update && \
    apk add --no-cache \
    curl \
    tar \
    tailscale

# Install AdGuard Home
RUN curl -L -o /tmp/AdGuardHome.tar.gz https://github.com/AdguardTeam/AdGuardHome/releases/latest/download/AdGuardHome_linux_amd64.tar.gz && \
    tar xvf /tmp/AdGuardHome.tar.gz -C /opt && \
    rm /tmp/AdGuardHome.tar.gz && \
    chmod +x /opt/AdGuardHome/AdGuardHome

# Create startup script
RUN echo '#!/bin/sh\n\
tailscaled --state=/var/lib/tailscale/tailscaled.state --socket=/var/run/tailscale/tailscaled.sock &\n\
sleep 5\n\
tailscale up --accept-dns=false --authkey=$TAILSCALE_AUTHKEY --hostname=adguard-home\n\
/opt/AdGuardHome/AdGuardHome --host 0.0.0.0' > /start.sh && \
    chmod +x /start.sh

# Expose AdGuard Home ports
EXPOSE 53/tcp 53/udp 80/tcp 443/tcp 853/tcp 3000/tcp

ENTRYPOINT ["/start.sh"]