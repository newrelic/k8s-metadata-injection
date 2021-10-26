#!/usr/bin/env sh
set -e

printf 'bootstrapping starts:\n'
# shellcheck disable=SC1090
. "$(dirname "$0")/k8s-e2e-bootstraping.sh"
printf 'bootstrapping complete\n'

WEBHOOK_LABEL="app.kubernetes.io/name=nri-metadata-injection,app.kubernetes.io/instance=nri-mdi"
DEPLOYMENT_NAME="nri-mdi-nri-metadata-injection"
DUMMY_POD_LABEL="app=dummy"
DUMMY_DEPLOYMENT_NAME="dummy-deployment"
ENV_VARS_PREFIX="NEW_RELIC_METADATA_KUBERNETES"

finish() {
    helm uninstall nri-mdi

    printf "webhook logs:\n"
    kubectl logs "$webhook_pod_name"
}

# ensure that we build docker image in minikube
[ "$E2E_MINIKUBE_DRIVER" = "none" ] || eval "$(minikube docker-env --shell bash)"

# build webhook docker image

# Set GOOS and GOARCH explicitly since Dockerfile expects them in the binary name
(
    cd ..
    GOOS=linux GOARCH=amd64 make compile build-container
    cd -
)

trap finish EXIT

# install the metadata-injection webhook
helm repo add newrelic https://helm-charts.newrelic.com
helm install nri-mdi newrelic/nri-metadata-injection --set cluster=YOUR-CLUSTER-NAME
if [ $? -ne 0 ]; then
    printf "Helm failed to install this release\n"
fi

webhook_pod_name=$(get_pod_name_by_label "$WEBHOOK_LABEL")
if [ "$webhook_pod_name" = "" ]; then
    printf "not found any pod with label %s\n" "$WEBHOOK_LABEL"
    kubectl get deployments
    kubectl describe deployment "$DEPLOYMENT_NAME"
    kubectl get pods
    exit 1
fi
wait_for_pod "$webhook_pod_name"
sleep 10 # Wait 10 second more so the API Server can settle its cache 

### Testing

# deploy a pod
kubectl create -f manifests/deployment.yaml

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
printf "+ date"
date
printf "+ kubectl exec \"${pod_name}\" -- env | grep \"${ENV_VARS_PREFIX}\""
env_vars="$(kubectl exec "${pod_name}" -- env | grep "${ENV_VARS_PREFIX}")"
set -e
printf "\nInjected environment variables:\n"
printf "%s\n" "$env_vars"

errors=""
for PAIR in \
           "CLUSTER_NAME            YOUR-CLUSTER-NAME" \
           "NODE_NAME               minikube" \
           "NAMESPACE_NAME          default" \
           "POD_NAME                ${pod_name}" \
           "CONTAINER_NAME          busybox" \
           "CONTAINER_IMAGE_NAME    busybox:latest" \
           "DEPLOYMENT_NAME         dummy-deployment"
do
    k=$(echo "$PAIR" | awk '{ print $1 }')
    v=$(echo "$PAIR" | awk '{ print $2 }')
    if ! echo "$env_vars" | grep -q "${ENV_VARS_PREFIX}_${k}=${v}$"; then
        errors="${errors}\n${ENV_VARS_PREFIX}_${k}=${v} is not present"
    fi
done

if [ -n "$errors" ]; then
    printf "Test errors:$errors\n"
    exit 1
else
    printf "Tests are passing successfully\n\n"
fi
