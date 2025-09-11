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

func GetExposedGenerators(root string) (*ExposedServices, error) {

	file, err := os.Open(filepath.Join(root, "clidata/exposedgenerators.json"))
	if err != nil {
		fmt.Printf("failed to open exposedgenerators.json \n")
		return nil, err
	}

	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var services ExposedServices

	err = json.Unmarshal(byteValue, &services)
	if err != nil {
		fmt.Printf("error while marhsaling values from exposedgenerators.json \n")
		return nil, err
	}

	return &services, nil
}

func GetMetaForExposedGenerators() (generator.GeneratorMetas, error) {
	root, err := FindRoot()
	if err != nil {
		fmt.Printf("Failed to get root")
		return nil, err
	}

	exposedGenerators, err := GetExposedGenerators(root)
	if err != nil {
		fmt.Printf("Failed to get exposed generators \n")
		return nil, err
	}

	allMetas := []generator.GeneratorMeta{}

	for k, v := range *exposedGenerators {
		joinedPath := filepath.Join(root, v)
		meta, err := GetGeneratorMeta(joinedPath)
		if err != nil {
			fmt.Printf("Failed to get meta for generator: \n %v \n", k)
			fmt.Printf("Reason: \n %v \n", err)
		} else {
			allMetas = append(allMetas, *meta)
		}
	}

	return allMetas, nil
}
