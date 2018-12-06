FROM alpine:latest
ADD k8s-env-inject /k8s-env-inject
ENTRYPOINT ["/k8s-env-inject"]