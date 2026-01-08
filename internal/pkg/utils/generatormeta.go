package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"kubernetes/pkg/schema/generator"
	"maps"
	"os"
	"os/exec"
	"path"
	"slices"
	"sync"
)

type DiscoveredGenerators struct {
	Apps           map[string]string
	Infrastructure map[string]string
	Monitoring     map[string]string
}

func GetDiscoveredGeneratorsMeta(rootDir string) (generator.GeneratorMetas, error) {
	data, err := os.ReadFile(path.Join(rootDir, "clidata/discoveredgenerators.json"))
	if err != nil {
		return nil, err
	}

	var discoveredGenerators DiscoveredGenerators
	if err := json.Unmarshal(data, &discoveredGenerators); err != nil {
		return nil, err
	}

	generatorPaths := []string{}
	appPaths := slices.Collect(maps.Values(discoveredGenerators.Apps))
	infraPaths := slices.Collect(maps.Values(discoveredGenerators.Infrastructure))
	monitoringPaths := slices.Collect(maps.Values(discoveredGenerators.Monitoring))
	generatorPaths = append(generatorPaths, appPaths...)
	generatorPaths = append(generatorPaths, infraPaths...)
	generatorPaths = append(generatorPaths, monitoringPaths...)

	generatorMetas := GetGeneratorMetasByPaths(rootDir, generatorPaths)

	return generatorMetas, nil
}

func RunGeneratorMain(path string, flags []string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("❌ Generator does not exist at specified location")
		return nil, err
	}

	arguments := []string{"run", path}
	arguments = append(arguments, flags...)

	cmd := exec.Command("go", arguments...)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Executing the generator failed")
		fmt.Println(err.Error())
		return nil, err
	}

	return output, nil
}

func GetGeneratorMeta(rootDir string, path string) (*generator.GeneratorMeta, error) {
	rootFlag := fmt.Sprintf("--root=%v", rootDir)
	output, err := RunGeneratorMain(path, []string{
		rootFlag,
		"--metadata",
	})
	if err != nil {
		fmt.Println("❌ Failed to get metadata for generator")
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
			meta, err := GetGeneratorMeta(rootDir, path)
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

func GetGeneratorFlags() *generator.GeneratorFlags {
	rootDir := flag.String("root", "", "The root directory of this project")
	metadataFlag := flag.Bool("metadata", false, "Output generator metadata")

	flag.Parse()

	if *rootDir == "" {
		fmt.Println("❌ Invalid rootDir was passed to generator")
		return nil
	}

	return &generator.GeneratorFlags{
		RootDir:          *rootDir,
		ShouldReturnMeta: *metadataFlag,
	}

}
