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

	name := shared.NVIDIAK8sDevicePlugin
	namespace := "nvidia-k8s-device-plugin"
	generatorType := generator.Infrastructure
	var Vault = generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Chart:   "nvidia-device-plugin",
			Url:     "https://nvidia.github.io/k8s-device-plugin",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:      "./cluster/infrastructure/nvidia-k8s-device-plugin",
			Prune:     true,
			Wait:      true,
			Timeout:   "10m",
			DependsOn: []kustomization.KustomizationDependency{},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             Vault,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "./cluster/infrastructure/nvidia-k8s-device-plugin"),
		CreateManifests:  createNVIDIAGPUOperatorManifests,
	})
}
