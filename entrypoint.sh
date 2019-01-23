#!/bin/sh -e

exec /app/k8s-metadata-injection --clusterName ${clusterName} 2>&1
