#!/bin/sh -e

exec NEW_RELIC_K8S_METADATA_INJECTION_CLUSTER_NAME=${clusterName} /app/k8s-metadata-injection 2>&1
