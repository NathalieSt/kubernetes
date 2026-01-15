package main

import (
	"fmt"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"os"
	"path"
)

func getSynapseConfigMap(name string, rootDir string, relativeDir string) (*core.ConfigMap, error) {
	homeserver, err := os.ReadFile(path.Join(rootDir, relativeDir, "homeserver.yaml"))
	if err != nil {
		fmt.Printf("Error while reading homeserver.yaml")
		return nil, err
	}

	logConfig, err := os.ReadFile(path.Join(rootDir, relativeDir, "log.config"))
	if err != nil {
		fmt.Printf("Error while reading config")
		return nil, err
	}

	discordRegistration, err := os.ReadFile(path.Join(rootDir, relativeDir, "discord-registration.yaml"))
	if err != nil {
		fmt.Printf("Error while reading discord-registration.yaml")
		return nil, err
	}

	whatsappRegistration, err := os.ReadFile(path.Join(rootDir, relativeDir, "whatsapp-registration.yaml"))
	if err != nil {
		fmt.Printf("Error while reading discord-registration.yaml")
		return nil, err
	}

	configMap := core.NewConfigMap(meta.ObjectMeta{
		Name: name,
	}, map[string]string{
		"homeserver.yaml": string(homeserver),
		"matrix.cloud.nathalie-stiefsohn.eu.log.config": string(logConfig),
		"discord-registration.yaml":                     string(discordRegistration),
		"whatsapp-registration.yaml":                    string(whatsappRegistration),
	})

	return &configMap, nil
}
