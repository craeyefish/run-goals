#!/bin/bash
# Script to apply all secrets to DOKS cluster
# Run this script after connecting to your DOKS cluster with:
# doctl kubernetes cluster kubeconfig save <cluster-name>

set -e

echo "======================================"
echo "Applying Secrets to DOKS Cluster"
echo "======================================"
echo ""

# Check if kubectl is connected to a cluster
if ! kubectl cluster-info &>/dev/null; then
    echo "❌ Error: kubectl is not connected to a cluster"
    echo "Please connect to your DOKS cluster first:"
    echo "  doctl kubernetes cluster kubeconfig save <cluster-name>"
    exit 1
fi

echo "Current cluster context:"
kubectl config current-context
echo ""

read -p "Is this the correct cluster? (yes/no): " confirm
if [[ ! "$confirm" =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "Aborted. Please switch to the correct cluster context."
    exit 0
fi

echo ""
echo "Applying secrets..."
echo ""

# Apply database credentials
echo "✅ Applying postgres-credentials..."
kubectl apply -f db-secret.yaml

# Apply Strava credentials
echo "✅ Applying strava-credentials..."
kubectl apply -f strava-secret.yaml

# Apply JWT secret
echo "✅ Applying jwt-secret..."
kubectl apply -f jwt-secret.yaml

# Apply Cloudflare tunnel credentials
if [ -f "cloudflared-credentials.yaml" ]; then
    echo "✅ Applying cloudflared-credentials..."
    kubectl apply -f cloudflared-credentials.yaml
else
    echo "⚠️  WARNING: cloudflared-credentials.yaml not found"
    echo "   You need to create this file from Cloudflare dashboard"
    echo "   See: cloudflared-credentials-template.yaml"
fi

echo ""
echo "======================================"
echo "✅ Secrets Applied Successfully!"
echo "======================================"
echo ""
echo "Verify secrets were created:"
echo "  kubectl get secrets -n default"
echo ""
echo "Next steps:"
echo "  1. If you haven't created cloudflared-credentials.yaml, do that now"
echo "  2. Deploy your applications with Flux CD"
echo "  3. Monitor pod status: kubectl get pods -n default -w"
