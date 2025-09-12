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

	"github.com/rivo/tview"
)

func appendGeneratorsToTreeNode(treeNode *tview.TreeNode, generators []generator.GeneratorMeta, generatorCommandHandler func()) {

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

		generatorNode.SetSelectedFunc(generatorCommandHandler)

		treeNode.AddChild(generatorNode)
	}
}

func logToOutput(outputView *tview.TextView, message string) {
	outputView.Write([]byte(message))
	outputView.Write([]byte("\n"))
}

func clearOutput(outputView *tview.TextView) {
	outputView.Clear()
}

func makeNodeExpandable(node *tview.TreeNode) {
	node.SetSelectedFunc(func() {
		node.SetExpanded(!node.IsExpanded())
	})
}

func appendGeneratorsToTree(
	tree *tview.TreeView,
	appsGenerators []generator.GeneratorMeta,
	infrastructureGenerators []generator.GeneratorMeta,
	istioGenerators []generator.GeneratorMeta,
	monitoringGenerators []generator.GeneratorMeta,
	generatorCommandHandler func(),
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

func initializeGeneratorTree(rootDir string, outputView *tview.TextView, generatorsTree *tview.TreeView, generatorCommandHandler func()) {
	discoveredGeneratorsBytes, err := os.ReadFile(path.Join(rootDir, "clidata/discoveredgenerators.json"))
	if err != nil {
		logToOutput(outputView, fmt.Sprintf("An error happened while reading discoveredgenerators.json: \n %v", err))
	}

	discoveredGenerators := DiscoveredGenerators{}
	err = json.Unmarshal(discoveredGeneratorsBytes, &discoveredGenerators)
	if err != nil {
		logToOutput(outputView, fmt.Sprintf("An error happened while unmarshalling discoveredgenerators.json: \n %v", err))
	} else {
		appendGeneratorsToTree(
			generatorsTree,
			utils.GetGeneratorMetasByPaths(slices.Collect(maps.Values(discoveredGenerators.Apps))),
			utils.GetGeneratorMetasByPaths(slices.Collect(maps.Values(discoveredGenerators.Infrastructure))),
			utils.GetGeneratorMetasByPaths(slices.Collect(maps.Values(discoveredGenerators.Istio))),
			utils.GetGeneratorMetasByPaths(slices.Collect(maps.Values(discoveredGenerators.Monitoring))),
			generatorCommandHandler,
		)
	}
}

type GeneratorCommand string

var (
	Run     GeneratorCommand = "Run"
	Version GeneratorCommand = "Version"
	None    GeneratorCommand = "None"
)

func Start() {
	app := tview.NewApplication()
	commandList := tview.NewList()
	rootTreeNode := tview.NewTreeNode("Generators")
	generatorsTree := tview.NewTreeView().SetRoot(rootTreeNode).SetCurrentNode(rootTreeNode)
	outputView := tview.NewTextView()
	actionArea := tview.NewFlex()
	flex := tview.NewFlex()

	rootDir, err := utils.FindRoot()
	if err != nil {
		fmt.Println("Error while finding root, reason:")
		fmt.Print(err)
		return
	}

	app.EnableMouse(true)

	outputView.
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	outputView.SetTitle("Command output")
	outputView.SetBorder(true)
	outputView.SetScrollable(true)

	currentGeneratorCommand := None
	generatorCommandHandler := func() {
		switch currentGeneratorCommand {
		case None:
			logToOutput(outputView, "Current command: None")
		case Run:
			logToOutput(outputView, "Current command: Run")
		case Version:
			logToOutput(outputView, "Current command: Version")
		}
		app.SetFocus(commandList)
	}

	generatorsTree.SetBorder(true)
	initializeGeneratorTree(rootDir, outputView, generatorsTree, generatorCommandHandler)

	commandList.
		AddItem("Run", "Run a generator", 'r', func() {
			app.SetFocus(generatorsTree)
			currentGeneratorCommand = Run
		}).
		AddItem("Version", "Upgrade Helm/Docker versions", 'v', func() {
			app.SetFocus(generatorsTree)
			currentGeneratorCommand = Version
		}).
		AddItem("Discover", "Discover new generators", 'd', func() {
			clearOutput(outputView)
			logToOutput(outputView, "Starting discovery process...")

			apps, infrastructure, istio, monitoring := discoverGeneratorsByCategoryViaPath(fmt.Sprintf("%v/internal/generators", rootDir), rootDir, outputView)

			logToOutput(outputView, "Appending generators to tree")
			appendGeneratorsToTree(generatorsTree, apps, infrastructure, istio, monitoring, generatorCommandHandler)
			logToOutput(outputView, "Discovery completed")
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	commandList.SetTitle("Commands")
	commandList.SetBorder(true)

	actionArea.SetDirection(tview.FlexRow)
	actionArea.AddItem(outputView, 0, 3, false)

	flex.
		AddItem(commandList, 40, 1, true).
		AddItem(actionArea, 0, 3, false).
		AddItem(generatorsTree, 0, 1, false)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
