# ♜ Rook NooBaa Design / System CRD

The operator will define a new kubernetes API group `noobaa.rook.io`.

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

## Default Spec

The operator will provide defaults when a system created without a spec:

```yaml
apiVersion: noobaa.rook.io/v1alpha1
kind: System
metadata:
  name: noobaa
  namespace: rook-noobaa
```

