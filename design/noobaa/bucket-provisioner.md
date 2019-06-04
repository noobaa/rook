♜ [Rook NooBaa Design](README.md) /
# Bucket Claim

Kubernetes natively supports dynamic provisioning for many types of file and block storage, but lacks support for object bucket provisioning.
In order to provide a native provisioning of object storage buckets, the concept of Object Bucket Claim (OBC/OB) was introduced in a similar manner to Persistent Volume Claim (PVC/PV)

See the following repository providing the concept design and an embeddable library code to simplify the implementation:

- [Object Bucket Provisioning Design](https://github.com/yard-turkey/lib-bucket-provisioner/blob/master/doc/design/object-bucket-lib.md)

Here are the details of the implementation in the rook noobaa operator:
- The operator will create object storage classes as part of the system creation.
- The operator will watch for OBC's on the created classes and fulfill the claims.
- A claim for a new bucket will create a bucket in NooBaa using the configured bucketPolicyTemplateName from `system.spec.bucketPolicies.templates`.
- The bucket policy template specifies which backing stores to use and how to use it by specifying the tiering policy, mirroring etc.
- Applications that require a bucket will create an OBC with its deployment and refer to a storage class name.
- The operator will share a config map and a secret with the application in order to give it all the needed details to work with the bucket.


### Example `StorageClass` created by the operator:

See https://github.com/yard-turkey/lib-bucket-provisioner/blob/master/deploy/storageClass.yaml 

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: noobaa-cloud-mirror
provisioner: noobaa.rook.io/bucket
parameters:
  bucketPolicyTemplateName: cloud-mirror
reclaimPolicy: Delete
```

### Example `OBC`:

See https://github.com/yard-turkey/lib-bucket-provisioner/blob/master/deploy/example-claim.yaml 

```yaml
apiVersion: objectbucket.io/v1alpha1
kind: ObjectBucketClaim
metadata:
  name: my-bucket-claim
  namespace: default
spec:
  generateBucketName: "my-bucket-"
  storageClassName: noobaa-cloud-mirror
  SSL: false
```
