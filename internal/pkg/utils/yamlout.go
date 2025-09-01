package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func WriteManifestsToYaml(path string, files map[string][]byte) {
	for key, value := range files {
		err := os.WriteFile(filepath.Join(path, key), value, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			log.Fatal(err)
		}
	}

}
