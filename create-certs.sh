#!/bin/sh -e

[ -z ${SERVICE} ] && SERVICE=newrelic-metadata-injection-svc
[ -z ${NAMESPACE} ] && NAMESPACE=default


cat << EOF > openssl.ext
[root]
keyUsage=critical, keyCertSign, digitalSignature, keyEncipherment
basicConstraints=critical, CA:TRUE

[server]
keyUsage=critical, digitalSignature, keyEncipherment
extendedKeyUsage=serverAuth
basicConstraints=critical, CA:FALSE
subjectAltName=DNS:${SERVICE}, DNS:${SERVICE}.${NAMESPACE}, DNS:${SERVICE}.${NAMESPACE}.svc
EOF

# Generate CA-cert
openssl genrsa -out ca-key.pem 2048
openssl req -new -key ca-key.pem -out ca-csr.pem -subj /CN=kubernetes
openssl x509 -req -sha256 -in ca-csr.pem -signkey ca-key.pem -out ca-cert.pem -days 3650 -extfile openssl.ext -extensions root

# Generate Server cert
openssl genrsa -out server-key.pem 2048
openssl req -new -key server-key.pem -out server-csr.pem -subj /CN=${SERVICE}.${NAMESPACE}.svc
openssl x509 -req -sha256 -in server-csr.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -days 3650 -extfile openssl.ext -extensions server

# Cleanup
rm *-csr.pem *.srl openssl.ext

if [ "$1" = "--generate-only" ]; then
    exit 0
fi

# Push the secret to Kubernetes
[ -z ${SECRET} ] && SECRET=newrelic-metadata-injection-certs

kubectl create secret generic ${SECRET} \
    --from-file=key.pem=server-key.pem \
    --from-file=cert.pem=server-cert.pem \
    --dry-run -o yaml |
kubectl -n ${NAMESPACE} apply -f -

# Put the caBundle into the newrelic-metadata-injection-manual-certs.yaml file
export CA_BUNDLE=$(cat ca-cert.pem | base64 | tr -d '\n')

if command -v envsubst >/dev/null 2>&1; then
    cat newrelic-metadata-injection-manual-nocert.yaml | envsubst > newrelic-metadata-injection-manual-cert.yaml
else
    cat newrelic-metadata-injection-manual-nocert.yaml | sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g" > newrelic-metadata-injection-manual-cert.yaml
fi
