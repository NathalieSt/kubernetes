package main

import (
	"fmt"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"os"
)

func getIPLeakConfigMap(name string) (*core.ConfigMap, error) {

	//FIXME: rootDir and relative dir need to be prepended
	config, err := os.ReadFile("./ipleak.sh")
	if err != nil {
		fmt.Printf("Error while reading ipleak.sh")
		return nil, err
	}

	configMap := core.NewConfigMap(meta.ObjectMeta{
		Name: name,
	}, map[string]string{
		"ipleak.sh": string(config),
	})

	return &configMap, nil
}
