package infrastructure

import (
	"kubernetes/pkg/schema/generator"
)

var Postgres = generator.GeneratorMeta{
	Name:       "postgres",
	Namespace:  "postgres",
	EntityType: generator.Infrastructure,
	ClusterUrl: "postgres-rw.postgres.svc.cluster.local",
	Port:       5432,
}
