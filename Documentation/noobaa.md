# NooBaa

## User Guide

### Install Operator
```
kubectl create -f https://raw.githubusercontent.com/rook/rook/master/cluster/examples/kubernetes/noobaa/operator.yaml
```

### Create a new NooBaa System
```
kubectl create -f https://...
```

### Claim an object bucket
```
kubectl create -f https://...
```

### Use the bucket from an application
```
kubectl create -f https://...
```

### Upgrade to latest version ( *** same as install *** )
```
kubectl create -f https://raw.githubusercontent.com/rook/rook/master/cluster/examples/kubernetes/noobaa/operator.yaml
```

### Delete NooBaa System
```
kubectl delete system.noobaa.rook.io/system-0 -n rook-noobaa
```

### Uninstall Operator
```
kubectl delete -f https://raw.githubusercontent.com/rook/rook/master/cluster/examples/kubernetes/noobaa/operator.yaml
```

### Management Console with port-forward
```
kubectl port-forward -n noobaa-rook noobaa-server-0 8080:5555
# the previous command runs to proxy your browser, start a new terminal to run:
open 127.0.0.1:5555
```

### Management Console with externalIP
```
kubectl describe noobaa -n noobaa-rook

...
Internal        :
  Management    : http://mgmt.rook-noobaa
  S3            : http://s3.rook-noobaa

External:
  ManagementExt : https://externalIP:8443/
  S3Ext         : https://externalIP/
...
