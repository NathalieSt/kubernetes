import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema Template:
    text: str

schema Transformation:
    excludes?: [str]
    templates: {str:Template}

schema Destination:
    create: bool
    name: str
    labels: {str:str} = {}
    annotations: {str:str} = {}
    transformation?: Transformation

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

