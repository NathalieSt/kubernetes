# Cause netbird really doesnt like to be uninstalled...

kubectl api-resources --verbs=list -o name | \
      xargs -I {} kubectl get {} --all-namespaces -o name 2>/dev/null | \
      grep -i netbird | \
while read resource; do
    kubectl patch $resource -p '{"metadata":{"finalizers":[]}}' --type=merge 2>/dev/null
    kubectl delete $resource --ignore-not-found --force --grace-period=0 2>/dev/null
done

kubectl get crd -o name | grep -i netbird | \
while read crd; do
    kubectl patch $crd -p '{"metadata":{"finalizers":[]}}' --type=merge
    kubectl delete $crd --ignore-not-found --force --grace-period=0
done    

kubectl get namespace netbird -o json \
  | jq '.spec.finalizers=[]' \
  | kubectl replace --raw "/api/v1/namespaces/netbird/finalize" -f -