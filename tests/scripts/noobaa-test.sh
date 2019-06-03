# build commands

function cmd_dev {
    make IMAGES=noobaa PLATFORM=darwin GO_STATIC_PACKAGES=github.com/rook/rook/cmd/rook go.build && \
        _output/bin/darwin/rook noobaa operator
}

function cmd_build {
    if ! docker info | grep 'Name: minikube'
    then
        echo 'NOTE: To build directly on minikube use: eval $(minikube docker-env)'
    fi
    make IMAGES=noobaa || exit 1
    . build/common.sh
    docker tag ${BUILD_REGISTRY}/noobaa-amd64 rook/noobaa:master
}

function cmd_unbuild {
    make clean
    rm -rf vendor/ pkg/client/
    dep ensure
    make codegen || exit 1
}

function cmd_rebuild {
    cmd_unbuild
    cmd_build
}

# operator commands

function cmd_install {
    kubectl create -f cluster/examples/kubernetes/noobaa/operator.yaml
}

function cmd_uninstall {
    kubectl delete -f cluster/examples/kubernetes/noobaa/operator.yaml
    kubectl delete crd systems.noobaa.rook.io
}

function cmd_reinstall {
    cmd_uninstall
    cmd_install
}

function cmd_restart {
    kubectl delete pod -l app=rook-noobaa-operator -n rook-noobaa-operator
}

function cmd_logs {
    kubectl logs --tail=1000 -l app=rook-noobaa-operator -n rook-noobaa-operator
}

# noobaa commands

function cmd_create {
    kubectl create -f cluster/examples/kubernetes/noobaa/system.yaml
}

function cmd_delete {
    kubectl delete -f cluster/examples/kubernetes/noobaa/system.yaml
}

function cmd_recreate {
    cmd_delete
    cmd_create
}

# general commands

function cmd_cleanup {
    cmd_delete
    cmd_uninstall
    cmd_unbuild
}

function cmd_status {
    RESOURCES="noobaa,all,pvc,serviceaccounts,rolebindings,secrets,configmaps"
    echo
    echo "*******************"
    echo "*** rook-noobaa ***"
    echo "*******************"
    echo
    kubectl get $RESOURCES -n rook-noobaa -o wide --show-kind
    echo
    echo "****************************"
    echo "*** rook-noobaa-operator ***"
    echo "****************************"
    echo
    kubectl get $RESOURCES -n rook-noobaa-operator -o wide --show-kind
    echo
}

# Execute the operator manually - not so useful but might be in some cases

function cmd_exec_operator {
    kubectl run rook-noobaa-operator \
        -it \
        --rm --restart=Never \
        --image=rook/noobaa:master \
        --serviceaccount=rook-noobaa-operator \
        -n rook-noobaa-operator \
        --env="POD_NAME=rook-noobaa-operator" \
        --env="POD_NAMESPACE=rook-noobaa-operator" \
        --env="ROOK_LOG_LEVEL=INFO" \
        -- noobaa operator
}

if [ "$(type -t cmd_$1)" = "function" ]
then
    cmd=$1
    shift
    cmd_$cmd $*
else
    echo "Commands:"
    echo "========"
    set | egrep '^cmd_[a-z]+ ()' | cut -d'_' -f2- | cut -d' ' -f1 | sort
    exit 1
fi
