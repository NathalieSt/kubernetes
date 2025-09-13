package cli

import (
	"encoding/json"
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"maps"
	"os"
	"path"
	"slices"
	"sync"

	"github.com/rivo/tview"
)

type GeneratorCommand string

var (
	Run     GeneratorCommand = "Run"
	Version GeneratorCommand = "Version"
	None    GeneratorCommand = "None"
)

func makeNodeExpandable(node *tview.TreeNode) {
	node.SetSelectedFunc(func() {
		node.SetExpanded(!node.IsExpanded())
	})
}

func appendGeneratorsToTreeNode(treeNode *tview.TreeNode, generators []generator.GeneratorMeta, generatorCommandHandler func(generator.GeneratorMeta)) {

	if len(generators) == 0 {
		treeNode.AddChild(tview.NewTreeNode("Generators").SetText("No generators available\nPlease try to run the \"discover\" command to find generators"))
		return
	}

	for _, generator := range generators {
		version := ""

		if generator.Helm != nil && generator.Helm.Version != "" {
			version = fmt.Sprintf("(%v)", generator.Helm.Version)
		} else if generator.Docker != nil && generator.Docker.Version != "" {
			version = fmt.Sprintf("(%v)", generator.Docker.Version)
		}

		generatorNode := tview.NewTreeNode(fmt.Sprintf("%v %v", generator.Name, version))

		generatorNode.SetSelectedFunc(func() {
			generatorCommandHandler(generator)
		})

		treeNode.AddChild(generatorNode)
	}
}

func appendGeneratorsToTree(
	tree *tview.TreeView,
	appsGenerators []generator.GeneratorMeta,
	infrastructureGenerators []generator.GeneratorMeta,
	istioGenerators []generator.GeneratorMeta,
	monitoringGenerators []generator.GeneratorMeta,
	generatorCommandHandler func(generator.GeneratorMeta),
) {
	rootNode := tree.GetRoot()
	rootNode.ClearChildren()

	applicationNode := tview.NewTreeNode("Applications").SetExpanded(false)
	makeNodeExpandable(applicationNode)

	infrastructureNode := tview.NewTreeNode("Infrastructure").SetExpanded(false)
	makeNodeExpandable(infrastructureNode)

	istioNode := tview.NewTreeNode("Istio").SetExpanded(false)
	makeNodeExpandable(istioNode)

	monitoringNode := tview.NewTreeNode("Monitoring").SetExpanded(false)
	makeNodeExpandable(monitoringNode)

	rootNode.AddChild(applicationNode).AddChild(infrastructureNode).AddChild(istioNode).AddChild(monitoringNode)

	appendGeneratorsToTreeNode(applicationNode, appsGenerators, generatorCommandHandler)
	appendGeneratorsToTreeNode(infrastructureNode, infrastructureGenerators, generatorCommandHandler)
	appendGeneratorsToTreeNode(istioNode, istioGenerators, generatorCommandHandler)
	appendGeneratorsToTreeNode(monitoringNode, monitoringGenerators, generatorCommandHandler)
}

func initializeGeneratorTree(rootDir string, outputView *tview.TextView, generatorsTree *tview.TreeView, generatorCommandHandler func(generator.GeneratorMeta)) {
	discoveredGeneratorsBytes, err := os.ReadFile(path.Join(rootDir, "clidata/discoveredgenerators.json"))
	if err != nil {
		logToOutput(outputView, fmt.Sprintf("An error happened while reading discoveredgenerators.json: \n %v", err))
	}

	discoveredGenerators := DiscoveredGenerators{}
	err = json.Unmarshal(discoveredGeneratorsBytes, &discoveredGenerators)
	if err != nil {
		logToOutput(outputView, fmt.Sprintf("An error happened while unmarshalling discoveredgenerators.json: \n %v", err))
	} else {

		appsGenerators := []generator.GeneratorMeta{}
		infrastructureGenerators := []generator.GeneratorMeta{}
		istioGenerators := []generator.GeneratorMeta{}
		monitoringGenerators := []generator.GeneratorMeta{}

		var wg sync.WaitGroup
		wg.Go(func() {
			appsGenerators = utils.GetGeneratorMetasByPaths(rootDir, slices.Collect(maps.Values(discoveredGenerators.Apps)))
		})
		wg.Go(func() {
			infrastructureGenerators = utils.GetGeneratorMetasByPaths(rootDir, slices.Collect(maps.Values(discoveredGenerators.Infrastructure)))
		})
		wg.Go(func() {
			istioGenerators = utils.GetGeneratorMetasByPaths(rootDir, slices.Collect(maps.Values(discoveredGenerators.Istio)))
		})
		wg.Go(func() {
			monitoringGenerators = utils.GetGeneratorMetasByPaths(rootDir, slices.Collect(maps.Values(discoveredGenerators.Monitoring)))
		})

		wg.Wait()

		appendGeneratorsToTree(
			generatorsTree,
			appsGenerators,
			infrastructureGenerators,
			istioGenerators,
			monitoringGenerators,
			generatorCommandHandler,
		)
	}
}

func runGeneratorFromJSON(rootDir string, meta generator.GeneratorMeta, outputView *tview.TextView) {
	logToOutput(outputView, "Getting location for generator")
	discoveredGeneratorsBytes, err := os.ReadFile(path.Join(rootDir, "clidata/discoveredgenerators.json"))
	if err != nil {
		logToOutput(outputView, fmt.Sprintf("An error happened while reading discoveredgenerators.json: \n %v", err))
	}

	discoveredGenerators := DiscoveredGenerators{}
	err = json.Unmarshal(discoveredGeneratorsBytes, &discoveredGenerators)
	if err != nil {
		logToOutput(outputView, fmt.Sprintf("Could not unmarshall discovered generators: \n %v", err))
	}

	generatorLocation := ""

	switch meta.GeneratorType {
	case generator.App:
		generatorLocation = discoveredGenerators.Apps[meta.Name]
	case generator.Infrastructure:
		generatorLocation = discoveredGenerators.Infrastructure[meta.Name]
	case generator.Istio:
		generatorLocation = discoveredGenerators.Istio[meta.Name]
	case generator.Monitoring:
		generatorLocation = discoveredGenerators.Monitoring[meta.Name]
	}

	logToOutput(outputView, "Running generator")
	_, err = utils.RunGeneratorMain(generatorLocation, []string{
		fmt.Sprintf("--root %v", rootDir),
	})
	if err != nil {
		logToOutput(outputView, fmt.Sprintf("Running the generator failed! Reason: \n %v", err))
	} else {
		logToOutput(outputView, "Succesfully ran the generator!")
	}
}
