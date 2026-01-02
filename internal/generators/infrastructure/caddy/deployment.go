package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func getDeployment(generatorMeta generator.GeneratorMeta, configmapName string, servicesDNSName string, pvcName string) apps.Deployment {
	configmapVolume := "caddy-config-volume"
	pvcVolume := "caddy-pvc-volume"
	return apps.NewDeployment(
		meta.ObjectMeta{
			Name: generatorMeta.Name,
			Labels: map[string]string{
				"app.kubernetes.io/name":    generatorMeta.Name,
				"app.kubernetes.io/version": generatorMeta.Docker.Version,
			},
		},
		apps.DeploymentSpec{
			Replicas: 1,
			Selector: meta.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": generatorMeta.Name,
				},
			},
			Template: core.PodTemplateSpec{
				Metadata: meta.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  generatorMeta.Name,
							Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
							Ports: []core.Port{
								{
									ContainerPort: generatorMeta.Port,
									Name:          generatorMeta.Name,
								},
							},
							VolumeMounts: []core.VolumeMount{
								{
									MountPath: "/etc/caddy/",
									Name:      configmapVolume,
								},
								{
									MountPath: "/data",
									Name:      pvcVolume,
								},
							},
							Env: []core.Env{
								{
									Name: "API_TOKEN",
									ValueFrom: core.ValueFrom{
										SecretKeyRef: core.SecretKeyRef{
											Name: generators.HetznerAPITokenSecretName,
											Key:  "api-token",
										},
									},
								},
							},
						},
						{
							Name:  "netbird-agent",
							Image: "netbirdio/netbird:latest",
							Env: []core.Env{
								{
									Name: "NB_SETUP_KEY",
									ValueFrom: core.ValueFrom{
										SecretKeyRef: core.SecretKeyRef{
											Name: generators.NetbirdSecretName,
											Key:  "setup-key",
										},
									},
								},
								{
									Name:  "NB_HOSTNAME",
									Value: "Caddy",
								},
								{
									Name:  "NB_MANAGEMENT_URL",
									Value: "https://netbird.nathalie-stiefsohn.eu",
								},
								{
									Name:  "NB_EXTRA_DNS_LABELS",
									Value: servicesDNSName,
								},
							},
							Resources: core.Resources{
								Requests: map[string]string{
									"cpu":    "50m",
									"memory": "64Mi",
								},
								Limits: map[string]string{
									"cpu":    "100m",
									"memory": "128Mi",
								},
							},
							SecurityContext: core.ContainerSecurityContext{
								Privileged: true,
							},
						},
					},
					Volumes: []core.Volume{
						{
							Name: configmapVolume,
							ConfigMap: core.ConfigMapVolumeSource{
								Name: configmapName,
							},
						},
						{
							Name: pvcVolume,
							PersistentVolumeClaim: core.PVCVolumeSource{
								ClaimName: pvcName,
							},
						},
					},
				},
			},
		},
	)
}
