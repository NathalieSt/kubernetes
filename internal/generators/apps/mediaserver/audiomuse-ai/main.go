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

	name := shared.AudiomuseAI
	namespace := "audiomuse-ai"
	generatorType := generator.App
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		ClusterUrl:    "audiomuse.jellyfin.svc.cluster.local",
		Port:          8000,
		Docker: &generator.Docker{
			Registry: "ghcr.io/neptunehub/audiomuse-ai",
			Version:  utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/apps/mediaserver/audiomuse-ai",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []kustomization.KustomizationDependency{
				{Name: shared.Jellyfin},
				{Name: shared.Redis},
				{Name: shared.MainPostgres},
				{Name: shared.VaultSecretsOperatorConfig},
			},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/apps/mediaserver/audiomuse-ai/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			return createAudiomuseAIManifests(flags.RootDir, meta)
		},
	})
}
