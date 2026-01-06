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

	name := "nfs-volumes"
	generatorType := generator.Infrastructure
	meta := generator.GeneratorMeta{
		Name:                name,
		GeneratorType:       generatorType,
		DependsOnGenerators: []string{},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/infrastructure/nfs-volumes/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			return createNFSVolumesManifests(flags.RootDir, gm)
		},
	})
}
