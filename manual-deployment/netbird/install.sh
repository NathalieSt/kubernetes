kubectl create namespace netbird
kubectl apply -f ../../../api.yaml

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.17.0/cert-manager.yaml

sleep 2m

helm install --create-namespace -f values.yaml -n netbird netbird-operator netbirdio/kubernetes-operator