package utils

import (
	"encoding/json"
	"fmt"
	"kubernetes/pkg/schema/generator"
	"os"
	"os/exec"
	"path/filepath"
)

func GetServiceMeta(root string, servicePath string) (*generator.GeneratorMeta, error) {

	joinedPath := filepath.Join(root, servicePath)

	if _, err := os.Stat(joinedPath); os.IsNotExist(err) {
		fmt.Println("❌ Generator does not exist at specified location")
		return nil, err
	}

	cmd := exec.Command("go", "run",
		joinedPath,
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
