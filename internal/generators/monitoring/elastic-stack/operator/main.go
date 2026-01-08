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

	name := shared.ElasticStackOperator
	namespace := "elastic-stack-operator"
	generatorType := generator.Monitoring
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Url:     "https://helm.elastic.co",
			Chart:   "eck-operator",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:      "./cluster/monitoring/elastic-stack/operator",
			Prune:     true,
			Wait:      true,
			Timeout:   "10m",
			DependsOn: []string{},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/elastic-stack/operator"),
		CreateManifests:  createElasticStackManifests,
	})
}
