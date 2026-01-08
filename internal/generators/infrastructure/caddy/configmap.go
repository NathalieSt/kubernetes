package main

import (
	"bytes"
	"fmt"
	"kubernetes/internal/generators/shared"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"strings"
)

func getCaddyFile(exposedServicesMeta []generator.GeneratorMeta) string {
	caddyfileBuffer := bytes.Buffer{}
	caddyfileBuffer.WriteString(`
*.cloud.nathalie-stiefsohn.eu {	
	tls {
        dns hetzner {$API_TOKEN}
        propagation_delay 30s
    }
	`)
	for index, meta := range exposedServicesMeta {
		caddyfileBuffer.WriteString(fmt.Sprintf(`
	@svc%v host %v.%v
	
	handle @svc%v {
		reverse_proxy %v:%v
	} 
		`,
			index,
			strings.ReplaceAll(meta.Caddy.DNSName, ".cluster", ""),
			shared.NetbirdDomainBase,
			index,
			meta.ClusterUrl,
			meta.Port,
		),
		)
	}
	caddyfileBuffer.WriteString(`
    handle {
        respond "Not Found" 404
    }
}
	`)

	return caddyfileBuffer.String()
}

func getCaddyConfigMap(name string, metas []generator.GeneratorMeta) core.ConfigMap {

	return core.NewConfigMap(meta.ObjectMeta{
		Name: name,
	}, map[string]string{
		"Caddyfile": getCaddyFile(metas),
	})
}
