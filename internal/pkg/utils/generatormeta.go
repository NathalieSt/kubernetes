package utils

import (
	"encoding/json"
	"fmt"
	"kubernetes/pkg/schema/generator"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func RunGeneratorMain(path string, flags []string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("❌ Generator does not exist at specified location")
		return nil, err
	}

	cmd := exec.Command("go", "run",
		path,
		strings.Join(flags, " "),
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ The command for getting the meta of the generator failed")
		return nil, err
	}

	return output, nil
}

func GetGeneratorMeta(path string) (*generator.GeneratorMeta, error) {

	output, err := RunGeneratorMain(path, []string{"--metadata"})
	if err != nil {
		fmt.Println("❌ Failed to run generator with metadata flag")
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

	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Go(func() {
			meta, err := GetGeneratorMeta(path)
			if err != nil {
				fmt.Printf("❌ Error while getting generator for path: %v \n Reason: \n %v", path, err)
			} else {
				generators = append(generators, *meta)
			}
		})
	}
	wg.Wait()
	return generators
}
