package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"kubernetes/pkg/schema/generator"
	"os"
)

type GeneratorConfig struct {
	Meta            generator.GeneratorMeta
	OutputDir       string
	CreateManifests func(generator.GeneratorMeta) map[string][]byte
}

func RunGenerator(config GeneratorConfig) {

	meta := config.Meta
	metadataFlag := flag.Bool("metadata", false, "Output generator metadata")

	flag.Parse()

	if *metadataFlag {
		fmt.Println("✅ Encoding Meta in JSON")
		json.NewEncoder(os.Stdout).Encode(meta)
		return
	}

	fmt.Println("✅ Creating K8s manifests")
	manifests := config.CreateManifests(meta)

	fmt.Println("✅ Writing K8s manifests to output directory")
	WriteManifestsToYaml(config.OutputDir, manifests)

}
