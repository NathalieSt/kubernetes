# First, remove finalizers from all netbird.io resources
kubectl get crd -o name | grep netbird.io | xargs -I {} kubectl patch {} -p '{"metadata":{"finalizers":[]}}' --type=merge

# Then delete all netbird.io CRDs
kubectl delete crd --all --selector=app.kubernetes.io/name=netbird || true
kubectl get crd -o name | grep netbird.io | xargs kubectl delete --force --grace-period=0

# If still stuck, force delete each CRD individually
kubectl get crd | grep netbird.io | awk '{print $1}' | xargs -I {} kubectl delete crd {} --force --grace-period=0

# If still stuck, force delete each CRD individually
kubectl get crd | grep netbird.io | awk '{print $1}' | xargs -I {} kubectl patch crd {} -p '{"metadata":{"finalizers":[]}}' --type=merge


# If specific resource still remains stuck
kubectl patch NBRoutingPeer router -n netbird \
  -p '{"metadata":{"finalizers":[]}}' --type=merge

kubectl delete NBRoutingPeer router -n netbird --force --grace-period=0



# First, remove finalizers from all netbird.io resources
kubectl get crd -o name | grep netbird.io | xargs -I {} kubectl patch {} -p '{"metadata":{"finalizers":[]}}' --type=merge

# Remove finalizers from all instances of netbird.io CRDs
kubectl get crd | grep netbird.io | awk '{print $1}' | xargs -I {} kubectl get {} --all-namespaces -o name | xargs -I {} kubectl patch {} -p '{"metadata":{"finalizers":[]}}' --type=merge

# Force delete all instances of netbird.io CRDs
kubectl get crd | grep netbird.io | awk '{print $1}' | xargs -I {} kubectl delete {} --all --all-namespaces --force --grace-period=0

# Then delete all netbird.io CRDs
kubectl delete crd --all --selector=app.kubernetes.io/name=netbird || true
kubectl get crd -o name | grep netbird.io | xargs kubectl delete --force --grace-period=0

# If still stuck, force delete each CRD individually
kubectl get crd | grep netbird.io | awk '{print $1}' | xargs -I {} kubectl delete crd {} --force --grace-period=0

# delete all services which have finalizers

kubectl get services --all-namespaces -o json | jq -r '.items[] | select(.metadata.finalizers[]? == "netbird.io/cleanup") | "\(.metadata.namespace)/\(.metadata.name)"' | xargs -I {} bash -c 'NS=$(echo {} | cut -d/ -f1); NAME=$(echo {} | cut -d/ -f2); kubectl patch service $NAME -n $NS -p "{\"metadata\":{\"finalizers\":[]}}" --type=merge && kubectl delete service $NAME -n $NS --force --grace-period=0'

kubectl delete namespace cert-manager