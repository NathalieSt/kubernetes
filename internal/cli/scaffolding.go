package cli

import (
	"fmt"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/rivo/tview"
)

/*
type GeneratorMeta struct {
    Name                string
    Namespace           string
    GeneratorType       GeneratorType
    ClusterUrl          string
    Port                int64
    Docker              *Docker
    Helm                *Helm
    Caddy               *Caddy
    VirtualService      *VirtualServiceConfig
    KedaScaling         *keda.ScaledObjectTriggerMeta
    DependsOnGenerators []string
}*/

func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func GetGeneratorTypeString(generatorType generator.GeneratorType) string {
	switch generatorType {
	case generator.App:
		return "generator.App"
	case generator.Infrastructure:
		return "generator.Infrastructure"
	case generator.Istio:
		return "generator.Istio"
	case generator.Monitoring:
		return "generator.Monitoring"
	default:
		return "generator.App"
	}
}

func GetGeneratorTypeSubfolderString(generatorType generator.GeneratorType) string {
	switch generatorType {
	case generator.App:
		return "apps"
	case generator.Infrastructure:
		return "infrastructure"
	case generator.Istio:
		return "istio"
	case generator.Monitoring:
		return "monitoring"
	default:
		return "apps"
	}
}

func RemoveWhitespaces(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func getMainTemplate() (*template.Template, error) {
	template, err := template.New("main.go").Funcs(template.FuncMap{
		"RemoveWithspaces":                RemoveWhitespaces,
		"CapitalizeFirst":                 CapitalizeFirst,
		"ToLower":                         strings.ToLower,
		"GetGeneratorTypeString":          GetGeneratorTypeString,
		"GetGeneratorTypeSubfolderString": GetGeneratorTypeSubfolderString,
	}).Parse(`
package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"path/filepath"
)

func main() {
	flags := utils.GetGeneratorFlags()
	if flags == nil {
		fmt.Println("An error happened while getting flags for generator")
		return
	}

	name := "{{.Name | ToLower}}"
	generatorType := {{.GeneratorType | GetGeneratorTypeString}}
	var meta = generator.GeneratorMeta{
		Name:          name,
		Namespace:     "{{.Namespace}}",
		GeneratorType: generatorType,
		ClusterUrl:    "{{.Name}}.{{.Namespace}}.svc.cluster.local",
		Port:          {{.Port}},
		{{if .Docker}}
		Docker: &generator.Docker{
			Registry: "{{.Docker.Registry}}",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		{{end}}
		{{if .Helm}}
		Helm: &generator.Helm{
			Url:     "{{.Helm.Url}}",
			Chart:   "{{.Helm.Chart}}",
			Version: utils.GetGeneratorVersionByType(flags.RootDir, name, generatorType),
		},
		{{end}}
		{{if .Caddy}}
		Caddy: &generator.Caddy{
			DNSName: "{{.Caddy.DNSName}}",
		},
		{{end}}
		{{if .KedaScaling}}
		KedaScaling: &keda.ScaledObjectTriggerMeta{
			Timezone:        "{{.KedaScaling.Timezone}}",
			Start:           "{{.KedaScaling.Start}}",
			End:             "{{.KedaScaling.End}}",
			DesiredReplicas: "{{.KedaScaling.DesiredReplicas}}",
		},
		{{end}}
		DependsOnGenerators: []string{
			{{range .DependsOnGenerators}}
			"{{.}}"
			{{end}}
		},
	}

	utils.RunGenerator(utils.GeneratorRunnerConfig{
		Meta:             meta,
		ShouldReturnMeta: flags.ShouldReturnMeta,
		OutputDir:        filepath.Join(flags.RootDir, "/cluster/{{.GeneratorType | GetGeneratorTypeSubfolderString}}/{{.Name | ToLower }}/"),
		CreateManifests: func(gm generator.GeneratorMeta) map[string][]byte {
			manifests, err := create{{.Name | CapitalizeFirst | RemoveWithspaces}}Manifests(gm, flags.RootDir)
			if err != nil {
				fmt.Println("An error happened while generating {{.Name}} Manifests")
				fmt.Printf("Reason:\n %v", err)
				return nil
			}
			return manifests
		},
	})
}
`)

	if err != nil {
		return nil, err
	}

	return template, nil
}

func getManifestsTemplate() (*template.Template, error) {
	template, err := template.New("manifests").Parse(`
	
	`)

	if err != nil {
		return nil, err
	}

	return template, nil
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func writeTemplatesToFiles(generatorMeta generator.GeneratorMeta, outDir string, templates ...*template.Template) {

	for _, template := range templates {
		f, err := create(path.Join(outDir, fmt.Sprintf("internal/generators/%v/%v", GetGeneratorTypeSubfolderString(generatorMeta.GeneratorType), strings.ReplaceAll(strings.ToLower(generatorMeta.Name), " ", "")), template.Name()))
		if err != nil {
			fmt.Println("Couldnt create file")
		}
		err = template.Execute(f, generatorMeta)
		if err != nil {
			fmt.Println("Error executing template")
		}
		f.Close()
	}
}

func generateGoFiles(generatorMeta generator.GeneratorMeta, outdir string) {
	main, err := getMainTemplate()
	if err != nil {
		fmt.Println("An error happened while getting Main template")
		fmt.Println(err.Error())
		return
	}

	writeTemplatesToFiles(generatorMeta, outdir, main)
}

func generateScaffoldingPageLayout(rootDir string, pages *tview.Pages, flex *tview.Flex) {
	generatorMeta := generator.GeneratorMeta{}
	name := ""
	namespace := ""
	typeString := ""
	clusterUrl := ""
	port := new(int64)
	var helm *generator.Helm
	var docker *generator.Docker
	var kedaScaling *keda.ScaledObjectTriggerMeta
	var caddy *generator.Caddy

	form := tview.NewForm()

	form.SetTitle("New Generator")
	form.SetBorder(true)

	form.AddInputField("Name:", "", 20, nil, func(text string) {
		name = text
	})
	form.AddInputField("Namespace:", "", 20, nil, func(text string) {
		namespace = text
	})
	form.AddDropDown("Type:", []string{"App", "Infrastructure", "Istio", "Monitoring"}, 0, func(option string, optionIndex int) {
		typeString = option
	})
	form.AddInputField("Cluster url:", "<service-name>.<namespace>.svc.cluster.local", 50, nil, func(text string) {
		clusterUrl = text
	})
	form.AddInputField("Port:", "<1-65535>", 20, func(textToCheck string, lastChar rune) bool {
		if parsed, err := strconv.ParseInt(textToCheck, 10, 64); err == nil {
			if parsed > 0 && parsed <= 65535 {
				return true
			}
			return false
		}
		return false
	}, func(text string) {
		if parsed, err := strconv.ParseInt(text, 10, 64); err == nil {
			port = &parsed
		}
	})

	registryInput := tview.NewInputField().SetLabel("Container Registry:").SetText("").SetChangedFunc(func(text string) {
		docker.Registry = text
	})
	containerVersionInput := tview.NewInputField().SetLabel("Container Version:").SetText("").SetChangedFunc(func(text string) {
		docker.Version = text
	})

	chartRepoInput := tview.NewInputField().SetLabel("Chart repo:").SetText("").SetChangedFunc(func(text string) {
		helm.Url = text
	})
	chartNameInput := tview.NewInputField().SetLabel("Chart name:").SetText("").SetChangedFunc(func(text string) {
		helm.Chart = text
	})
	chartVersionInput := tview.NewInputField().SetLabel("Chart version").SetText("").SetChangedFunc(func(text string) {
		helm.Version = text
	})

	form.AddDropDown("Deployment type:", []string{"Manual Deployment", "Helm Chart", "None"}, 2, func(option string, optionIndex int) {
		switch option {
		case "Manual Deployment":
			helm = nil
			docker = &generator.Docker{}
			chartRepoInput.SetDisabled(true)
			chartNameInput.SetDisabled(true)
			chartVersionInput.SetDisabled(true)

			registryInput.SetDisabled(false)
			containerVersionInput.SetDisabled(false)

		case "Helm Chart":
			docker = nil
			helm = &generator.Helm{}
			registryInput.SetDisabled(true)
			containerVersionInput.SetDisabled(true)

			chartRepoInput.SetDisabled(false)
			chartNameInput.SetDisabled(false)
			chartVersionInput.SetDisabled(false)

		case "None":
			chartRepoInput.SetDisabled(true)
			chartNameInput.SetDisabled(true)
			chartVersionInput.SetDisabled(true)
			registryInput.SetDisabled(true)
			containerVersionInput.SetDisabled(true)
			helm = nil
			docker = nil
		}
	})

	form.AddFormItem(registryInput).
		AddFormItem(containerVersionInput)

	form.AddFormItem(chartRepoInput).
		AddFormItem(chartNameInput).
		AddFormItem(chartVersionInput)

	kedaTimezoneInput := tview.NewInputField().SetLabel("Timezone:").SetText("").SetChangedFunc(func(text string) {
		kedaScaling.Timezone = text
	}).SetDisabled(true)
	kedaStartInput := tview.NewInputField().SetLabel("Start:").SetText("").SetChangedFunc(func(text string) {
		kedaScaling.Start = text
	}).SetDisabled(true)
	kedaEndInput := tview.NewInputField().SetLabel("End:").SetText("").SetChangedFunc(func(text string) {
		kedaScaling.End = text
	}).SetDisabled(true)
	kedaReplicasInput := tview.NewInputField().SetLabel("Desired Replicas:").SetText("").SetChangedFunc(func(text string) {
		kedaScaling.DesiredReplicas = text
	}).SetDisabled(true)

	form.AddCheckbox("Keda scaling required:", false, func(checked bool) {
		if checked {
			kedaScaling = &keda.ScaledObjectTriggerMeta{}
			kedaTimezoneInput.SetDisabled(false)
			kedaStartInput.SetDisabled(false)
			kedaEndInput.SetDisabled(false)
			kedaReplicasInput.SetDisabled(false)
		} else {
			kedaTimezoneInput.SetDisabled(true)
			kedaStartInput.SetDisabled(true)
			kedaEndInput.SetDisabled(true)
			kedaReplicasInput.SetDisabled(true)
			kedaScaling = nil
		}
	})

	form.AddFormItem(kedaTimezoneInput).
		AddFormItem(kedaStartInput).
		AddFormItem(kedaEndInput).
		AddFormItem(kedaReplicasInput)

	dnsNameInput := tview.NewInputField().SetLabel("DNS Name:").SetText("<name>.cluster").SetChangedFunc(func(text string) {
		caddy.DNSName = text
	}).SetDisabled(true)
	requiresWebsocketsInput := tview.NewCheckbox().SetLabel("Requires Websockets:").SetChecked(false).SetChangedFunc(func(checked bool) {
		if checked {
			caddy.WebsocketSupportIsRequired = true
		} else {
			caddy.WebsocketSupportIsRequired = false
		}
	}).SetDisabled(true)

	form.AddCheckbox("Is externally available:", false, func(checked bool) {
		if checked {
			caddy = &generator.Caddy{}
			dnsNameInput.SetDisabled(false)
			requiresWebsocketsInput.SetDisabled(false)
		} else {
			dnsNameInput.SetDisabled(true)
			requiresWebsocketsInput.SetDisabled(true)
			caddy = nil
		}
	})

	form.AddFormItem(dnsNameInput).
		AddFormItem(requiresWebsocketsInput)

	form.AddButton("Submit", func() {
		generatorMeta.Name = name
		generatorMeta.Namespace = namespace
		generatorMeta.ClusterUrl = clusterUrl
		switch typeString {
		case "App":
			generatorMeta.GeneratorType = generator.App
		case "Infrastructure":
			generatorMeta.GeneratorType = generator.Infrastructure
		case "Istio":
			generatorMeta.GeneratorType = generator.Istio
		case "Monitoring":
			generatorMeta.GeneratorType = generator.Monitoring
		default:
			generatorMeta.GeneratorType = generator.App
		}
		if port != nil {
			generatorMeta.Port = *port
		}
		if helm != nil {
			generatorMeta.Helm = helm
		}
		if docker != nil {
			generatorMeta.Docker = docker
		}
		if caddy != nil {
			generatorMeta.Caddy = caddy
		}
		if kedaScaling != nil {
			generatorMeta.KedaScaling = kedaScaling
		}

		generateGoFiles(generatorMeta, rootDir)

		pages.SwitchToPage("Main")
	})
	form.AddButton("Cancel", func() {
		pages.SwitchToPage("Main")
	})

	flex.AddItem(form, 0, 1, true)
}
