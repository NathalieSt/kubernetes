package main

import (
	"encoding/json"
	"fmt"
	"io"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"os"
	"os/exec"
	"path/filepath"
)

type ExposedServices map[string]string

func getExposedServices(root string) (*ExposedServices, error) {

	file, err := os.Open(filepath.Join(root, "exposedservices.json"))
	if err != nil {
		fmt.Printf("❌ failed to open exposedservices.json: \n %v", err)
		return nil, err
	}

	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var services ExposedServices

	err = json.Unmarshal(byteValue, &services)
	if err != nil {
		fmt.Printf("❌ error while marhsaling values from exposedservices.json: \n %v", err)
	}

	fmt.Printf("Services from exposedServices.json:  %v", services)

	return &services, nil
}

func getExposedServiceMeta(root string, servicePath string) (*generator.GeneratorMeta, error) {

	joinedPath := filepath.Join(root, servicePath)

	cmd := exec.Command("go", "run",
		joinedPath,
		"--metadata",
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Generator failed: %v \n", err)
		return nil, err
	}

	var generatorMeta generator.GeneratorMeta
	if err := json.Unmarshal(output, &generatorMeta); err != nil {
		fmt.Printf("❌ Failed to parse objects: %v \n", err)
		return nil, err
	}

	return &generatorMeta, nil
}

func getAllExposedServicesMeta(root string, services ExposedServices) []generator.GeneratorMeta {
	allMetas := []generator.GeneratorMeta{}

	for k, v := range services {
		meta, err := getExposedServiceMeta(root, v)
		if err != nil {
			fmt.Printf("Failed to get meta for service: \n %v", k)
			fmt.Printf("Reason: \n %v", err)
		}
		allMetas = append(allMetas, *meta)
	}

	return allMetas
}

func getCaddyConfigMap() {
	root, err := utils.FindRoot()
	if err != nil {
		fmt.Printf("Failed to get root: %v", err)
		return
	}

	exposedServices, err := getExposedServices(root)
	if err != nil {
		fmt.Printf("Failed to get exposed services: %v", err)
	}

	metas := getAllExposedServicesMeta(root, *exposedServices)

	fmt.Printf("found metas: \n %v", metas)
}
