package cli

import (
	"fmt"
	"io/fs"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
)

func getGeneratorLocations(workingDirectory string) ([]string, error) {
	locations := []string{}

	err := filepath.WalkDir(workingDirectory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if d.Name() == "main.go" {
				locations = append(locations, strings.Replace(path, "main.go", "", 1))
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func discoverGeneratorsViaPath(workingDirectory string, logOutputView *tview.TextView) generator.GeneratorMetas {
	logToOutput(logOutputView, "Getting locations")
	locations, err := getGeneratorLocations(workingDirectory)
	if err != nil {
		logToOutput(logOutputView, fmt.Sprintf("Failed to get locations for generators: \n %v", err))
	}

	metas := []generator.GeneratorMeta{}

	logToOutput(logOutputView, "Getting meta for locations")
	for _, location := range locations {
		meta, err := utils.GetGeneratorMeta(location)
		if err != nil || meta == nil {
			logToOutput(logOutputView, fmt.Sprintf("Failed getting meta for location: %v \n Error: %v \n", location, err))
		} else {
			metas = append(metas, *meta)
		}
	}

	return metas
}
