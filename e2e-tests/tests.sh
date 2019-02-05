#!/usr/bin/env sh
set -e

E2E_KUBERNETES_VERSION=${E2E_E2E_KUBERNETES_VERSION:-v1.10.0}
E2E_MINIKUBE_VERSION=${E2E_E2E_MINIKUBE_VERSION:-v0.33.1}
E2E_SETUP_MINIKUBE=${E2E_SETUP_MINIKUBE:-}
E2E_SETUP_KUBECTL=${E2E_SETUP_KUBECTL:-}
E2E_START_MINIKUBE=${E2E_START_MINIKUBE:-}
E2E_MINIKUBE_DRIVER=${E2E_MINIKUBE_DRIVER:-virtualbox}
E2E_SUDO=${E2E_SUDO:-}

finish() {
    printf "calling cleanup function\n"
    kubectl delete -f ../deploy/ || true
    kubectl delete -f manifests/ || true
}

setup_minikube() {
    curl -sLo minikube https://storage.googleapis.com/minikube/releases/"$E2E_MINIKUBE_VERSION"/minikube-linux-amd64 \
        && chmod +x minikube \
        && $E2E_SUDO mv minikube /usr/local/bin/
}

setup_kubectl() {
    curl -sLo kubectl https://storage.googleapis.com/kubernetes-release/release/"$E2E_KUBERNETES_VERSION"/bin/linux/amd64/kubectl \
        && chmod +x kubectl \
        && $E2E_SUDO mv kubectl /usr/local/bin/
}

start_minikube() {
    $E2E_SUDO minikube start --vm-driver="$E2E_MINIKUBE_DRIVER" --bootstrapper=kubeadm --kubernetes-version="$E2E_KUBERNETES_VERSION" --logtostderr
}

get_pod_name_by_label() {
    pod_name=""
    i=1
    while [ "$i" -ne 10 ]
    do
        pod_name=$(kubectl get pods -l "$1" -o name | sed 's/pod\///g; s/pods\///g')
        if [ "$pod_name" != "" ]; then
            break
        fi
        sleep 1
        i=$((i + 1))
    done
    printf "%s" "$pod_name"
}

wait_for_pod() {
    set +e
    is_pod_running=false
    i=1
    while [ "$i" -ne 30 ]
    do
        pod_status="$(kubectl get pod "$1" -o jsonpath='{.status.phase}')"

        if [ "$pod_status" = "Running" ]; then
            is_pod_running=true
            printf "pod %s is running\n" "$1"
            break
        fi

        printf "Waiting for pod %s to be running\n" "$1"
        printf "job/newrelic-metadata-setup logs starts:\n"
        kubectl logs job/newrelic-metadata-setup
        printf "job/newrelic-metadata-setup logs ends:\n"
        sleep 3
        i=$((i + 1))
    done
    if [ $is_pod_running = "false" ]; then
        printf "pod %s does not start within 1 minute 30 seconds\n" "$1"
        kubectl describe job/metadata-setup
        kubectl get pods
        kubectl describe pod "$1"
        exit 1
    fi
    set -e
}

### Bootstraping

cd "$(dirname "$0")"

[ -n "$E2E_SETUP_MINIKUBE" ] && setup_minikube

minikube version

[ -n "$E2E_SETUP_KUBECTL" ] && setup_kubectl

export MINIKUBE_WANTREPORTERRORPROMPT=false
export MINIKUBE_HOME=$HOME
export CHANGE_MINIKUBE_NONE_USER=true
mkdir "$HOME"/.kube || true
touch "$HOME"/.kube/config
export KUBECONFIG=$HOME/.kube/config

[ -n "$E2E_START_MINIKUBE" ] && start_minikube

minikube update-context
kubectl config use-context minikube

is_kube_running="false"

set +e
# this for loop waits until kubectl can access the api server that Minikube has created
i=1
while [ "$i" -ne 90 ] # timeout for 3 minutes
do
   kubectl get po 1>/dev/null 2>&1
   if [ $? -ne 1 ]; then
      is_kube_running="true"
      break
   fi

   printf "waiting for Kubernetes cluster up\n"
   sleep 2
   i=$((i + 1))
done

if [ $is_kube_running = "false" ]; then
   minikube logs
   printf "Kubernetes did not start within 3 minutes. Something went wrong.\n"
   exit 1
fi
set -e

kubectl version

# ensure that we build docker image in minikube
[ "$E2E_MINIKUBE_DRIVER" != "none" ] && eval "$(minikube docker-env)"

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
printf "webhook logs:\n"
kubectl logs "$webhook_pod_name"

label="app=dummy"
pod_name="$(get_pod_name_by_label "$label")"
if [ "$pod_name" = "" ]; then
    printf "not found any pod with label %s" "$label"
    kubectl describe deployment dummy-deployment
    exit 1
fi
wait_for_pod "$pod_name"

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
