import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema SourceRef:
    kind: "HelmRepository" | "OCIRepository"
    name: str

schema ChartSpec:
    chart: str
    version: str
    sourceRef: SourceRef
    interval: str

schema Chart:
    spec: ChartSpec

schema Remediation:
    retries: int

schema Install:
    remediation: Remediation

schema ValuesFrom:
    kind: "Secret" | "ConfigMap"
    name: str
    valuesKey: str
    targetPath: str
    optional: bool = False

schema ChartRef:
    kind: "HelmChart" | "OCIRepository"
    name: str

schema ReleaseSpec:
    releaseName: str
    chart?: Chart
    chartRef?: ChartRef
    interval: str
    install: Install = {
        remediation = {}
    }
    timeout: str
    valuesFrom?: [ValuesFrom]
    values: any

schema HelmRelease:
    apiVersion: "helm.toolkit.fluxcd.io/v2"
    kind: "HelmRelease"
    metadata: meta.ObjectMeta
    spec: ReleaseSpec