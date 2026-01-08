package generator

import (
	"kubernetes/pkg/schema/cluster/flux/kustomization"
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
	Monitoring
)

type GeneratorFlags struct {
	RootDir          string
	ShouldReturnMeta bool
}

type GeneratorNFSVolume struct {
	Name         string
	StorageClass string
	Path         string
	Capacity     string
}

type GeneratorMeta struct {
	Name          string
	Namespace     string
	GeneratorType GeneratorType
	ClusterUrl    string
	Port          int64
	Flux          *kustomization.KustomizationSpec
	Docker        *Docker
	Helm          *Helm
	Caddy         *Caddy
	KedaScaling   *keda.ScaledObjectTriggerMeta
	NFSVolumes    map[string]GeneratorNFSVolume
}

type GeneratorMetas []GeneratorMeta

func (metas GeneratorMetas) GetDNSNames() []string {
	var list []string
	for _, meta := range metas {
		list = append(list, meta.Caddy.DNSName)
	}
	return list
}

func (metas GeneratorMetas) GetMetasSeparatedByCategories() ([]GeneratorMeta, []GeneratorMeta, []GeneratorMeta) {
	apps := []GeneratorMeta{}
	infrastructure := []GeneratorMeta{}
	monitoring := []GeneratorMeta{}

	for _, meta := range metas {
		switch meta.GeneratorType {
		case App:
			apps = append(apps, meta)
		case Infrastructure:
			infrastructure = append(infrastructure, meta)
		case Monitoring:
			monitoring = append(monitoring, meta)
		}
	}
	return apps, infrastructure, monitoring
}
