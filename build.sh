#!/bin/bash -e

IMAGE=quay.io/newrelic/k8s-metadata-injector-dev:latest

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-metadata-injection .
docker build --no-cache -t ${IMAGE} .
rm -rf k8s-metadata-injection