#!/bin/bash

# Script to verify cert-manager and certificate status
# Usage: ./verify-cert-manager.sh

set -e

echo "🔍 Checking cert-manager deployment status..."

# Check if cert-manager namespace exists
if kubectl get namespace cert-manager >/dev/null 2>&1; then
    echo "✅ cert-manager namespace exists"
else
    echo "❌ cert-manager namespace not found"
    exit 1
fi

# Check cert-manager pods
echo ""
echo "📋 cert-manager pods status:"
kubectl get pods -n cert-manager

# Check if ClusterIssuer is ready
echo ""
echo "🔑 ClusterIssuer status:"
kubectl get clusterissuer selfsigned-cluster-issuer -o wide

# Check certificate status
echo ""
echo "📜 Certificate status:"
if kubectl get certificate nginx-gateway-tls -n nginx-gateway >/dev/null 2>&1; then
    kubectl get certificate nginx-gateway-tls -n nginx-gateway -o wide
    echo ""
    echo "📝 Certificate details:"
    kubectl describe certificate nginx-gateway-tls -n nginx-gateway | grep -A 10 "Status:"
else
    echo "❌ nginx-gateway-tls certificate not found"
fi

# Check if secret exists and contains certificate data
echo ""
echo "🔐 TLS Secret status:"
if kubectl get secret nginx-gateway-tls -n nginx-gateway >/dev/null 2>&1; then
    echo "✅ nginx-gateway-tls secret exists"
    
    # Check certificate validity
    echo ""
    echo "📅 Certificate validity:"
    kubectl get secret nginx-gateway-tls -n nginx-gateway -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -noout -dates
    
    echo ""
    echo "🌐 Certificate domains:"
    kubectl get secret nginx-gateway-tls -n nginx-gateway -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -noout -text | grep -A1 "Subject Alternative Name"
else
    echo "❌ nginx-gateway-tls secret not found"
fi

echo ""
echo "🚀 All checks completed!"
