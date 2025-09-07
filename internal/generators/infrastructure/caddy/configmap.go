package main

import (
	"bytes"
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
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

func getWebsocketSupportIfRequired(required bool) string {
	if required {
		return fmt.Sprintf(`
	@websockets {
		header Connection *Upgrade*
		header Upgrade    websocket
	}
	reverse_proxy @websockets %v:80 {
		header_up Host {host}
	}
		`, generators.IstioGatewayIP)
	}
	return ""
}

func getCaddyFile(exposedServicesMeta []generator.GeneratorMeta) string {
	caddyfileBuffer := bytes.Buffer{}
	for _, meta := range exposedServicesMeta {
		caddyfileBuffer.WriteString(fmt.Sprintf(`
%v.netbird.selfhosted:443 {
	tls internal
	reverse_proxy %v:80 {
		header_up Host {host}
		%v
	}
	%v
}
		`,
			meta.Caddy.DNSName,
			generators.IstioGatewayIP,
			forwardHeadersIfRequried(meta.Caddy.HeaderForwardingIsRequired),
			getWebsocketSupportIfRequired(meta.Caddy.WebsocketSupportIsRequired),
		),
		)
	}

	return caddyfileBuffer.String()
}

func getCaddyConfigMap(name string, metas []generator.GeneratorMeta) core.ConfigMap {

	return core.NewConfigMap(meta.ObjectMeta{
		Name: name,
	}, map[string]string{
		"Caddyfile": getCaddyFile(metas),
	})
}
