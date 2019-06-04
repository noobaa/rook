♜ [Rook NooBaa Design](README.md) /
# System CRD

The operator will define a new kubernetes API group `noobaa.rook.io`. \
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

The operator will provide working defaults when a system created with an empty spec. \
It assumes the user expects a simple zero-config deployment, with a simple way to configure the system later from the NooBaa UI. \
This is the minimal system yaml:

```yaml
apiVersion: noobaa.rook.io/v1alpha1
kind: System
metadata:
  name: noobaa
  namespace: rook-noobaa
```

## Full Spec Example

Here is a full blown example of a NooBaa system spec with inline comments.

```yaml
apiVersion: noobaa.rook.io/v1alpha1
kind: System
metadata:
  name: noobaa-1
  namespace: rook-noobaa
spec:

  # Access specifies settings for accessing the system services
  access:

    # Accounts is a list of named accounts that should get access credentials to NooBaa.
    accounts:
      - name: app1-s3-account
        namespace: app1
        s3CredentialsSecretName: app1-s3-credentials
        defaultBackingStoreForNewBuckets: aws-s3

    virtualHost:
      - virtual.host.name

    security:
      certificate: TODO
      disableNonSecureAccess: true

  # Backing stores is a list of storage targets to be used as underlying storage for NooBaa buckets.
  # NooBaa uses these storage targets to store deduped+compressed+encrypted chunks of data (encryption keys are stored separatly).
  # Each item will have a locally unique name to identify it when defining bucket placement policies.
  # Multiple types of storage targets are supported: aws-s3, s3-compatible, google-cloud-storage, azure-blob, obc, pvc.
  # See specific examples below.
  backingStores:

    # AWS S3 bucket
    - name: aws-s3
      type: aws-s3
      region: us-east-1
      bucketName: noobaa1-aws-backing-store
      credentialsSecretName: aws-credentials-secret

    # S3 compatible storage service
    - name: rgw
      type: s3-compatible
      endpoint: s3.rook-ceph
      bucketName: noobaa1-rgw-backing-store
      credentialsSecretName: rgw-credentials-secret
      options:
        sslEnabled: false
        s3ForcePathStyle: true
        signatureVersion: v2

    # Google cloud storage
    - name: gcs
      type: google-cloud-storage
      region: us-west1
      bucketName: noobaa1-gcs-backing-store
      credentialsSecretName: gcs-credentials-secret

    # Azure blob
    - name: azure-blob
      type: azure-blob
      bucketName: noobaa1-azure-blob-backing-store
      credentialsSecretName: azure-blob-credentials-secret

    # OBC - Object Bucket Claim
    - name: aws-obc
      type: obc
      bucketName: aws-obc-backing-store
      storageClassName: aws-obc-provisioner

    # PVC - Persistent Volume Claim
    # This item will create a StatefulSet of NooBaa Agents (pods) each with a PVC of a certain size as specified.
    # NooBaa Agents connect to the noobaa system and allow the system to use the PV filesystem to store encrypted chunks.
    - name: local-pv
      type: pvc
      storageClassName: ceph-block-provisioner
      count: 3
      size: 30
      sizeUnits: GB

  # TODO
  bucketPolicies:

    # 
    templates:

      # a class that uses a single underlying resource
      - name: default
        tiering:
          - mirroring:
            - backingStoreName: noobaa-aws-resource

      # a class that uses multiple resources in mirror policy
      - name: cloud-mirror
        tiering:
          - mirroring:
            - backingStoreName: noobaa-ceph-resource
            - backingStoreName: noobaa-aws-resource

      # a class that uses multiple resources in tiering policy
      - name: tiering-to-cloud
        tiering:
          - mirroring:
            - backingStoreName: noobaa-ceph-resource
          - mirroring:
            - backingStoreName: noobaa-aws-resource

  # TODO
  bucketProvisioner:
    exportedStorageClasses:
      - name: noobaa-object-storage
        bucketPolicyTemplateName: object-storage
      - name: noobaa-cloud-mirror
        bucketPolicyTemplateName: cloud-mirror
      - name: noobaa-tiering-to-cloud
        bucketPolicyTemplateName: tiering-to-cloud

  # TODO
  endpoints:
    min: 1
    maxAutoScale: 10
  
```
