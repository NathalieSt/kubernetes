package main

import (
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/flux/helm"
	"kubernetes/pkg/schema/generator"
)

func createPipedManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {

	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"frontend": map[string]any{
				"securityContext": map[string]any{
					"capabilities": map[string]any{
						"add": []string{
							"NET_BIND_SERVICE",
						},
					},
				},
				"env": map[string]any{
					"BACKEND_HOSTNAME": "backend",
				},
			},
			"backend": map[string]any{
				"config": map[string]any{
					"API_URL":      "http://backend",
					"PROXY_PART":   "http://proxy",
					"FRONTEND_URL": "http://frontend",
					"database": map[string]any{
						"connection_url": "jdbc:postgresql://postgres-rw.postgres.svc.cluster.local:5432/piped",
						"driver_class":   "org.postgresql.Driver",
						"dialect":        "org.hibernate.dialect.PostgreSQLDialect",
					},
				},
			},
			"ingress": map[string]any{
				"main": map[string]any{
					"enabled": false,
				},
				"backend": map[string]any{
					"enabled": false,
				},
				"ytproxy": map[string]any{
					"enabled": false,
				},
			},
		},
		[]helm.ReleaseValuesFrom{
			{
				Kind:       helm.Secret,
				Name:       generators.PostgresCredsSecret,
				ValuesKey:  "username",
				TargetPath: "backend.config.database.username",
				Optional:   false,
			},
			{
				Kind:       helm.Secret,
				Name:       generators.PostgresCredsSecret,
				ValuesKey:  "password",
				TargetPath: "backend.config.database.password",
				Optional:   false,
			},
		},
	)

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				namespace.Filename,
				repo.Filename,
				chart.Filename,
				release.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release})
}
