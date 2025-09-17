package main

import (
	"fmt"
	"kubernetes/internal/generators/istio"
	"kubernetes/internal/pkg/utils"
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

	name := "kiali"
	generatorType := generator.Monitoring
	kiali := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "kiali-operator",
		GeneratorType: generatorType,
		ClusterUrl:    fmt.Sprintf("kiali.%v.svc.cluster.local", istio.Namespace),
		Port:          20001,
		Caddy: &generator.Caddy{
			DNSName: "kiali.cluster",
		},
		Helm: &generator.Helm{
			Chart:   "kiali-operator",
			Url:     "https://kiali.org/helm-charts",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "Europe/Vienna",
			Start:           "0 9 * * *",
			End:             "0 21 * * *",
			DesiredReplicas: "1",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             kiali,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/monitoring/kiali/"),
		CreateManifests:  createKialiManifests,
	})
}
