#!/bin/sh -e

[ -z ${CERTS_DIR} ] && CERTS_DIR=/etc/webhook/certs

if [ ! -e ${CERTS_DIR}/cert.pem ]; then
    echo "Certificate ${CERTS_DIR}/cert.pem is missing, auto-generating certs..."
    echo

    ./create-certs.sh --generate-only

    # Install the certificates
    mkdir -p ${CERTS_DIR}
    cp server-cert.pem ${CERTS_DIR}/cert.pem
    cp server-key.pem ${CERTS_DIR}/key.pem

    # Create the caBundle variable
    CA_OPTS="-caBundle=$(cat ca-cert.pem | base64 | tr -d '\n')"

    # Cleanup
    rm *.pem
fi

exec /k8s-metadata-injection ${CA_OPTS} -tlsCertFile=${CERTS_DIR}/cert.pem -tlsKeyFile=${CERTS_DIR}/key.pem -alsologtostderr -v=4 2>&1