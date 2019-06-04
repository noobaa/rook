♜ [Rook NooBaa Design](README.md) /
# System Read

The operator will set the status of the NooBaa system to represent the current state of reconciling to the desired state.\
Here is the example status structure as would be returned by a `kubectl get noobaa <system-name> -o yaml`:

```yaml
apiVersion: noobaa.rook.io/v1alpha1
kind: System
metadata:
  name: noobaa-1
  namespace: rook-noobaa
spec:
    # ...
status:

  health:
    accounts: OK
    backingStores: OK
    bucketPolicies: OK
    bucketProvisioner: OK
    issues: []

  counters:
    accounts: 3
    backingStores: 2
    bucketPolicies:
      templates: 3
    bucketProvisioner:
      exportedStorageClasses: 3
    buckets: 33
  
  readme: |

Welcome to NooBaa

Management UI
-------------
- Username/password     : admin@noobaa.io / TgerOAZ9MOd/Xg== 
- External address      : https://111.111.111.111:8443
- Node port address     : http://192.168.99.100:30785
- Cluster local address : https://noobaa-mgmt.noobaa:8443
- Port forwarding       : kubectl port-forward -n noobaa service/noobaa-mgmt 11443:8443 # open https://localhost:11443

S3 Endpoint
-----------
- Access/secret         : iqDe8ubjD26kPJXgjvlR / PEh03V5ed3fFxriJIHwcoZc5A8+sshdhVO3LINBj 
- External address      : https://222.222.222.222:8443
- Node port address     : http://192.168.99.100:30361
- Cluster local address : https://s3.noobaa
- Port forwarding       : kubectl port-forward -n noobaa service/s3 10443:443 # open https://localhost:10443
- aws-cli               : alias s3='AWS_ACCESS_KEY_ID=iqDe8ubjD26kPJXgjvlR AWS_SECRET_ACCESS_KEY=PEh03V5ed3fFxriJIHwcoZc5A8+sshdhVO3LINBj aws --endpoint https://localhost:10443 s3'
```

Example health status when there is an issue with the availability of a backing store:
```yaml
status:
  health:
    accounts: OK
    backingStores: WARNING
    bucketPolicies: OK
    bucketProvisioner: OK
    issues:
      - title: backingStore "aws" is not accessible
        createTime: "2019-06-04T13:05:35.473Z"
        lastTime: "2019-06-04T13:05:35.473Z"
```
