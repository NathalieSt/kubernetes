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
	DNSName                    string
	HeaderForwardingIsRequired bool
	WebsocketSupportIsRequired bool
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

type GeneratorFlags struct {
	RootDir          string
	ShouldReturnMeta bool
}

type GeneratorMeta struct {
	Name                string
	Namespace           string
	GeneratorType       GeneratorType
	ClusterUrl          string
	Port                int64
	Docker              *Docker
	Helm                *Helm
	Caddy               *Caddy
	VirtualService      *VirtualServiceConfig
	KedaScaling         *keda.ScaledObjectTriggerMeta
	DependsOnGenerators []string
}

type GeneratorMetas []GeneratorMeta

func (metas GeneratorMetas) GetDNSNames() []string {
	var list []string
	for _, meta := range metas {
		list = append(list, meta.Caddy.DNSName)
	}
	return list
}

func (metas GeneratorMetas) GetMetasSeparatedByCategories() ([]GeneratorMeta, []GeneratorMeta, []GeneratorMeta, []GeneratorMeta) {
	apps := []GeneratorMeta{}
	infrastructure := []GeneratorMeta{}
	istio := []GeneratorMeta{}
	monitoring := []GeneratorMeta{}

	for _, meta := range metas {
		switch meta.GeneratorType {
		case App:
			apps = append(apps, meta)
		case Infrastructure:
			infrastructure = append(infrastructure, meta)
		case Istio:
			istio = append(istio, meta)
		case Monitoring:
			monitoring = append(monitoring, meta)
		}
	}
	return apps, infrastructure, istio, monitoring
}
