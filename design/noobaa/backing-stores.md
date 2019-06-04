♜ [Rook NooBaa Design](README.md) /
# Backing Stores

Each item in `system.spec.backingStores` will be provisioned according to its type:

### Type: `aws-s3`

Assume the bucket exists and set the aws credentials in NooBaa.

### Type: `s3-compatible`

Assume the bucket exists and set the endpoint and aws credentials in NooBaa.

### Type: `google-cloud-storage`

Assume the bucket exists and set the google credentials in NooBaa.

### Type: `azure-blob`

Assume the bucket exists and set the google credentials in NooBaa.

### Type: `obc`

The operator will create a claim and the appropriate provisioner will create a new bucket or connect to existing one depending on the obc options. \
Once the claim is ready its details will be used to configure a cloud resource in NooBaa.

### Type: `pvc`

Create a NooBaa storage agent StatefulSet with PVC mounted in each pod.\
Each agent will connect to the NooBaa brain and provide the PV filesystem storage to be used for storing encrypted chunks of data.

