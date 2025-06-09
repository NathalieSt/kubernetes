import ...lib.k8s_wrapper
import ...lib.helm_flux
import ...lib.vault
import ...lib.cert_manager
import ..shared
import yaml

out_dir = shared.cluster_path + "infrastructure/reflector/"
secrets_sub_dir = "secrets/"
secrets_out_dir = out_dir + secrets_sub_dir

# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "namespace.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
netbird_setup_key_secret_file = "netbird-setup-key-secret.yaml"
cert_manager_certificate_file = "cert-manager-certificate.yaml"
vault_auth_file = "vault-auth.yaml"
vault_rbac_file = "vault-rbac.yaml"

# k8s configuration
kustomization = k8s_wrapper.create_kustomization(
    "reflector",
    [
        namespace_file,
        release_file,
        repository_file,
        # Secrets
        secrets_sub_dir + netbird_setup_key_secret_file,
        secrets_sub_dir + cert_manager_certificate_file,
        secrets_sub_dir + vault_auth_file,
        secrets_sub_dir + vault_rbac_file
    ]
)
namespace = k8s_wrapper.create_namespace("reflector")
# Helm configuration
repo = helm_flux.create_helm_repository(
    "https://emberstack.github.io/helm-charts",
    "reflector",
)
release = helm_flux.create_helm_release(
    {},
    "reflector",
    "reflector",
    "reflector",
    "reflector",
     []
)
# Reflector annotations
reflector_annotations = {
        "reflector.v1.k8s.emberstack.com/reflection-allowed": "true"
        "reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces": "authentik"  # Control destination namespaces
        "reflector.v1.k8s.emberstack.com/reflection-auto-enabled": "true" # Auto create reflection for matching namespaces
        "reflector.v1.k8s.emberstack.com/reflection-auto-namespaces": "authentik" # Control auto-reflection namespaces
}

# Vault configuration
vault_auth = vault.create_vault_auth(
    "reflector",
    "reflector-vault-serviceaccount",
    True,
    shared.global_vault_auth_namespace
)
vault_rbac = vault.create_vault_rbac(
    "reflector",
)
netbird_setup_key_secret = vault.create_vault_static_secret(
    "reflector-vault-secret",
    vault_auth.metadata.name,
    "kvv2",
    "kv-v2",
    "netbird/setup-key",
    "10s",
    shared.netbird.setup_key,
    # Allow reflector to replicate the netbird setup key secret across namespaces
    reflector_annotations 
)
# Cert Manager configuration
cert_manager_certificate = cert_manager.create_certificate(
    "nginx-reverse-proxy-cert", 
    ["cert-manager"], 
    "nginx-tls-secret", 
    shared.ca_issuer, 
    "ClusterIssuer", 
    ["*.cluster.netbird.selfhosted"],
    # Allow reflector to replicate the certificate across namespaces
    reflector_annotations
)
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
# Write helm configuration to files
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
# Write secrets to files
yaml.dump_to_file(cert_manager_certificate, secrets_out_dir + cert_manager_certificate_file)
yaml.dump_to_file(netbird_setup_key_secret, secrets_out_dir + netbird_setup_key_secret_file)
yaml.dump_to_file(vault_auth, secrets_out_dir + vault_auth_file)
yaml.dump_all_to_file(vault_rbac, secrets_out_dir + vault_rbac_file)