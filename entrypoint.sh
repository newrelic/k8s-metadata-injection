#!/bin/sh -e

export NEW_RELIC_K8S_METADATA_INJECTION_CLUSTER_NAME=${clusterName}
# TODO remove
stat /app/k8s-metadata-injection
ldd /app/k8s-metadata-injection || true
exec /app/k8s-metadata-injection 2>&1
