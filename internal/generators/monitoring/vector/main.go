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

	name := shared.Vector
	namespace := "vector"
	generatorType := generator.Monitoring
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		Docker: &generator.Docker{
			Registry: "timberio/vector",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/monitoring/vector",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []string{
				shared.CSIDriverNFS,
			},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/vector/"),
		CreateManifests:  createVectorManifests,
	})
}
