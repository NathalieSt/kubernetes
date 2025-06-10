import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema SourceRef:
    kind: "HelmRepository"
    name: str

schema ChartSpec:
    chart: str
    version: "*"
    sourceRef: SourceRef
    interval: "24h"

schema Chart:
    spec: ChartSpec

schema Remediation:
    retries: 3

schema Install:
    remediation: Remediation

schema ValuesFrom:
    kind: "Secret" | "ConfigMap"
    name: str
    valuesKey: str
    targetPath: str
    optional: bool = False

schema ReleaseSpec:
    releaseName: str
    chart: Chart
    interval: "24h"
    install: Install = {
        remediation = {}
    }
    timeout: "5m"
    valuesFrom: [ValuesFrom] = []
    values: any

schema HelmRelease:
    apiVersion: "helm.toolkit.fluxcd.io/v2"
    kind: "HelmRelease"
    metadata: meta.ObjectMeta
    spec: ReleaseSpec

create_helm_release = lambda values: any, metadata_name: str, release_name: str, chart_name: str, source_ref_name: str, valuesFrom: [ValuesFrom] -> HelmRelease {
    {
        metadata.name = metadata_name
        spec.releaseName = release_name
        spec.chart.spec.chart = chart_name
        spec.chart.spec.sourceRef.name = source_ref_name
        spec.values = values
        spec.valuesFrom = valuesFrom
    }
}