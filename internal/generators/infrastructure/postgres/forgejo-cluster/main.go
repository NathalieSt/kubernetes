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

	name := shared.ForgejoPostgres
	namespace := "forgejo-pg-cluster"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		ClusterUrl:    "forgejo-pg-rw.forgejo-pg-cluster.svc.cluster.local",
		Docker: &generator.Docker{
			Registry: "ghcr.io/cloudnative-pg/postgresql",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Port: 5432,
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/infrastructure/postgres/forgejo",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []string{
				shared.CSIDriverNFS,
				shared.VaultSecretsOperatorConfig,
			},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/postgres/forgejo"),
		CreateManifests:  createPostgresManifests,
	})
}
