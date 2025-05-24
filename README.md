# Kubernetes
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