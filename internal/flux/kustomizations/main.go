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

	name := "flux-kustomizations"
	meta := generator.GeneratorMeta{
		Name:      name,
		Namespace: "flux-system",
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/flux/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			return createFluxKustomizationManifests(flags.RootDir)
		},
	})
}
