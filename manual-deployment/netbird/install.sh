kubectl create namespace netbird

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.17.0/cert-manager.yaml

helm install netbird -f ./values.yaml -n netbird netbirdio/netbird-operator