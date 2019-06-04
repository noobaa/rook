# ♜ Rook NooBaa Design

> NooBaa is an object data service for hybrid and multi cloud environments.\
NooBaa runs on kubernetes, provides an S3 object store service (and Lambda with bucket triggers) to clients both inside and outside the cluster, and uses storage resources from within or outside the cluster, with flexible placement policies to automate data use cases.

Overview:

- [NooBaa Github Project](https://github.com/noobaa/noobaa-core)
- [NooBaa From Zero to Multi Cloud (youtube)](https://www.youtube.com/watch?v=fuTKXBMwOes)
- [NooBaa on Ceph Tech Talk (youtube)](https://www.youtube.com/watch?v=uW6NvsYFX-s)
- [NooBaa Architecture Slide](media/noobaa-architecture.png)
- [Object Bucket Provisioning Design](https://github.com/yard-turkey/lib-bucket-provisioner/blob/master/doc/design/object-bucket-lib.md)

## Operator Design Topics

1. [System CRD](system-crd.md) - Specify NooBaa system structure.
1. [System create](system-create.md) - Install a new system.
1. [System read](system-read.md) - Report the system status.
1. [System update](system-update.md) - Update system spec.
1. [System delete](system-delete.md) - Uninstall a system.
1. [Backing stores](backing-stores.md)
1. [Bucket policy templates](bucket-policy.md)
1. [Bucket provisioner](bucket-provisioner.md) - Provision a bucket to application claim
1. [Upgrades](upgrades.md) - Upgrade the operator and systems to a new release
1. **TODO:** [Security](security.md)
1. **TODO:** [Operator Lifecycle Manager integration](noobaa-olm.md)
1. **TODO:** [Scaling up/down S3 endpoints](endpoint-scale.md)
1. **TODO:** [Multi-cluster federation](multi-cluster.md)
