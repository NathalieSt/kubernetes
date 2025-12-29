# Dockerfile
FROM caddy:2.10.2-builder-alpine AS builder

# Add the Hetzner DNS plugin
RUN xcaddy build \
    --with github.com/libdns/hetzner/v2 \
    --with github.com/caddy-dns/hetzner/v2 

# -----------------------------------------------------------------
# Runtime stage – copy the compiled binary into a minimal image
FROM caddy:alpine

# Copy the custom-built Caddy binary from the builder stage
COPY --from=builder /usr/bin/caddy /usr/bin/caddy

# Optional: expose the default HTTP/HTTPS ports
EXPOSE 80 443

# Keep the default entrypoint (Caddy)
ENTRYPOINT ["caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile"]