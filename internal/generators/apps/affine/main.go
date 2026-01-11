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

	name := shared.Affine
	namespace := "affine"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		ClusterUrl:    "affine.affine.svc.cluster.local",
		Port:          3010,
		Caddy: &generator.Caddy{
			DNSName: "affine",
		},
		Docker: &generator.Docker{
			Registry: "ghcr.io/toeverything/affine",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/apps/affine",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []kustomization.KustomizationDependency{
				{Name: shared.CSIDriverNFS},
				{Name: shared.Reflector},
				{Name: shared.AffinePostgres},
			},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/affine/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			return createAffineManifests(flags.RootDir, gm)
		},
	})
}
