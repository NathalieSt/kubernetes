#!/bin/bash

# Script to generate a self-signed certificate for nginx-gateway
# This certificate will be valid for all subdomains under cluster.netbird.selfhosted
# Usage: ./generate-cert.sh [primary-domain]

DOMAIN=${1:-"*.cluster.netbird.selfhosted"}
NAMESPACE="nginx-gateway"

echo "Generating wildcard self-signed certificate for: $DOMAIN"
echo "This certificate will be valid for all *.cluster.netbird.selfhosted subdomains"

# Store original directory
ORIGINAL_DIR=$(pwd)

# Create temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Generate private key
openssl genrsa -out tls.key 2048

# Create certificate configuration
cat > cert.conf <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = US
ST = Some-State
L = City
O = Internet Widgits Pty Ltd
CN = $DOMAIN

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = *.cluster.netbird.selfhosted
DNS.2 = cluster.netbird.selfhosted
DNS.3 = jellyfin.cluster.netbird.selfhosted
EOF

# Generate certificate signing request
openssl req -new -key tls.key -out tls.csr -config cert.conf

# Generate self-signed certificate valid for 1 year
openssl x509 -req -in tls.csr -signkey tls.key -out tls.crt -days 365 -extensions v3_req -extfile cert.conf

echo "Certificate generated successfully!"
echo "Certificate details:"
openssl x509 -in tls.crt -text -noout | grep -A2 "Subject:"
openssl x509 -in tls.crt -text -noout | grep -A3 "Subject Alternative Name"

# Create Kubernetes secret YAML in the original directory
kubectl create secret tls nginx-gateway-tls \
  --cert=tls.crt \
  --key=tls.key \
  --namespace=$NAMESPACE \
  --dry-run=client -o yaml > "$ORIGINAL_DIR/nginx-gateway-tls-secret.yaml"

echo ""
echo "Kubernetes secret YAML created: nginx-gateway-tls-secret.yaml"
echo "You can apply it with: kubectl apply -f nginx-gateway-tls-secret.yaml"

# Show certificate validity
echo ""
echo "Certificate validity:"
openssl x509 -in tls.crt -dates -noout

# Clean up
cd - > /dev/null
rm -rf "$TEMP_DIR"

echo ""
echo "Certificate generation complete!"
echo "Remember to update the tls-secret.yaml file with the new certificate data."
