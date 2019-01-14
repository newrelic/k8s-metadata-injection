#!/bin/bash -e

REGISTRY=docker.coscale.com:5000
IMAGE=quay.io/newrelic/k8s-metadata-injector-dev:latest

dep ensure
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-metadata-injection .
docker build --no-cache -t ${REGISTRY}/${IMAGE} .
rm -rf k8s-metadata-injection