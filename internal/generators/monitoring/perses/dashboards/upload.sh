cat "Applying datasource ..."
cat datasource.yaml | yq  | curl -s -X POST http://127.0.0.1:44293/api/v1/projects/kubernetes/datasources -H "Content-Type: application/json" -d @-
# requires yq and curl
for f in cluster-overview cluster-nodes cluster-workloads cluster-control-plane; do
  echo "→ Applying $f ..."
  cat $f.yaml | yq  | curl -s -X POST http://127.0.0.1:44293/api/v1/projects/kubernetes/dashboards \
    -H "Content-Type: application/json" \
    -d @-
  echo ""
done