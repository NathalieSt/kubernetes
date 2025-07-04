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

## Bootstrap flux
```sh
flux bootstrap git \
  --components-extra=image-reflector-controller,image-automation-controller \
  --url=ssh://git@codeberg.org/NathalieStiefsohn/kubernetes.git \
  --branch=main \
  --private-key-file=./id_ed25519 \
  --path=cluster/flux
```

## Remove stuck pod
```sh
kubectl delete pod <podname> -n <namespace> --force --grace-period=0
```

## Delete netbird operator
```sh
helm uninstall netbird-operator -n netbird
```

## Trigger Cronjob manually
```sh
kubectl create job --from=cronjob/postgres-postgresql-pgdumpall manual-postgres -n postgres
```