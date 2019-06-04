♜ [Rook NooBaa Design](README.md) /
# System Create

The operator watches for NooBaa System CRD creations and in reaction deploys a new system with the following steps:

- Kubernetes resources:
    - Role, RoleBinding, ServiceAccount, StatefulSet, Service.
- NooBaa owner-account:
    - With s3 credentials and UI user/password.
    - Save owner auth token in a secret of the operator for subsequent operations.
- Backing Stores:
    - For every `system.spec.backingStores` the operator will setup the storage resource per the specified type.
    - It will then configure it in NooBaa using API calls to the NooBaa brain.
    - See [Backing stores](backing-stores.md).
- Bucket Provisioner:
    - For every `system.spec.bucketProvisioner.exportedStorageClasses` the operator will create a StorageClass for Object Bucket Provisioning.
    - See [Bucket provisioner](bucket-provisioner.md).
- Ready to serve:
    - NooBaa-System is serving S3 requests
    - Rook-NooBaa-Operator is serving NooBaa-OBC requests
