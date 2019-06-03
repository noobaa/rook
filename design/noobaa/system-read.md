# ♜ Rook NooBaa Design / System Read

## System Status

The operator will update the status of the NooBaa system to represent the current state of its operation on that system.

Here is the planned status structure as would be returned by a `kubectl get noobaa noobaa-1 -o yaml`:

```yaml
apiVersion: noobaa.rook.io/v1alpha1
kind: System
metadata:
  name: noobaa-1
  namespace: rook-noobaa
spec:
    # ...
status:
    # TODO
```
