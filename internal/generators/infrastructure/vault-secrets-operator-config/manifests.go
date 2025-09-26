package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/cluster/infrastructure/vaultsecretsoperator"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/meta"
	"path"
)

func createVaultSecretsOperatorManifests(rootDir string, generatorMeta generator.GeneratorMeta) (map[string][]byte, error) {

	serviceAccount, role, rolebinding := utils.GenerateRBAC(generatorMeta.Name)

	vaultMeta, err := utils.GetGeneratorMeta(rootDir, path.Join(rootDir, "internal/generators/infrastructure/vault"))
	if err != nil {
		fmt.Println("An error happened while getting vault meta ")
		return nil, err
	}

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
				vaultConfigs.Filename,
			},
		),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{kustomization, vaultConfigs}), nil
}
