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

	name := shared.Openclarity
	namespace := "openclarity"
	generatorType := generator.Monitoring
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		ClusterUrl:    "openclarity-ui.openclarity.svc.cluster.local",
		Port:          8080,
		Helm: &generator.Helm{
			Url:     "oci://ghcr.io/openclarity/charts/openclarity",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "openclarity",
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/monitoring/openclarity",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []kustomization.KustomizationDependency{
				{Name: shared.MainPostgres},
				{Name: shared.Reflector},
			},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/openclarity/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests := createOpenclarityManifests(flags.RootDir, gm)
			return manifests
		},
	})
}
