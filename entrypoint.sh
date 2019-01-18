#!/bin/sh -e

[ -z ${CERTS_DIR} ] && CERTS_DIR=/etc/tls-key-cert-pair

exec /k8s-metadata-injection ${CA_OPTS} -tlsCertFile=${CERTS_DIR}/tls.crt -tlsKeyFile=${CERTS_DIR}/tls.key -alsologtostderr -v=4 2>&1
