import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema Destination:
    create: bool
    name: str
    labels: {str:str} = {}
    annotations: {str:str} = {}

schema VaultStaticSecretSpec:
    vaultAuthRef: str
    mount: str
    type: str
    path: str
    refreshAfter: str
    destination: Destination

schema VaultStaticSecret:
    apiVersion: "secrets.hashicorp.com/v1beta1"
    kind: "VaultStaticSecret"
    metadata: meta.ObjectMeta
    spec: VaultStaticSecretSpec

create_vault_static_secret = lambda secret_name: str, vault_auth_ref: str, mount: str, secret_type: str, path: str, refresh_after: str, destination_name: str, secret_annotations: {str: str} = {} -> VaultStaticSecret {
    vault_static_secret = VaultStaticSecret{
        metadata: meta.ObjectMeta{
            name: secret_name,
        },
        spec: VaultStaticSecretSpec{
            vaultAuthRef: vault_auth_ref,
            mount: mount,
            type: secret_type,
            path: path,
            refreshAfter: refresh_after,
            destination: Destination{
                create: True,
                name: destination_name
                annotations: secret_annotations
            }
        }
    }

    vault_static_secret
}
