import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema ScaledObjectTriggerMeta:
    timezone: str
    start: str
    end: str
    desiredReplicas: str

schema ScaledObjectTrigger:
    type: "cron"
    metadata: ScaledObjectTriggerMeta

schema ScaleTargetRef:
    apiVersion?: str
    kind?: str
    name: str

schema ScaledObjectSpec:
    scaleTargetRef: ScaleTargetRef
    minReplicaCount: int
    cooldownPeriod: int
    triggers: [ScaledObjectTrigger]

schema ScaledObject:
    apiVersion = "keda.sh/v1alpha1"
    kind = "ScaledObject"
    metadata: meta.ObjectMeta
    spec: ScaledObjectSpec

