package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"os"
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

type DiscoveredGenerators struct {
	Apps           []string
	Infrastructure []string
	Istio          []string
	Monitoring     []string
}

func discoverGeneratorsByCategoryViaPath(workingDirectory string, rootDir string, logOutputView *tview.TextView) ([]generator.GeneratorMeta, []generator.GeneratorMeta, []generator.GeneratorMeta, []generator.GeneratorMeta) {
	logToOutput(logOutputView, "Getting locations")
	locations, err := getGeneratorLocations(workingDirectory)
	if err != nil {
		logToOutput(logOutputView, fmt.Sprintf("Failed to get locations for generators: \n %v", err))
	}

	logToOutput(logOutputView, "Getting meta for locations")

	appsPaths := []string{}
	infrastructurePaths := []string{}
	istioPaths := []string{}
	monitoringPaths := []string{}

	apps := []generator.GeneratorMeta{}
	infrastructure := []generator.GeneratorMeta{}
	istio := []generator.GeneratorMeta{}
	monitoring := []generator.GeneratorMeta{}

	for _, location := range locations {
		meta, err := utils.GetGeneratorMeta(location)
		if err != nil || meta == nil {
			logToOutput(logOutputView, fmt.Sprintf("Failed getting meta for location: %v \n Error: %v \n", location, err))
		} else {
			switch meta.GeneratorType {
			case generator.App:
				apps = append(apps, *meta)
				appsPaths = append(appsPaths, location)
			case generator.Infrastructure:
				infrastructure = append(infrastructure, *meta)
				infrastructurePaths = append(infrastructurePaths, location)
			case generator.Istio:
				istio = append(istio, *meta)
				istioPaths = append(istioPaths, location)
			case generator.Monitoring:
				monitoring = append(monitoring, *meta)
				monitoringPaths = append(monitoringPaths, location)
			}
		}
	}

	discoveredGenerators := DiscoveredGenerators{
		Apps:           appsPaths,
		Infrastructure: infrastructurePaths,
		Istio:          istioPaths,
		Monitoring:     monitoringPaths,
	}

	discoveredGeneratorsBytes, err := json.Marshal(discoveredGenerators)
	if err != nil {
		logToOutput(logOutputView, fmt.Sprintf("Marshalling discovered generators to json failed\n Reason: %v", err))
	}

	var out bytes.Buffer
	json.Indent(&out, discoveredGeneratorsBytes, "", "")

	err = os.WriteFile(filepath.Join(rootDir, "clidata/discoveredgenerators.json"), out.Bytes(), 0644)
	if err != nil {
		logToOutput(logOutputView, fmt.Sprintf("Error writing discovered generators to json file \n Reason: %v", err))
	}

	return apps, infrastructure, istio, monitoring
}
