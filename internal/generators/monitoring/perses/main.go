package main

import (
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/kustomization"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
		return
	}

	name := shared.Perses
	namespace := "perses"
	generatorType := generator.Monitoring
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "perses",
		GeneratorType: generatorType,
		ClusterUrl:    "perses.perses.svc.cluster.local",
		Port:          8080,
		Helm: &generator.Helm{
			Url:     "https://perses.github.io/helm-charts",
			Chart:   "perses",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Caddy: &generator.Caddy{
			DNSName: "perses",
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 8 * * *",
			End:             "0 22 * * *",
			DesiredReplicas: "1",
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:    "./cluster/monitoring/perses",
			Prune:   true,
			Wait:    true,
			Timeout: "10m",
			DependsOn: []string{
				shared.VictoriaMetrics,
				shared.CSIDriverNFS,
			},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/perses/"),
		CreateManifests:  createPersesManifests,
	})
}
