package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func getDeployment(generatorMeta generator.GeneratorMeta, configmapName string, certpvcName string, servicesDNSName string) apps.Deployment {
	certpvcVolume := "caddy-cert-pvc-volume"
	configmapVolume := "caddy-config-volume"
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
									MountPath: "/data/caddy/pki/authorities/local/",
									Name:      certpvcVolume,
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
											Name: generators.NetbirdNecretName,
											Key:  "setup-key",
										},
									},
								},
								{
									Name:  "NB_HOSTNAME",
									Value: "cluster",
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
									"cpu":    "50mi",
									"memory": "64mi",
								},
								Limits: map[string]string{
									"cpu":    "100mi",
									"memory": "128mi",
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
							Name: certpvcVolume,
							PersistentVolumeClaim: core.PVCVolumeSource{
								ClaimName: certpvcName,
							},
						},
					},
				},
			},
		},
	)
}
