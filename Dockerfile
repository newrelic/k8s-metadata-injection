FROM alpine:latest
RUN apk add --update openssl
COPY entrypoint.sh create-certs.sh bin/k8s-metadata-injection /
CMD ["/entrypoint.sh"]