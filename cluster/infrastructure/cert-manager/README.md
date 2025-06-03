# cert-manager Infrastructure

This directory contains the cert-manager deployment for the Kubernetes cluster using GitOps with Flux.

## Components

- **namespace.yaml**: Creates the cert-manager namespace
- **repository.yaml**: Jetstack Helm repository for cert-manager charts
- **release.yaml**: Helm release configuration for cert-manager
- **cluster-issuer.yaml**: Self-signed ClusterIssuer for certificate generation
- **nginx-gateway-certificate.yaml**: Certificate for nginx-gateway with wildcard support

## Features

- Self-signed certificate issuer for internal use
- Automatic certificate generation and renewal
- Wildcard certificate for `*.cluster.netbird.selfhosted`
- Integration with nginx-gateway for TLS termination

## Certificate Details

The nginx-gateway certificate includes the following domains:
- `*.cluster.netbird.selfhosted` (wildcard)
- `cluster.netbird.selfhosted`
- `jellyfin.cluster.netbird.selfhosted`
- `uptime.cluster.netbird.selfhosted`

## Deployment

The cert-manager is deployed automatically via Flux when changes are pushed to the repository. The nginx-gateway deployment depends on cert-manager being ready first.

## Manual Operations

If you need to manually check the certificate status:

```bash
# Check certificate status
kubectl get certificate -n nginx-gateway

# Check certificate details
kubectl describe certificate nginx-gateway-tls -n nginx-gateway

# Check the generated secret
kubectl get secret nginx-gateway-tls -n nginx-gateway
```

## Migration from Manual Certificates

This setup replaces the manual certificate generation script (`generate-cert.sh`) with automated certificate management. The certificate will be automatically renewed before expiration.
