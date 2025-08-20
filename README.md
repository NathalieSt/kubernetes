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

## Creating a token for kiali
```sh
kubectl create token kiali-service-account -n istio-system
```

## Copy cert from caddy
For Firefox p12 is required to import,
for OwnTracks on Android to work only the root.crt needs to be imported as CA.
```sh
kubectl cp <pod-name>:/data/caddy/pki/authorities/local/root.crt ./root.crt -n caddy
kubectl cp <pod-name>:/data/caddy/pki/authorities/local/root.key ./root.key -n caddy
openssl pkcs12 -export -in root.crt -inkey root.key -out caddy-server.p12
```

## Copying files from/to vault
```sh
# copy from
kubectl cp <pod-name>:/vault/data ./vault/ -n vault
# copy to
kubectl cp ./vault-backup/ <pod-name>:/vault/data -n vault
```
## Mounting nfs
```sh
sudo mount -t nfs <server name>:<location on server> <local location>
```