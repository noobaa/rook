# ♜ Rook NooBaa Design / System Create

The operator watches for NooBaa System CRD creations and in reaction will deploy a new system.

## Create Flow

- Kubernetes resources:
    - Role, RoleBinding, ServiceAccount, StatefulSet, Service.
- NooBaa owner-account:
    - With s3 credentials and UI user/password.
- Underlying storage resources:
    - For every `spec.storageResources` the operator will claim the resource and configure it in NooBaa.
    - Each underlying OBC will be provisioned and used to configure a cloud resource in NooBaa.
    - Each underlying PVC will be mounted with a NooBaa storage agent (aka daemon).
- StorageClass:
    - For every `spec.storageClasses` the operator will create a StorageClass provisioned by the operator from a target system.
    - Each storage class will refer to a noobaa system name, and specify a tiering policy with references to which underlying resources to use.
    - Applications that would use a NooBaa-OBC will refer to a storage class name, and the operator will use that to satisfy the bucket claim.
    - **TODO:** StorageClasses are global and have a unique name across all namespaces, therefore we should decide who is responsible for these names - the cluster administrator or the NooBaa admin.
- Endpoints:
    - Create a horizontal-pod-autoscale group of NooBaa endpoints.
- Create `first.bucket`:
    - **TODO:**  Using internal storage?
- Ready to serve:
    - NooBaa-System is serving S3 requests
    - Rook-NooBaa-Operator is serving NooBaa-OBC requests
