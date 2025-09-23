package main

import (
	"fmt"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"os"
)

func getDiscordBridgeConfigMap(name string) (*core.ConfigMap, error) {

	//FIXME: rootDir and relative dir need to be prepended
	config, err := os.ReadFile("./config.yaml")
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
