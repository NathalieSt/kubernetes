package main

import (
	"bytes"
	"fmt"
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

func getWebsocketSupportIfRequired(required bool, clusterUrl string) string {
	if required {
		return fmt.Sprintf(`
	@websockets {
		header Connection *Upgrade*
		header Upgrade    websocket
	}
	reverse_proxy @websockets %v:80 {
		header_up Host {host}
	}
		`, clusterUrl)
	}
	return ""
}

func getCaddyFile(exposedServicesMeta []generator.GeneratorMeta) string {
	caddyfileBuffer := bytes.Buffer{}
	for _, meta := range exposedServicesMeta {
		caddyfileBuffer.WriteString(fmt.Sprintf(`
%v.cloud.nathalie-stiefsohn.eu:443 {
	tls internal
	reverse_proxy %v:%v {
		header_up Host {host}
		%v
	}
	%v
}
		`,
			meta.Caddy.DNSName,
			meta.ClusterUrl,
			meta.Port,
			forwardHeadersIfRequried(meta.Caddy.HeaderForwardingIsRequired),
			getWebsocketSupportIfRequired(meta.Caddy.WebsocketSupportIsRequired, meta.ClusterUrl),
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
