package cli

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"

	"github.com/rivo/tview"
)

func getGeneratorsFromJSON() []generator.GeneratorMeta {
	return []generator.GeneratorMeta{}
}

func appendGeneratorsToTreeNode(treeNode *tview.TreeNode, generators []generator.GeneratorMeta) {
	//treeNode.ClearChildren()

	if len(generators) == 0 {
		treeNode.AddChild(tview.NewTreeNode("Generators").SetText("No generators available\nPlease try to run the \"discover\" command to find generators"))
		return
	}

	for _, generator := range generators {
		version := ""

		if generator.Helm.Version != "" {
			version = fmt.Sprintf("(%v)", generator.Helm.Version)
		} else if generator.Docker.Version != "" {
			version = fmt.Sprintf("(%v)", generator.Docker.Version)
		}

		treeNode.AddChild(tview.NewTreeNode(fmt.Sprintf("%v %v", generator.Name, version)))
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
) {
	rootNode := tree.GetRoot()

	applicationNode := tview.NewTreeNode("Applications").SetExpanded(false)
	makeNodeExpandable(applicationNode)

	infrastructureNode := tview.NewTreeNode("Infrastructure").SetExpanded(false)
	makeNodeExpandable(infrastructureNode)

	istioNode := tview.NewTreeNode("Istio").SetExpanded(false)
	makeNodeExpandable(istioNode)

	monitoringNode := tview.NewTreeNode("Monitoring").SetExpanded(false)
	makeNodeExpandable(monitoringNode)

	rootNode.AddChild(applicationNode).AddChild(infrastructureNode).AddChild(istioNode).AddChild(monitoringNode)

	appendGeneratorsToTreeNode(applicationNode, appsGenerators)
	appendGeneratorsToTreeNode(infrastructureNode, infrastructureGenerators)
	appendGeneratorsToTreeNode(istioNode, istioGenerators)
	appendGeneratorsToTreeNode(monitoringNode, monitoringGenerators)
}

func Start() {

	app := tview.NewApplication()

	app.EnableMouse(true)

	rootDir, err := utils.FindRoot()
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
	root := tview.NewTreeNode("Generators")
	generatorsTree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	commandList := tview.NewList().
		AddItem("Run", "Run a generator", 'r', nil).
		AddItem("Version", "Upgrade Helm/Docker versions", 'v', nil).
		AddItem("Discover", "Discover new generators", 'd', func() {
			clearOutput(outputView)
			logToOutput(outputView, "Starting discovery process...")

			apps, infrastructure, istio, monitoring := discoverGeneratorsByCategoryViaPath(fmt.Sprintf("%v/internal/generators", rootDir), rootDir, outputView)

			logToOutput(outputView, "Appending generators to tree")
			appendGeneratorsToTree(generatorsTree, apps, infrastructure, istio, monitoring)
			logToOutput(outputView, "Discovery completed")
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	commandList.SetTitle("Commands")
	commandList.SetBorder(true)

	actionArea := tview.NewFlex().SetDirection(tview.FlexRow)
	actionArea.AddItem(outputView, 0, 3, false)

	flex := tview.NewFlex().
		AddItem(commandList, 40, 1, true).
		AddItem(actionArea, 0, 3, false).
		AddItem(generatorsTree, 0, 1, false)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
