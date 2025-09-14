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

	fmt.Println("✅ Writing K8s manifests to output directory")
	WriteManifestsToYaml(config.OutputDir, manifests)

}
