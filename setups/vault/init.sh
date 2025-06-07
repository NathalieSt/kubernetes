#!/bin/bash
set -e

echo "Waiting for Vault to be ready..."
until vault status 2>/dev/null; do
    echo "Vault not ready yet, waiting..."
    sleep 5
done

echo "Checking if Vault is already initialized..."
if vault status | grep -q "Initialized.*true"; then
    echo "Vault is already initialized"
    
    # For dev mode, try the dev root token
    echo "Attempting to login with dev root token..."
    if vault auth root 2>/dev/null; then
    echo "Successfully authenticated with dev root token"
    else
    echo "Dev root token failed, trying to find stored token..."
    if [ -f /vault/config/root-token ]; then
        vault auth "$(cat /vault/config/root-token)"
    else
        echo "No authentication method available. This job expects Vault to be in dev mode."
        echo "Please configure Vault with a known root token."
        exit 1
    fi
    fi
else
    echo "Initializing Vault..."
    vault operator init -key-shares=1 -key-threshold=1 > /tmp/vault-init.txt
    
    # Extract keys and token
    UNSEAL_KEY=$(grep 'Unseal Key 1:' /tmp/vault-init.txt | awk '{print $NF}')
    ROOT_TOKEN=$(grep 'Initial Root Token:' /tmp/vault-init.txt | awk '{print $NF}')
    
    echo "Vault initialized. Unseal key and root token:"
    echo "Unseal Key: $UNSEAL_KEY"
    echo "Root Token: $ROOT_TOKEN"
    echo "IMPORTANT: Store these securely!"
    
    # Unseal Vault
    vault operator unseal "$UNSEAL_KEY"
    
    # Login with root token
    vault auth "$ROOT_TOKEN"
fi

echo "Configuring Kubernetes auth backend..."

# Enable kubernetes auth if not already enabled
if ! vault auth list | grep -q kubernetes; then
    vault auth enable kubernetes
fi

# Get Kubernetes API details
K8S_HOST="https://kubernetes.default.svc.cluster.local"
SA_JWT_TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
SA_CA_CRT=$(cat /var/run/secrets/kubernetes.io/serviceaccount/ca.crt)

# Configure kubernetes auth
vault write auth/kubernetes/config \
    token_reviewer_jwt="$SA_JWT_TOKEN" \
    kubernetes_host="$K8S_HOST" \
    kubernetes_ca_cert="$SA_CA_CRT"

# Create a policy for k8s role
vault policy write k8s-policy - <<EOF
path "kv/*" {
    capabilities = ["read", "list"]
}
path "kvv2/*" {
    capabilities = ["read", "list"]
}
EOF

# Create k8s role
vault write auth/kubernetes/role/k8s \
    bound_service_account_names=vault-auth-serviceaccount \
    bound_service_account_namespaces="*" \
    policies=k8s-policy \
    ttl=1h

# Enable KV v2 secrets engine if not already enabled
if ! vault secrets list | grep -q kvv2; then
    vault secrets enable -path=kvv2 kv-v2
fi

# Create a test secret
vault kv put kvv2/netbird username=test password=test123

echo "Vault configuration completed successfully!"