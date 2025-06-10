import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema Subject:
  organizations: [str]

schema IssuerRef:
  name: str
  kind?: str
  group?: str

schema SecretTemplate:
  annotations?: {str:str}
  labels?: {str:str}

schema PrivateKey:
  algorithm?: "ECDSA" | "RSA" | "Ed25519"
  size?: int = 521

schema CertificateSpec:
  subject?: Subject
  secretName: str
  commonName?: str
  isCA?: bool
  issuerRef: IssuerRef
  dnsNames?: [str]
  secretTemplate?: SecretTemplate
  privateKey?: PrivateKey

schema Certificate:
  apiVersion: "cert-manager.io/v1"
  kind: "Certificate"
  metadata?: meta.ObjectMeta
  spec?: CertificateSpec
    
create_certificate = lambda certificate_name: str, organizations: [str], secret_name: str, issuer_ref_name: str, issuer_ref_kind: str, dns_names: [str], secret_annotations: {str: str} -> Certificate {
  {
    metadata.name = certificate_name
    spec.subject.organizations = organizations
    spec.secretName = secret_name
    spec.issuerRef.name = issuer_ref_name
    spec.issuerRef.kind = issuer_ref_kind
    spec.dnsNames = dns_names
    spec.secretTemplate.annotations = secret_annotations
  }
}