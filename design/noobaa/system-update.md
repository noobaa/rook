♜ [Rook NooBaa Design](README.md) /
# System Update

Updates to the `system.spec` represent a desired change requested by the user. \
The operator will watch for such updates and react by reconciling them in a safe manner, \
or set the `system.status` to represent the problem in achieving the desired state.

The following sections describe non-trivial cases:

## Removal of backing stores

Since backing stores have data on them, it is usually not possible to remove them immediately.\
Instead they switch to a *decommissioning* state, in which NooBaa will attempt to rebuild the data to a new backing store location.\
Once the decomissioning process completes it will clean up from the status.

However there are cases where the decommissioning cannot complete due to inability to read the data from the backing store that is already not serving - for example if the target bucket was already deleted or the credentials were invalidated or there is no network from the system to the backing store service. In such cases the system status will be used to report these issues and suggest manual resolutions.

## Update backing store credentials

In case the credentials of a backing store need to be updated due to a periodic security policy or concern, the appropriate secret should be updated by the user, and the operator will be responsible for watching changes in those secrets and propagating the new credential update to the NooBaa system server.

