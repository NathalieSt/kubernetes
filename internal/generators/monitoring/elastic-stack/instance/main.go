package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/kustomization"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
		return
	}

	name := shared.ElasticStackInstance
	namespace := "elastic-stack"
	generatorType := generator.Monitoring
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Url:     "https://helm.elastic.co",
			Chart:   "eck-stack",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		ClusterUrl: "elastic-stack-eck-kibana-kb-http.elastic-stack.svc.cluster.local",
		Port:       5601,
		Caddy: &generator.Caddy{
			DNSName: "kibana",
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/monitoring/elastic-stack/instance",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []string{
				shared.CSIDriverNFS,
				shared.Vector,
			},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/elastic-stack/instance"),
		CreateManifests:  createElasticStackManifests,
	})
}
