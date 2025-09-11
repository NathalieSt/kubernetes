package main

import (
	"kubernetes/pkg/schema/generator"
)

var Postgres = generator.GeneratorMeta{
	Name:          "postgres",
	Namespace:     "postgres",
	GeneratorType: generator.Infrastructure,
	ClusterUrl:    "postgres-rw.postgres.svc.cluster.local",
	Docker: generator.Docker{
		Registry: "ghcr.io/cloudnative-pg/postgis",
		Version:  "17",
	},
	Port:                5432,
	DependsOnGenerators: []string{},
}
