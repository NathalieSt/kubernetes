package main

import (
	"fmt"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"os"
	"path"
)

func getDiscordBridgeConfigMap(rootDir string, relativeDir string, name string) (*core.ConfigMap, error) {
	config, err := os.ReadFile(path.Join(rootDir, relativeDir, "config.yaml"))
	if err != nil {
		fmt.Printf("Error while reading config.yaml")
		return nil, err
	}

	configMap := core.NewConfigMap(meta.ObjectMeta{
		Name: name,
	}, map[string]string{
		"config.yaml": string(config),
	})

	return &configMap, nil
}
