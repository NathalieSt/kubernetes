package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
		return
	}

	name := "csi-driver-nfs"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:          name,
		Namespace:     "csi-driver-nfs",
		GeneratorType: generatorType,
		Helm: &generator.Helm{
			Chart:   "csi-driver-nfs",
			Url:     "https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/charts",
			Version: "4.11.0",
		},
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/csi-driver-nfs/"),
		CreateManifests:  createCSIDriverNFSManifests,
	})
}
