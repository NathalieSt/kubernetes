package cluster

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
	"kubernetes/pkg/schema/cluster/monitoring/keda"
)

type VirtualServiceConfig struct{}

type GatewayConfig struct {
}

type EntityType = int

const (
	App EntityType = iota
	Infrastructure
	Istio
	Monitoring
)

type Entity struct {
	Name                 string
	Namespace            string
	EntityType           EntityType
	ClusterUrl           string
	Port                 int
	GatewayConfig        *GatewayConfig
	VirtualServiceConfig *VirtualServiceConfig
	KedaScaling          *keda.ScaledObject
}
