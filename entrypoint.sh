#!/bin/sh -e

[ -z ${CERTS_DIR} ] && CERTS_DIR=/etc/webhook/certs

exec //k8s-metadata-injection ${CA_OPTS} -tlsCertFile=${CERTS_DIR}/tls.crt -tlsKeyFile=${CERTS_DIR}/tls.key 2>&1
