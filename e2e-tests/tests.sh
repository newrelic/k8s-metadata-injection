#!/usr/bin/env sh
set -e

# shellcheck disable=SC1090
. "$(dirname "$0")/k8s-e2e-bootstraping.sh"

finish() {
    printf "calling cleanup function\n"
    kubectl delete -f ../deploy/ || true
    kubectl delete -f manifests/ || true
}

# build webhook docker image
(
    cd ..
    make build-container
    cd -
)

trap finish EXIT

# install the metadata-injection webhook
kubectl create -f ../deploy/job.yaml
awk '/image: / { print; print "        imagePullPolicy: Never"; next }1' ../deploy/newrelic-metadata-injection.yaml | kubectl create -f -

label="app=newrelic-metadata-injection"
webhook_pod_name=$(get_pod_name_by_label "$label")
if [ "$webhook_pod_name" = "" ]; then
    printf "not found any pod with label %s\n" "$label"
    kubectl get deployments
    kubectl describe deployment newrelic-metadata-injection-deployment
    kubectl get pods
    exit 1
fi
wait_for_pod "$webhook_pod_name"

### Testing

# deploy a pod
kubectl create -f manifests/deployment.yaml

label="app=dummy"
pod_name="$(get_pod_name_by_label "$label")"
if [ "$pod_name" = "" ]; then
    printf "not found any pod with label %s" "$label"
    kubectl describe deployment dummy-deployment
    exit 1
fi
wait_for_pod "$pod_name"

printf "webhook logs:\n"
kubectl logs "$webhook_pod_name"

kubectl get pods
kubectl describe pod "${pod_name}"
printf "getting env vars for %s\n" "${pod_name}"
kubectl exec "${pod_name}" env
env_vars="$(kubectl exec "${pod_name}" env | grep "NEW_RELIC_METADATA_KUBERNETES")"
printf "\nInjected environment variables:\n"
printf "%s\n" "$env_vars"

errors=""
for PAIR in \
           "CLUSTER_NAME      <YOUR_CLUSTER_NAME>" \
           "NODE_NAME         minikube" \
           "NAMESPACE_NAME    default" \
           "POD_NAME          ${pod_name}" \
           "CONTAINER_NAME    busybox" \
           "DEPLOYMENT_NAME   dummy-deployment"
do
    k=$(echo "$PAIR" | awk '{ print $1 }')
    v=$(echo "$PAIR" | awk '{ print $2 }')
    if ! echo "$env_vars" | grep -q "NEW_RELIC_METADATA_KUBERNETES_${k}=${v}"; then
        errors="${errors}\nNEW_RELIC_METADATA_KUBERNETES_${k}=${v} is not present"
    fi
done

if [ -n "$errors" ]; then
    printf "Test errors:\n"
    printf '%s\n' "$errors"
else
    printf "Tests are passing successfully\n\n"
fi
