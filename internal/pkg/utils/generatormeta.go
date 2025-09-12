package utils

import (
	"encoding/json"
	"fmt"
	"kubernetes/pkg/schema/generator"
	"os"
	"os/exec"
)

func GetGeneratorMeta(path string) (*generator.GeneratorMeta, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("❌ Generator does not exist at specified location")
		return nil, err
	}

	cmd := exec.Command("go", "run",
		path,
		"--metadata",
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ The command for getting the meta of the generator failed")
		return nil, err
	}

	var generatorMeta generator.GeneratorMeta
	if err := json.Unmarshal(output, &generatorMeta); err != nil {
		fmt.Println("❌ Failed to parse meta for generator")
		return nil, err
	}

	return &generatorMeta, nil
}

func GetGeneratorMetasByPaths(paths []string) []generator.GeneratorMeta {
	generators := []generator.GeneratorMeta{}
	for _, path := range paths {
		meta, err := GetGeneratorMeta(path)
		if err != nil {
			fmt.Printf("Error while getting generator for path: %v \n Reason: \n %v", path, err)
		} else {
			generators = append(generators, *meta)
		}
	}
	return generators
}
