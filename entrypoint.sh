#!/bin/sh -e

export NEW_RELIC_K8S_METADATA_INJECTION_CLUSTER_NAME=${clusterName}

exec /app/k8s-metadata-injection 2>&1
