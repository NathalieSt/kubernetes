package cli

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"time"

	"github.com/rivo/tview"
)

func logToOutput(outputView *tview.TextView, message string) {
	outputView.Write([]byte(message))
	outputView.Write([]byte("\n"))
}

func clearOutput(outputView *tview.TextView) {
	outputView.Clear()
}

func Start() {
	start := time.Now()
	app := tview.NewApplication()
	commandList := tview.NewList()
	rootTreeNode := tview.NewTreeNode("Generators")
	generatorsTree := tview.NewTreeView().SetRoot(rootTreeNode).SetCurrentNode(rootTreeNode)
	outputView := tview.NewTextView()
	actionArea := tview.NewFlex()
	pages := tview.NewPages()
	mainFlex := tview.NewFlex()
	scaffoldingFlex := tview.NewFlex()
	newGeneratorMeta := generator.GeneratorMeta{}

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
	generatorCommandHandler := func(meta generator.GeneratorMeta) {
		switch currentGeneratorCommand {
		case None:
			logToOutput(outputView, "Current command: None")
		case Run:
			logToOutput(outputView, fmt.Sprintf("Running generator: %v", meta.Name))
			runGeneratorFromJSON(rootDir, meta, outputView)
			logToOutput(outputView, "")
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
			logToOutput(outputView, "")
		}).
		AddItem("Scaffolding", "Create scaffolding for new generator", 's', func() {
			pages.SwitchToPage("Scaffolding")
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	commandList.SetTitle("Commands")
	commandList.SetBorder(true)

	actionArea.SetDirection(tview.FlexRow)
	actionArea.AddItem(outputView, 0, 3, false)

	mainFlex.
		AddItem(commandList, 40, 1, true).
		AddItem(actionArea, 0, 3, false).
		AddItem(generatorsTree, 0, 1, false)

	generateScaffoldingPageLayout(pages, scaffoldingFlex, newGeneratorMeta)

	pages.AddPage("Main", mainFlex, true, true)
	pages.AddPage("Scaffolding", scaffoldingFlex, true, false)

	elapsed := time.Since(start)
	logToOutput(outputView, fmt.Sprintf("Initialization took: %s", elapsed))

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
