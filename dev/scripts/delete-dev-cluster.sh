#!/bin/bash
set -e

CLUSTER_NAME="autorollout-dev"

echo "Starting Autorollout Dev Cluster Cleanup..."

if ! kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    echo "WARNING: Cluster $CLUSTER_NAME does not exist. Nothing to clean up."
else
    echo "Deleting Kind cluster: $CLUSTER_NAME"
    kind delete cluster --name=$CLUSTER_NAME
    echo "Cluster deleted successfully"
fi

echo ""
echo "Autorollout Dev Cluster Cleanup Complete!"
echo ""
