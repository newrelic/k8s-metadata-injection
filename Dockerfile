FROM alpine:latest
RUN apk add --update openssl
COPY entrypoint.sh create-certs.sh k8s-env-inject /
CMD ["/entrypoint.sh"]