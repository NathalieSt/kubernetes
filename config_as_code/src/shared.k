schema Netbird:
    setup_key: str
    key: str

cluster_path = "../../../../cluster/"
global_vault_auth_namespace = "vault-secrets-operator"
ca_issuer = "my-ca-issuer"
netbird: Netbird = {
    setup_key = "netbird-setup-key"
    key = "setup-key"
}
tls_cert_secret = "nginx-tls-secret"