import ...lib.k8s_wrapper
import ...lib.helm_flux
import ...lib.vault
import ..shared
import yaml

out_dir = shared.cluster_path + "infrastructure/authentik/"

# Filenames
kustomization_file = "kustomization.yaml"
namespace_file = "authentik.yaml"
release_file = "release.yaml"
repository_file = "repository.yaml"
authentik_secret_file = "authentik-secret.yaml"
authentik_postgres_secret_file = "authentik-postgres-secret.yaml"
vault_auth_file = "vault-auth.yaml"
vault_rbac_file = "vault-rbac.yaml"

# k8s configuration
kustomization = k8s_wrapper.create_kustomization(
    "authentik",
    [
        namespace_file,
        release_file,
        repository_file,
        authentik_secret_file,
        authentik_postgres_secret_file,
        vault_auth_file,
        vault_rbac_file
    ]
)
namespace = k8s_wrapper.create_namespace("authentik")
# Helm configuration
repo = helm_flux.create_helm_repository(
    "https://charts.goauthentik.io",
    "authentik",
)
release = helm_flux.create_helm_release(
    {
       authentik: {
            error_reporting: {
                enabled: True
            },
            server: {
                ingress: {
                    enabled: False
                }
            }
            postgresql: {
                enabled: True
            }
            redis: {
                enabled: True
            }
       }
    },
    "authentik",
    "authentik",
    "authentik",
    "authentik",
     [
        {
            kind: "Secret",
            name: "authentik-secret",
            targetPath: "authentik.secret_key",
            valuesKey: "key",
            optional: False
        }, 
        {
            kind: "Secret",
            name: "authentik-postgres-secret",
            targetPath: "postgresql.auth.password",
            valuesKey: "key",
            optional: False
        }
     ]
    
)
# Vault configuration
vault_auth = vault.create_vault_auth(
    "authentik",
    "authentik-vault-serviceaccount",
    True,
    "vault-secrets-operator"
)

vault_rbac = vault.create_vault_rbac(
    "authentik",
)
authentik_secret = vault.create_vault_static_secret(
    "authentik-vault-secret",
    vault_auth.metadata.name,
    "kvv2",
    "kv-v2",
    "authentik/authentik",
    "10s",
    "authentik-secret"
)
authentik_postgres_secret = vault.create_vault_static_secret(
    "authentik-postgres-vault-secret",
    vault_auth.metadata.name,
    "kvv2",
    "kv-v2",
    "authentik/postgres",
    "10s",
    "authentik-postgres-secret"
)
# Write to files
yaml.dump_to_file(kustomization, out_dir + kustomization_file)
yaml.dump_to_file(namespace, out_dir + namespace_file)
yaml.dump_to_file(repo, out_dir + repository_file)
yaml.dump_to_file(release, out_dir + release_file)
yaml.dump_to_file(authentik_secret, out_dir + authentik_secret_file)
yaml.dump_to_file(authentik_postgres_secret, out_dir + authentik_postgres_secret_file)
yaml.dump_to_file(vault_auth, out_dir + vault_auth_file)
yaml.dump_all_to_file(vault_rbac, out_dir + vault_rbac_file)