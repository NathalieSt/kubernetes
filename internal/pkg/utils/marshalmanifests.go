package utils

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type ManifestConfig struct {
	Filename  string
	Manifests []any
}

func MarshalManifests(manifestConfigs []ManifestConfig) map[string][]byte {
	result := make(map[string][]byte)
	for _, cfg := range manifestConfigs {

		manifestsBytes := []byte{}

		for _, manifest := range cfg.Manifests {
			data, err := yaml.MarshalWithOptions(manifest, yaml.UseLiteralStyleIfMultiline(true))
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}
			manifestsBytes = append(manifestsBytes, []byte("---\n")...)
			manifestsBytes = append(manifestsBytes, data...)
		}

		result[cfg.Filename] = manifestsBytes
	}
	return result
}
