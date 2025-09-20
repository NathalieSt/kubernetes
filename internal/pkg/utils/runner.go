package utils

import (
	"encoding/json"
	"fmt"
	"kubernetes/pkg/schema/generator"
	"os"
)

type GeneratorRunnerConfig struct {
	Meta             generator.GeneratorMeta
	ShouldReturnMeta bool
	OutputDir        string
	CreateManifests  func(generator.GeneratorMeta) map[string][]byte
}

func RunGenerator(config GeneratorRunnerConfig) {
	meta := config.Meta
	if config.ShouldReturnMeta {
		json.NewEncoder(os.Stdout).Encode(meta)
		return
	}

	fmt.Println("✅ Creating K8s manifests")
	manifests := config.CreateManifests(meta)
	if manifests == nil {
		fmt.Println("\nAn error happened while running CreateManifests, can't write yaml to output")
		return
	}

	fmt.Println("✅ Writing K8s manifests to output directory")
	WriteManifestsToYaml(config.OutputDir, manifests)

}
