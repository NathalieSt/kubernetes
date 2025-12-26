package main

import (
	"bytes"
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"log"
	"os"
	"strings"
)

func forwardHeadersIfRequried(required bool) string {
	if required {
		return `
		header_up X-Forwarded-Proto {scheme}
		header_up X-Forwarded-Host {host}
		header_up X-Forwarded-Port {server_port}
		header_up X-Real-IP {remote_host}
		header_up X-Forwarded-Uri {uri}
		`
	}
	return ""
}

func getWebsocketSupportIfRequired(required bool, clusterUrl string, port int64) string {
	if required {
		return fmt.Sprintf(`
	@websockets {
		header Connection *Upgrade*
		header Upgrade    websocket
	}
	reverse_proxy @websockets http://%v:%v {
		header_up Host {host}
	}
		`, clusterUrl, port)
	}
	return ""
}

func getCaddyFile(exposedServicesMeta []generator.GeneratorMeta) string {
	caddyfileBuffer := bytes.Buffer{}
	for _, meta := range exposedServicesMeta {
		caddyfileBuffer.WriteString(fmt.Sprintf(`
%v.%v {
	tls {
        dns hetzner {$API_TOKEN}
        propagation_delay 30s
    }    
	reverse_proxy http://%v:%v {
		header_up Host {host}
		%v
	}
	%v
}
		`,
			strings.ReplaceAll(meta.Caddy.DNSName, ".cluster", ""),
			generators.NetbirdDomainBase,
			meta.ClusterUrl,
			meta.Port,
			forwardHeadersIfRequried(meta.Caddy.HeaderForwardingIsRequired),
			getWebsocketSupportIfRequired(meta.Caddy.WebsocketSupportIsRequired, meta.ClusterUrl, meta.Port),
		),
		)
	}

	return caddyfileBuffer.String()
}

func main() {
	root, err := utils.FindRoot()
	if err != nil {
		println("An error while finding the root has occurred")
		return
	}
	exposedGenerators, err := utils.GetMetaForExposedGenerators(root)
	if err != nil {
		println("An error while finding exposed generators")
		return
	}

	caddyfilestring := getCaddyFile(
		exposedGenerators,
	)
	println("Use the following script:")
	print(fmt.Sprintf(`
		netbird deregister
		netbird up --management-url https://netbird.nathalie-stiefsohn.eu --setup-key <insert-token here> --extra-dns-labels %v
	`, strings.Join(exposedGenerators.GetDNSNames(), ",")))

	if err := os.WriteFile("./out/caddyfile", []byte(caddyfilestring), 0666); err != nil {
		log.Fatal(err)
	}

}
