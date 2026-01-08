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

	name := shared.IntelDevicePluginsOperator
	namespace := "inteldeviceplugins-system"
	generatorType := generator.Infrastructure
	var Vault = generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		Port:          8200,
		Helm: &generator.Helm{
			Chart:   "intel-device-plugins-operator",
			Url:     "https://intel.github.io/helm-charts/",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:      "./cluster/infrastructure/intel-device-plugins/operator",
			Prune:     true,
			Wait:      true,
			Timeout:   "10m",
			DependsOn: []string{},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             Vault,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/intel-device-plugins/operator"),
		CreateManifests:  createVaultManifests,
	})
}
