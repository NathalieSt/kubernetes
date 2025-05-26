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