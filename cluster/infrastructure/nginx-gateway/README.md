# Nginx Gateway TLS Configuration

This directory contains the configuration for the nginx gateway with TLS support.

## TLS Certificate Setup

The current configuration uses a **wildcard self-signed certificate** for `*.cluster.netbird.selfhosted`, which means it's valid for:
- `jellyfin.cluster.netbird.selfhosted`
- `grafana.cluster.netbird.selfhosted` 
- `prometheus.cluster.netbird.selfhosted`
- Any other subdomain you want to add under `cluster.netbird.selfhosted`

### To replace with a proper certificate:

1. **Let's Encrypt Certificate** (recommended for production):
   ```bash
   # Install cert-manager in your cluster first
   kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.2/cert-manager.yaml
   
   # Create a ClusterIssuer for Let's Encrypt
   kubectl apply -f - <<EOF
   apiVersion: cert-manager.io/v1
   kind: ClusterIssuer
   metadata:
     name: letsencrypt-prod
   spec:
     acme:
       server: https://acme-v02.api.letsencrypt.org/directory
       email: your-email@example.com
       privateKeySecretRef:
         name: letsencrypt-prod
       solvers:
       - http01:
           ingress:
             class: nginx
   EOF
   
   # Update the tls-secret.yaml to use cert-manager annotation:
   # Add this annotation to the secret metadata:
   # cert-manager.io/cluster-issuer: "letsencrypt-prod"
   ```

2. **Custom Certificate Authority**:
   Replace the `tls.crt` and `tls.key` data in `tls-secret.yaml` with your CA-signed certificate.

3. **Generate a new self-signed certificate**:
   ```bash
   # Generate private key
   openssl genrsa -out tls.key 2048
   
   # Generate certificate signing request
   openssl req -new -key tls.key -out tls.csr -subj "/CN=jellyfin.cluster.netbird.selfhosted"
   
   # Generate self-signed certificate
   openssl x509 -req -in tls.csr -signkey tls.key -out tls.crt -days 365
   
   # Update the secret
   kubectl create secret tls nginx-gateway-tls \
     --cert=tls.crt \
     --key=tls.key \
     --namespace=nginx-gateway \
     --dry-run=client -o yaml > tls-secret.yaml
   ```

## Security Features

The nginx configuration includes:
- HTTP to HTTPS redirect
- Modern TLS protocols (TLSv1.2 and TLSv1.3)
- Strong cipher suites
- Security headers (HSTS, X-Frame-Options, etc.)
- WebSocket support for Jellyfin
- Optimized proxy buffers

## Network Access

The nginx gateway uses a netbird sidecar for external network access, eliminating the need for a LoadBalancer service. The pod is accessible directly through the netbird mesh network at `jellyfin.cluster.netbird.selfhosted`.

## Monitoring

You can check the TLS configuration with:
```bash
# Check certificate details
openssl s_client -connect jellyfin.cluster.netbird.selfhosted:443 -servername jellyfin.cluster.netbird.selfhosted

# Test SSL rating
curl -s "https://www.ssllabs.com/ssltest/analyze.html?d=jellyfin.cluster.netbird.selfhosted&hideResults=on"
```
