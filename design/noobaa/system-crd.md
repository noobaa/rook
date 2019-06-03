# ♜ Rook NooBaa Design / System CRD

The operator will define a new kubernetes API group `noobaa.rook.io`.

In that group it will define a CRD `System` representing a NooBaa system.

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: systems.noobaa.rook.io
spec:
  scope: Namespaced
  group: noobaa.rook.io
  version: v1alpha1
  names:
    kind: System
    listKind: SystemList
    singular: system
    plural: systems
    shortNames:
      - nb
      - noobaa
```

## Minimal Spec

The operator will provide working defaults when a system created without a spec:

```yaml
apiVersion: noobaa.rook.io/v1alpha1
kind: System
metadata:
  name: noobaa
  namespace: rook-noobaa
```

## Full Spec Features

```yaml
apiVersion: noobaa.rook.io/v1alpha1
kind: System
metadata:
  name: noobaa-1
  namespace: rook-noobaa
spec:

  storageResources:

    - name: noobaa-ceph-resource
      region: on-prem-ohio
      type: obc
      storageClassName: ceph-rgw-bucket-provisioner
      # bucketName: rgw-target-bucket-name

    - name: noobaa-aws-resource
      region: aws-paris
      type: obc
      storageClassName: aws-bucket-provisioner
      # bucketName: aws-target-bucket-name

    - name: local-pv-storage
      region: on-prem-ohio
      type: pvc
      storageClassName: ceph-block-provisioner
      size: "30 GB"
      count: 3

    # - name: s3-storage
    #   region: 
    #   type: s3-bucket
    #   endpoint: 
    #   credentialsSecretName: xxx
    #   bucketName: aws-target-bucket-name

  storageClasses:

    # a class that uses a single underlying resource
    - name: noobaa-local-storage
      storageResourceName: local-pv-storage

    # a class that uses multiple resources in mirror policy
    - name: noobaa-cloud-mirror
      tiering:
        - mirroring:
          - storageResourceName: noobaa-ceph-resource
          - storageResourceName: noobaa-aws-resource

    # a class that uses multiple resources in tiering policy
    - name: noobaa-tiering-to-cloud
      tiering:
        - mirroring:
          - storageResourceName: noobaa-ceph-resource
        - mirroring:
          - storageResourceName: noobaa-aws-resource

  endpoints:
    maxAutoScale: 100000
```


## Sample `StorageClass`

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: noobaa-cloud-mirror
provisioner: s3.noobaa.io/bucket
parameters:
  noobaaSystemName: noobaa-1
reclaimPolicy: Delete
```

## Sample `OBC`

See https://github.com/yard-turkey/lib-bucket-provisioner/blob/master/deploy/example-claim.yaml 

```yaml
apiVersion: objectbucket.io/v1alpha1
kind: ObjectBucketClaim
metadata:
  name: "example-bucket"
  namespace: default
spec:
  generateBucketName: "objectbucket-io-"
  storageClassName: noobaa-cloud-mirror
  SSL: false
  versioned: false
```
