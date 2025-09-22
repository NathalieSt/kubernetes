package main

import (
	"fmt"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"os"
)

func getSynapseConfigMap(name string) (*core.ConfigMap, error) {

	//FIXME: rootDir and relative dir need to be prepended
	homeserver, err := os.ReadFile("homeserver.yaml")
	if err != nil {
		fmt.Printf("Error while reading homeserver.yaml")
		return nil, err
	}

	//FIXME: rootDir and relative dir need to be prepended
	logConfig, err := os.ReadFile("log.config")
	if err != nil {
		fmt.Printf("Error while reading config")
		return nil, err
	}

	//FIXME: rootDir and relative dir need to be prepended
	discordRegistration, err := os.ReadFile("discord-registration.yaml")
	if err != nil {
		fmt.Printf("Error while reading discord-registration.yaml")
		return nil, err
	}

	configMap := core.NewConfigMap(meta.ObjectMeta{
		Name: name,
	}, map[string]string{
		"homeserver.yaml": string(homeserver),
		"matrix.cluster.netbird.selfhosted.log.config": string(logConfig),
		"discord-registration.yaml":                    string(discordRegistration),
	})

	return &configMap, nil
}
