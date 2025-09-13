package utils

import (
	"encoding/json"
	"fmt"
	"kubernetes/pkg/schema/generator"
	"os"
	"os/exec"
	"path"
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

func GetGeneratorMeta(rootDir string, path string) (*generator.GeneratorMeta, error) {

	output, err := RunGeneratorMain(path, []string{
		fmt.Sprintf("--root %v", rootDir),
		"--metadata",
	})
	if err != nil {
		fmt.Println("❌ Failed to run generator with metadata flag")
		return nil, err
	}

	var generatorMeta generator.GeneratorMeta
	if err := json.Unmarshal(output, &generatorMeta); err != nil {
		fmt.Println("❌ Failed to parse meta for generator")
		fmt.Printf("Offending meta: %v\n", string(output))
		return nil, err
	}

	return &generatorMeta, nil
}

func GetGeneratorMetasByPaths(rootDir string, paths []string) []generator.GeneratorMeta {
	generators := []generator.GeneratorMeta{}

	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Go(func() {
			meta, err := GetGeneratorMeta(path, rootDir)
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

type Versions map[string]string

func GetGeneratorVersionFromFile(fileLocation string, name string) (string, error) {
	versionsBytes, err := os.ReadFile(fileLocation)
	if err != nil {
		fmt.Printf("❌ Error while getting version for generator : %v", name)
		return "", err
	}

	versions := Versions{}

	err = json.Unmarshal(versionsBytes, &versions)
	if err != nil {
		fmt.Printf("❌ Error while marshalling versions json for generator : %v", name)
		return "", err
	}

	return versions[name], nil
}

func GetGeneratorVersionByType(rootDir string, name string, generatorType generator.GeneratorType) string {

	fileName := ""
	switch generatorType {
	case generator.App:
		fileName = "apps.json"
	case generator.Infrastructure:
		fileName = "infrastructure.json"
	case generator.Istio:
		fileName = "istio.json"
	case generator.Monitoring:
		fileName = "monitoring.json"
	}

	version, err := GetGeneratorVersionFromFile(path.Join(rootDir, "versions", fileName), name)
	if err != nil {
		fmt.Printf("❌ An error occurred while getting the version for generator: %v \n", name)
		fmt.Println("Error: " + err.Error())
	}

	return version
}
