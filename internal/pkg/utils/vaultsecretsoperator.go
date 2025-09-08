package utils

import (
	"fmt"
	"kubernetes/pkg/schema/cluster/infrastructure/vaultsecretsoperator"
	"kubernetes/pkg/schema/k8s/authorization"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

type StaticSecretConfig struct {
	Name              string
	SecretName        string
	Path              string
	SecretLabels      map[string]string
	SecretAnnotations map[string]string
}

func generateStaticSecrets(secretsConfig []StaticSecretConfig, authName string) []any {
	secrets := []any{}
	for _, secret := range secretsConfig {
		secrets = append(secrets, vaultsecretsoperator.NewStaticSecret(
			meta.ObjectMeta{
				Name: secret.Name,
			},
			vaultsecretsoperator.StaticSecretSpec{
				AuthRef:      authName,
				Mount:        "kvv2",
				Type:         "kv-v2",
				Path:         secret.Path,
				RefreshAfter: "10s",
				Destination: vaultsecretsoperator.Destination{
					Create:      true,
					Name:        secret.SecretName,
					Labels:      secret.SecretLabels,
					Annotations: secret.SecretAnnotations,
				},
			},
		))
	}
	return secrets
}

func generateAuth(serviceName string, serviceAccountName string, globalAuthNamespace string) vaultsecretsoperator.Auth {
	return vaultsecretsoperator.NewAuth(
		meta.ObjectMeta{
			Name: fmt.Sprintf("%v-vault-auth", serviceName),
		},
		vaultsecretsoperator.AuthSpec{
			Kubernetes: vaultsecretsoperator.Kubernetes{
				Role:           serviceName,
				ServiceAccount: serviceAccountName,
			},
			VaultAuthGlobalRef: vaultsecretsoperator.AuthGlobalRef{
				AllowDefault: true,
				Namespace:    globalAuthNamespace,
			},
		})
}

func generateRBAC(serviceName string) (core.ServiceAccount, authorization.Role, authorization.RoleBinding) {
	serviceAccount := core.NewServiceAccount(meta.ObjectMeta{
		Name: fmt.Sprintf("%v-vault-serviceaccount", serviceName),
	})

	role := authorization.NewRole(meta.ObjectMeta{
		Name: "vault-token-reviewer",
	}, []authorization.Rule{
		{
			APIGroups: []string{"authentication.k8s.io"},
			Resources: []string{"tokenreviews"},
			Verbs:     []string{"create"},
		},
		{
			APIGroups: []string{"authorization.k8s.io"},
			Resources: []string{"subjectaccessreviews"},
			Verbs:     []string{"create"},
		},
	})

	rolebinding := authorization.NewRoleBinding(meta.ObjectMeta{
		Name: fmt.Sprintf("%v-vault-rolebinding", serviceName),
	}, authorization.RoleRef{
		Kind:     role.Kind,
		Name:     role.Metadata.Name,
		APIGroup: "rbac.authorization.k8s.io",
	}, []authorization.Subject{
		{
			Kind: serviceAccount.Kind,
			Name: serviceAccount.Metadata.Name,
		},
	},
	)

	return serviceAccount, role, rolebinding
}

func GenerateVaultAccessManifests(serviceName string, globalAuthNamespace string, secretsConfig []StaticSecretConfig) []any {
	serviceAccount, role, rolebinding := generateRBAC(serviceName)

	auth := generateAuth(serviceName, serviceAccount.Metadata.Name, globalAuthNamespace)

	staticSecrets := generateStaticSecrets(secretsConfig, auth.Metadata.Name)

	manifests := append([]any{serviceAccount, role, rolebinding, auth}, staticSecrets...)

	return manifests
}
