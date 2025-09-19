package cli

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/generator"
	"strconv"

	"github.com/rivo/tview"
)

func generateScaffoldingPageLayout(pages *tview.Pages, flex *tview.Flex, generatorMeta generator.GeneratorMeta) {
	name := ""
	namespace := ""
	typeString := ""
	clusterUrl := ""
	port := new(int64)
	helm := &generator.Helm{}
	docker := &generator.Docker{}
	keda := &keda.ScaledObjectTriggerMeta{}
	caddy := &generator.Caddy{}

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
			chartRepoInput.SetDisabled(true)
			chartNameInput.SetDisabled(true)
			chartVersionInput.SetDisabled(true)

			registryInput.SetDisabled(false)
			containerVersionInput.SetDisabled(false)

		case "Helm Chart":
			docker = nil
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
		keda.Timezone = text
	}).SetDisabled(true)
	kedaStartInput := tview.NewInputField().SetLabel("Start:").SetText("").SetChangedFunc(func(text string) {
		keda.Start = text
	}).SetDisabled(true)
	kedaEndInput := tview.NewInputField().SetLabel("End:").SetText("").SetChangedFunc(func(text string) {
		keda.End = text
	}).SetDisabled(true)
	kedaReplicasInput := tview.NewInputField().SetLabel("Desired Replicas:").SetText("").SetChangedFunc(func(text string) {
		keda.DesiredReplicas = text
	}).SetDisabled(true)

	form.AddCheckbox("Keda scaling required:", false, func(checked bool) {
		if checked {
			kedaTimezoneInput.SetDisabled(false)
			kedaStartInput.SetDisabled(false)
			kedaEndInput.SetDisabled(false)
			kedaReplicasInput.SetDisabled(false)
		} else {
			kedaTimezoneInput.SetDisabled(true)
			kedaStartInput.SetDisabled(true)
			kedaEndInput.SetDisabled(true)
			kedaReplicasInput.SetDisabled(true)
			keda = nil
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

		pages.SwitchToPage("Main")
	})
	form.AddButton("Cancel", func() {
		pages.SwitchToPage("Main")
	})

	flex.AddItem(form, 0, 1, true)
}
