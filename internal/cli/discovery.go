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

func WriteToJSONFile(destination string, toWrite any, logOutputView *tview.TextView) {
	marshalledBytes, err := json.Marshal(toWrite)
	if err != nil {
		logToOutput(logOutputView, fmt.Sprintf("Marshalling discovered generators to json failed\n Reason: %v", err))
	}

	var out bytes.Buffer
	json.Indent(&out, marshalledBytes, "", "")

	err = os.WriteFile(destination, out.Bytes(), 0644)
	if err != nil {
		logToOutput(logOutputView, fmt.Sprintf("Error writing discovered generators to json file \n Reason: %v", err))
	}
}

type DiscoveredGenerators struct {
	Apps           map[string]string
	Infrastructure map[string]string
	Istio          map[string]string
	Monitoring     map[string]string
}

type ExposedGenerators map[string]string

func discoverGeneratorsByCategoryViaPath(workingDirectory string, rootDir string, logOutputView *tview.TextView) ([]generator.GeneratorMeta, []generator.GeneratorMeta, []generator.GeneratorMeta, []generator.GeneratorMeta) {
	logToOutput(logOutputView, "Getting locations")
	locations, err := getGeneratorLocations(workingDirectory)
	if err != nil {
		logToOutput(logOutputView, fmt.Sprintf("Failed to get locations for generators: \n %v", err))
	}

	logToOutput(logOutputView, "Getting meta for locations")

	appsPaths := map[string]string{}
	infrastructurePaths := map[string]string{}
	istioPaths := map[string]string{}
	monitoringPaths := map[string]string{}

	apps := []generator.GeneratorMeta{}
	infrastructure := []generator.GeneratorMeta{}
	istio := []generator.GeneratorMeta{}
	monitoring := []generator.GeneratorMeta{}

	exposedGenerators := ExposedGenerators{}

	for _, location := range locations {
		meta, err := utils.GetGeneratorMeta(location)
		if err != nil || meta == nil {
			logToOutput(logOutputView, fmt.Sprintf("Failed getting meta for location: %v \n Error: %v \n", location, err))
		} else {
			if meta.Caddy != nil {
				exposedGenerators[meta.Name] = location
			}
			switch meta.GeneratorType {
			case generator.App:
				apps = append(apps, *meta)
				appsPaths[meta.Name] = location
			case generator.Infrastructure:
				infrastructure = append(infrastructure, *meta)
				infrastructurePaths[meta.Name] = location
			case generator.Istio:
				istio = append(istio, *meta)
				istioPaths[meta.Name] = location
			case generator.Monitoring:
				monitoring = append(monitoring, *meta)
				monitoringPaths[meta.Name] = location
			}
		}
	}

	discoveredGenerators := DiscoveredGenerators{
		Apps:           appsPaths,
		Infrastructure: infrastructurePaths,
		Istio:          istioPaths,
		Monitoring:     monitoringPaths,
	}

	WriteToJSONFile(filepath.Join(rootDir, "clidata/discoveredgenerators.json"), discoveredGenerators, logOutputView)
	WriteToJSONFile(filepath.Join(rootDir, "clidata/exposedgenerators.json"), exposedGenerators, logOutputView)

	return apps, infrastructure, istio, monitoring
}
