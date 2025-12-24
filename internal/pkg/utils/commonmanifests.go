package utils

import (
	"kubernetes/pkg/schema/cluster/infrastructure/keda"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/kustomize"
)

func GenerateCronScaler(name string, targetName string, kind keda.Kind, kedaScaling *keda.ScaledObjectTriggerMeta) []any {
	return []any{
		keda.NewScaledObject(
			meta.ObjectMeta{
				Name: name,
			}, keda.ScaledObjectSpec{
				ScaleTargetRef: keda.ScaleTargetRef{
					Name: targetName,
					Kind: kind,
				},
				MinReplicaCount: 0,
				CooldownPeriod:  300,
				Triggers: []keda.ScaledObjectTrigger{
					{
						ScalerType: keda.Cron,
						Metadata:   *kedaScaling,
					},
				},
			},
		),
	}
}

func GenerateKustomization(name string, resources []string) []any {
	return []any{
		kustomize.NewKustomization(
			meta.ObjectMeta{
				Name: name,
			},
			resources,
		),
	}
}

func GenerateNamespace(name string) []any {

	return []any{
		core.NewNamespace(meta.ObjectMeta{
			Name: name,
		}),
	}

}
