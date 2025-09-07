package utils

import (
	"encoding/json"
	"fmt"
	"kubernetes/pkg/schema/generator"
	"os/exec"
	"path/filepath"
)

func GetServiceMeta(root string, servicePath string) (*generator.GeneratorMeta, error) {

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
