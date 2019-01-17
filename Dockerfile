FROM alpine:latest
RUN apk add --update openssl
COPY entrypoint.sh bin/k8s-metadata-injection /
CMD ["/entrypoint.sh"]