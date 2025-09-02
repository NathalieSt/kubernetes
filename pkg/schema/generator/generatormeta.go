package generator

/*
schema ClusterService:
    name: str
    namespace: str
    type: "app" | "infrastructure" | "istio" | "monitoring"
    cluster_url?: str
    port?: int
    gateway_config?: GatewayConfig
    virtual_service_config?: VirtualServiceConfig
    flux_kustomization: ClusterServiceFluxKustomization
    keda_scaling?: keda.ScaledObjectTriggerMeta
*/
import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
)

type VirtualServiceConfig struct{}

type Caddy struct {
}

type Docker struct {
	Registry string
	Version  string
}

type Helm struct {
	Url     string
	Chart   string
	Version string
}

type GeneratorType = int

const (
	App GeneratorType = iota
	Infrastructure
	Istio
	Monitoring
)

type GeneratorMeta struct {
	Name                string
	Namespace           string
	EntityType          GeneratorType
	ClusterUrl          string
	Port                int
	Docker              Docker
	Helm                Helm
	Caddy               Caddy
	VirtualService      VirtualServiceConfig
	KedaScaling         keda.ScaledObjectTriggerMeta
	DependsOnGenerators []string
}
