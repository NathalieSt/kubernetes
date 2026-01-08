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

	name := shared.CSIDriverNFS
	namespace := "csi-driver-nfs"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     namespace,
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Chart:   "csi-driver-nfs",
			Url:     "https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/charts",
			Version: "4.11.0",
		},
		Flux: &kustomization.KustomizationSpec{
			Interval:        "24h",
			TargetNamespace: namespace,
			SourceRef: kustomization.SourceRef{
				Kind: kustomization.GitRepository,
				Name: "flux-system",
			},
			Path:      "./cluster/infrastructure/csi-driver-nfs",
			Prune:     true,
			Wait:      true,
			Timeout:   "10m",
			DependsOn: []kustomization.KustomizationDependency{},
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/csi-driver-nfs/"),
		CreateManifests:  createCSIDriverNFSManifests,
	})
}
