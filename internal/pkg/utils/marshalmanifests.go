package utils

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type ManifestConfig struct {
	Filename string
	Generate func() any
}

func MarshalManifests(manifestConfigs []ManifestConfig) map[string][]byte {
	result := make(map[string][]byte)
	for _, cfg := range manifestConfigs {
		data, err := yaml.Marshal(cfg.Generate())
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		result[cfg.Filename] = data
	}
	return result
}
