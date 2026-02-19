# install server node
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC='--flannel-backend=none --disable-network-policy --disable=traefik' sh -

# register workers
curl -sfL https://get.k3s.io | K3S_URL='https://debian:6443' K3S_TOKEN=<token> sh -

