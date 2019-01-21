#!/bin/sh -e

[ -z ${CERTS_DIR} ] && CERTS_DIR=/etc/tls-key-cert-pair

exec ./k8s-metadata-injection ${CA_OPTS} -tlsCertFile=${CERTS_DIR}/cert.pem -tlsKeyFile=${CERTS_DIR}/key.pem 2>&1
