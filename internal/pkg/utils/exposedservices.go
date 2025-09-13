package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"kubernetes/pkg/schema/generator"
	"os"
	"path/filepath"
)

type ExposedGenerators map[string]string

func GetExposedGenerators(rootDir string) (*ExposedGenerators, error) {

	file, err := os.Open(filepath.Join(rootDir, "clidata/exposedgenerators.json"))
	if err != nil {
		fmt.Printf("failed to open exposedgenerators.json \n")
		return nil, err
	}

	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var generators ExposedGenerators

	err = json.Unmarshal(byteValue, &generators)
	if err != nil {
		fmt.Printf("error while marhsaling values from exposedgenerators.json \n")
		return nil, err
	}

	return &generators, nil
}

func GetMetaForExposedGenerators(rootDir string) (generator.GeneratorMetas, error) {
	exposedGenerators, err := GetExposedGenerators(rootDir)
	if err != nil {
		fmt.Printf("Failed to get exposed generators \n")
		return nil, err
	}

	allMetas := []generator.GeneratorMeta{}

	for name, location := range *exposedGenerators {
		meta, err := GetGeneratorMeta(rootDir, location)
		if err != nil {
			fmt.Printf("Failed to get meta for generator: \n %v \n", name)
			fmt.Printf("Reason: \n %v \n", err)
		} else {
			allMetas = append(allMetas, *meta)
		}
	}

	return allMetas, nil
}
