#!/usr/bin/env sh
set -e

printf 'bootstrapping starts:\n'
# shellcheck disable=SC1090
. "$(dirname "$0")/k8s-e2e-bootstraping.sh"
printf 'bootstrapping complete\n'

HELM_RELEASE_NAME="nri-metadata-injection"
WEBHOOK_LABEL="app.kubernetes.io/name=nri-metadata-injection,app.kubernetes.io/instance=${HELM_RELEASE_NAME}"
DUMMY_DEPLOYMENT_NAME="dummy-deployment"
DUMMY_POD_LABEL="app=${DUMMY_DEPLOYMENT_NAME}"
ENV_VARS_PREFIX="NEW_RELIC_METADATA_KUBERNETES"
NAMESPACE_NAME="$(kubectl config view --minify --output 'jsonpath={..namespace}')"
IMAGE_TAG=${IMAGE_TAG:-e2e-test}

IMAGE_TAG="$(curl --silent "https://api.github.com/repos/newrelic/k8s-metadata-injection/releases" | jq -r 'map(select(.prerelease)) | first | .tag_name')"


finish() {
    printf "webhook logs:\n"
    kubectl logs "$(get_pod_name_by_label "$WEBHOOK_LABEL")" || true

    helm uninstall "$HELM_RELEASE_NAME" || true
    kubectl delete deployment ${DUMMY_DEPLOYMENT_NAME} || true
}

# ensure that we build docker image in minikube
[ "$E2E_MINIKUBE_DRIVER" = "none" ] || eval "$(minikube docker-env --shell bash)"

# build webhook docker image

# Set GOOS and GOARCH explicitly since Dockerfile expects them in the binary name
GOOS="linux" GOARCH="amd64" DOCKER_IMAGE_TAG="e2e-test" make -C .. compile build-container

trap finish EXIT

# install the metadata-injection webhook
helm repo add newrelic https://helm-charts.newrelic.com
helm dependency build ../charts/nri-metadata-injection
if ! helm upgrade --install "$HELM_RELEASE_NAME" ../charts/nri-metadata-injection \
                --wait \
                --set cluster=YOUR-CLUSTER-NAME \
                --set image.pullPolicy=Never \
                --set image.tag="$IMAGE_TAG"
then
    printf "Helm failed to install this release\n"
    exit 1
fi


### Testing

# deploy a pod
kubectl create deployment "$DUMMY_DEPLOYMENT_NAME" --image=nginx:latest --dry-run=client -o yaml | kubectl apply -f-

pod_name="$(get_pod_name_by_label "$DUMMY_POD_LABEL")"
if [ "$pod_name" = "" ]; then
    printf "not found any pod with label %s\n" "$DUMMY_POD_LABEL"
    kubectl describe deployment "$DUMMY_DEPLOYMENT_NAME"
    exit 1
fi
wait_for_pod "$pod_name"

kubectl get pods
kubectl describe pod "${pod_name}"

printf "getting env vars for %s\n" "${pod_name}"
set +e # This grep can be empty in the webhook is not correctly running and we want logs and a proper error
date
env_vars="$(kubectl exec "${pod_name}" -- env | grep "${ENV_VARS_PREFIX}")"
set -e
printf "\nInjected environment variables:\n"
printf "%s\n" "$env_vars"

errors=""
for PAIR in \
           "CLUSTER_NAME            YOUR-CLUSTER-NAME" \
           "NODE_NAME               minikube" \
           "NAMESPACE_NAME          ${NAMESPACE_NAME}" \
           "POD_NAME                ${pod_name}" \
           "CONTAINER_NAME          nginx" \
           "CONTAINER_IMAGE_NAME    nginx:latest" \
           "DEPLOYMENT_NAME         ${DUMMY_DEPLOYMENT_NAME}"
do
    k=$(echo "$PAIR" | awk '{ print $1 }')
    v=$(echo "$PAIR" | awk '{ print $2 }')
    if ! echo "$env_vars" | grep -q "${ENV_VARS_PREFIX}_${k}=${v}$"; then
        errors="${errors}\n${ENV_VARS_PREFIX}_${k}=${v} is not present"
    fi
done

if [ -n "$errors" ]; then
    printf "Test errors:%s\n" "$errors"
    exit 1
else
    printf "Tests are passing successfully\n\n"
fi
