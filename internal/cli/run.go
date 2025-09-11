package cli

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"

	"github.com/rivo/tview"
)

func getGeneratorsFromJSON() []generator.GeneratorMeta {
	return []generator.GeneratorMeta{
		{
			Name:          "mealie",
			Namespace:     "mealie",
			GeneratorType: generator.App,
			ClusterUrl:    "mealie.mealie.svc.cluster.local",
			Port:          9000,
			Docker: generator.Docker{
				Registry: "ghcr.io/mealie-recipes/mealie",
				//FIXME: set to nil, later fetch in generator from version.json
				Version: "v3.0.2",
			},
			Caddy: generator.Caddy{
				DNSName: "mealie.cluster",
			},
			KedaScaling: keda.ScaledObjectTriggerMeta{
				Timezone:        "Europe/Vienna",
				Start:           "0 9 * * *",
				End:             "0 21 * * *",
				DesiredReplicas: "1",
			},
			DependsOnGenerators: []string{
				"postgres",
			},
		}, {
			Name:          "searxng",
			Namespace:     "searxng",
			GeneratorType: generator.App,
			ClusterUrl:    "searxng.searxng.svc.cluster.local",
			Port:          8080,
			Docker: generator.Docker{
				Registry: "searxng/searxng",
				//FIXME: set to nil, later fetch in generator from version.json
				Version: "2025.8.3-2e62eb5",
			},
			Caddy: generator.Caddy{
				DNSName: "searxng.cluster",
			},
			KedaScaling: keda.ScaledObjectTriggerMeta{
				Timezone:        "Europe/Vienna",
				Start:           "0 7 * * *",
				End:             "0 23 * * *",
				DesiredReplicas: "1",
			},
			DependsOnGenerators: []string{
				"valkey",
				"gluetun-proxy",
			},
		},
	}
}

func loadGeneratorsToList(list *tview.List) {
	list.Clear()
	generators := getGeneratorsFromJSON()

	for _, generator := range generators {
		version := ""

		if generator.Helm.Version != "" {
			version = generator.Helm.Version
		} else if generator.Docker.Version != "" {
			version = generator.Docker.Version
		} else {
			return
		}

		list.AddItem(fmt.Sprintf("%v (%v)", generator.Name, version), "Test Description", '>', nil)
	}
}

func logToOutput(outputView *tview.TextView, message string) {
	outputView.Write([]byte(message))
	outputView.Write([]byte("\n"))
}

func clearOutput(outputView *tview.TextView) {
	outputView.Clear()
}

func Start() {

	app := tview.NewApplication()

	root, err := utils.FindRoot()
	if err != nil {
		fmt.Println("Error while finding root, reason:")
		fmt.Print(err)
		return
	}

	outputView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	outputView.SetTitle("Command output")
	outputView.SetBorder(true)
	outputView.SetScrollable(true)

	generatorsList := tview.NewList()
	generatorsList.SetBorder(true)
	generatorsList.SetTitle("Discovered generators")
	loadGeneratorsToList(generatorsList)

	commandList := tview.NewList().
		AddItem("Discover", "Discover new generators", 'd', func() {
			clearOutput(outputView)
			logToOutput(outputView, "Starting discovery process...")
			discoverGenerators(fmt.Sprintf("%v/internal/generators", root), outputView)
			logToOutput(outputView, "Reloading discovered generators")
			loadGeneratorsToList(generatorsList)
			logToOutput(outputView, "Discovery completed")
		}).
		AddItem("Run", "Run a generator", 'r', nil).
		AddItem("Version", "Upgrade Helm/Docker versions", 'v', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	commandList.SetTitle("Commands")
	commandList.SetBorder(true)

	actionArea := tview.NewFlex().SetDirection(tview.FlexRow)
	actionArea.AddItem(outputView, 0, 3, false)

	flex := tview.NewFlex().
		AddItem(commandList, 35, 1, true).
		AddItem(actionArea, 0, 3, false).
		AddItem(generatorsList, 0, 1, false)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
