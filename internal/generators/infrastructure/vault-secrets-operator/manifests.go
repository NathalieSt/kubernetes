package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/vaultsecretsoperator"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
)

func createVaultSecretsOperatorManifests(generatorMeta generator.GeneratorMeta, rootDir string) (map[string][]byte, error) {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace, true),
	}

	repo, chart, release := utils.GetGenericHelmDeploymentManifests(generatorMeta.Name, generatorMeta.Helm,
		map[string]any{
			"controller": map[string]any{
				"annotations": map[string]any{
					"traffic.sidecar.istio.io/excludeOutboundPorts": "8200",
				},
			},
		},
		nil,
	)

	vaultMeta, err := utils.GetServiceMeta(rootDir, "internal/generators/infrastructure/vault")
	if err != nil {
		fmt.Println("An error happened while getting vault meta ")
		return nil, err
	}

	serviceAccount, role, rolebinding := utils.GenerateRBAC(generatorMeta.Name)

	vaultConfigs := utils.ManifestConfig{
		Filename: "vault-configs.yaml",
		Manifests: []any{
			vaultsecretsoperator.NewAuthGlobal(meta.ObjectMeta{
				Name: "default",
			}, vaultsecretsoperator.AuthGlobalSpec{
				AllowedNamespaces: []string{"reflector", "gluetun-proxy"},
				DefaultAuthMethod: "kubernetes",
				Kubernetes: vaultsecretsoperator.Kubernetes{
					Audiences:              []string{"vault"},
					Mount:                  "kubernetes",
					Role:                   "global-vault-auth",
					ServiceAccount:         serviceAccount.Metadata.Name,
					TokenExpirationSeconds: 600,
				},
			}),
			serviceAccount,
			role,
			rolebinding,
			vaultsecretsoperator.NewConnection(meta.ObjectMeta{
				Name: "default",
			}, vaultsecretsoperator.ConnectionSpec{
				Address: fmt.Sprintf("http://%v:8200", vaultMeta.ClusterUrl),
			}),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(
			generatorMeta.Name,
			[]string{
				namespace.Filename,
				repo.Filename,
				chart.Filename,
				release.Filename,
				vaultConfigs.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, repo, chart, release, vaultConfigs}), nil
}
