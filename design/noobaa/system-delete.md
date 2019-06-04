♜ [Rook NooBaa Design](README.md) /
# System Delete

The operator will detect deletion of a system CR and will folloup by deleting all the owned resources.

In general we can leave this to Garbage Collection as described here:

https://kubernetes.io/docs/concepts/workloads/controllers/garbage-collection/

While this works, and there are no known dependencies on the order to deletions, we will prefer to call it explicitly and leave the GC to edge cases.
