package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"kubernetes/pkg/schema/generator"
	"os"
	"path/filepath"
)

type ExposedServices map[string]string

func GetExposedServices(root string) (*ExposedServices, error) {

	file, err := os.Open(filepath.Join(root, "exposedservices.json"))
	if err != nil {
		fmt.Printf("failed to open exposedservices.json \n")
		return nil, err
	}

	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var services ExposedServices

	err = json.Unmarshal(byteValue, &services)
	if err != nil {
		fmt.Printf("error while marhsaling values from exposedservices.json \n")
		return nil, err
	}

	return &services, nil
}

func GetMetaForExposedServices() ([]generator.GeneratorMeta, error) {

	root, err := FindRoot()
	if err != nil {
		fmt.Printf("Failed to get root")
		return nil, err
	}

	exposedServices, err := GetExposedServices(root)
	if err != nil {
		fmt.Printf("Failed to get exposed services \n")
		return nil, err
	}

	allMetas := []generator.GeneratorMeta{}

	for k, v := range *exposedServices {
		meta, err := GetServiceMeta(root, v)
		if err != nil {
			fmt.Printf("Failed to get meta for service: \n %v \n", k)
			fmt.Printf("Reason: \n %v \n", err)
		} else {
			allMetas = append(allMetas, *meta)
		}
	}

	return allMetas, nil
}
