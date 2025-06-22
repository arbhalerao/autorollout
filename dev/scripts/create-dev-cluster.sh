#!/bin/bash
set -e

CLUSTER_NAME="autorollout-dev"

echo "Starting Autorollout Dev Cluster Setup..."

if ! command -v kind &> /dev/null; then
    echo "ERROR: Kind is not installed. Please install Kind first."
    exit 1
fi

if ! command -v kubectl &> /dev/null; then
    echo "ERROR: kubectl is not installed. Please install kubectl first."
    exit 1
fi

echo "Prerequisites check passed"

if kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    echo "WARNING: Cluster $CLUSTER_NAME already exists. Deleting and recreating..."
    kind delete cluster --name=$CLUSTER_NAME
fi

echo "Creating Kind cluster: $CLUSTER_NAME"
kind create cluster --config=dev/clusterconfig/kind-config.yaml --name=$CLUSTER_NAME

echo "Setting kubectl context"
kubectl config use-context kind-$CLUSTER_NAME
kubectl cluster-info --context kind-$CLUSTER_NAME

echo ""
echo "Autorollout Local Dev Infra Setup Complete!"
echo ""
