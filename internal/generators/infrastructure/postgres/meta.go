package main

import (
	"kubernetes/pkg/schema/generator"
)

var Postgres = generator.GeneratorMeta{
	Name:          "postgres",
	Namespace:     "postgres",
	GeneratorType: generator.App,
	ClusterUrl:    "postgres-rw.postgres.svc.cluster.local",
	Docker: generator.Docker{
		Registry: "ghcr.io/cloudnative-pg/postgis",
		Version:  "17-3.5-177",
	},
	Port:                5432,
	DependsOnGenerators: []string{},
}
