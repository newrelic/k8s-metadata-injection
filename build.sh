#!/bin/bash -e

REGISTRY=docker.coscale.com:5000
IMAGE=fryckbosch/newrelic-k8s-metadata-injector:v1

dep ensure
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-env-inject .
docker build --no-cache -t ${REGISTRY}/${IMAGE} .
rm -rf k8s-env-inject

docker push ${REGISTRY}/${IMAGE}
