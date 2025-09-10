package main

import (
	"kubernetes/pkg/schema/generator"
)

var VictoriaMetrics = generator.GeneratorMeta{
	Name:          "victoria-metrics",
	Namespace:     "victoria-metrics",
	GeneratorType: generator.Monitoring,
	ClusterUrl:    "vmsingle-victoria-metrics-vmks.victoria-metrics.svc.cluster.local",
	Port:          20001,
	Helm: generator.Helm{
		Url:     "oci://ghcr.io/victoriametrics/helm-charts/victoria-metrics-k8s-stack",
		Version: "0.59.3",
	},
	DependsOnGenerators: []string{},
}
