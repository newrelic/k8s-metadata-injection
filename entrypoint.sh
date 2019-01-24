#!/bin/sh -e

exec NR_K8S_METADATA_INJECTION_CLUSTER_NAME=${clusterName} /app/k8s-metadata-injection 2>&1
