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

# Copy and setup startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh && dos2unix /start.sh

# Expose AdGuard Home ports
EXPOSE 53/tcp 53/udp 80/tcp 443/tcp 853/tcp 3000/tcp

ENTRYPOINT ["/start.sh"]