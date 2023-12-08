#!/usr/bin/env sh

E2E_KUBERNETES_VERSION=${E2E_KUBERNETES_VERSION:-v1.28.3}
E2E_MINIKUBE_DRIVER=${E2E_MINIKUBE_DRIVER:-docker}
E2E_SUDO=${E2E_SUDO:-}

start_minikube() {
    export MINIKUBE_WANTREPORTERRORPROMPT=false
    export MINIKUBE_HOME=$HOME
    export CHANGE_MINIKUBE_NONE_USER=true
    mkdir -p "$HOME"/.kube
    touch "$HOME"/.kube/config
    export KUBECONFIG=$HOME/.kube/config

    printf "Starting Minikube with Kubernetes version %s...\n" "${E2E_KUBERNETES_VERSION}"
    $E2E_SUDO minikube start --driver="$E2E_MINIKUBE_DRIVER" --kubernetes-version="$E2E_KUBERNETES_VERSION"
}

get_pod_name_by_label() {
    pod_name=""
    i=1
    while [ "$i" -ne 10 ]
    do
        pod_name=$(kubectl -n default get pods -l "$1" -o name | sed 's/pod\///g; s/pods\///g')
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
    desired_status=${2:-'Running'}
    is_pod_in_desired_status=false
    i=1
    while [ "$i" -ne 30 ]
    do
        pod_status="$(kubectl -n default get pod "$1" -o jsonpath='{.status.phase}')"
        if [ "$pod_status" = "$desired_status" ]; then
            is_pod_in_desired_status=true
            printf "pod %s is %s\n" "$1" "$desired_status"
            break
        fi

        printf "Waiting for pod %s to be %s\n" "$1" "$desired_status"
        sleep 3
        i=$((i + 1))
    done
    if [ $is_pod_in_desired_status = "false" ]; then
        printf "pod %s does not transition to %s within 1 minute 30 seconds\n" "$1" "$desired_status"
        kubectl -n default get pods
        kubectl -n default describe pod "$1"
        exit 1
    fi
    set -e
}

### Bootstraping

cd "$(dirname "$0")"

start_minikube
minikube version
minikube update-context

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
