# ♜ Rook NooBaa Design

NooBaa is an object data service for hybrid and multi cloud environments. 

NooBaa runs on kubernetes, providing an S3 object store service (and Lambda with bucket triggers) to clients both inside and outside the cluster, and using storage resources from within or outside the cluster, with flexible placement policies to automate data use cases.

- https://github.com/noobaa/noobaa-core       NooBaa Github Project
- https://www.youtube.com/watch?v=fuTKXBMwOes NooBaa From Zero to Multi Cloud
- https://www.youtube.com/watch?v=uW6NvsYFX-s NooBaa on Ceph Tech Talk

## Design Topics

1. [System CRD](system-crd.md)
2. Operator features:
    1. [System install](system-install.md) - Install a new system
    2. [System upgrade](system-upgrade.md) - Upgrade a system to a new release
    3. [System status](system-status.md) - Report the system status
    4. [ObjectBucketClaim provisioner](noobaa-obc.md) - provision a bucket to application claim
    5. [OperatorLifecycleManager integration](noobaa-olm.md)
    6. [Adding/removing storage resources](storage-resources.md)
    7. [Defining storage classes](storage-classes.md)
    8. [Scaling up/down S3 endpoints](endpoint-scale.md)
    9. [Multi-cluster federation](multi-cluster.md)
3. Manual operations - The operator command-line utility will provide additional manual operations:
    1. [Operator installer](operator-installer.md) - install the operator
    2. [Open management console](mgmt-console.md) - open the management console
    3. [Diagnostics package](diagnostics.md)
