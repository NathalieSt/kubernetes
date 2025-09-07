package main

import (
	"kubernetes/pkg/schema/generator"
)

var Postgres = generator.GeneratorMeta{
	Name:                "postgres",
	Namespace:           "postgres",
	GeneratorType:       generator.App,
	ClusterUrl:          "postgres-rw.postgres.svc.cluster.local",
	Port:                5432,
	DependsOnGenerators: []string{},
}
