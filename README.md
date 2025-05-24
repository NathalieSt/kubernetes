# Kubernetees
## Delete terminating namespace
```sh
kubectl get namespace "netbird" -o json \
  | tr -d "\n" | sed "s/\"finalizers\": \[[^]]\+\]/\"finalizers\": []/" \
  | kubectl replace --raw /api/v1/namespaces/netbird/finalize -f -
```
## Delete stuck CRD
```sh
kubectl patch crd nbgroups.netbird.io -p '{"metadata":{"finalizers":[]}}' --type=merge
```

```sh
flux bootstrap git \
  --components-extra=image-reflector-controller,image-automation-controller \
  --url=ssh://git@codeberg.org/NathalieStiefsohn/kubernetes.git \
  --branch=main \
  --private-key-file=./id_ed25519 \
  --path=cluster/flux
```