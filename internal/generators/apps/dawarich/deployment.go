package main

import (
	"fmt"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func getDeployment(
	generatorMeta generator.GeneratorMeta,
	redisMeta generator.GeneratorMeta,
	postgresMeta generator.GeneratorMeta,
	publicpvcName string,
	watchedpvcName string,
	storagepvcName string,
) apps.Deployment {

	sharedEnv := []core.Env{
		{
			Name:  "RAILS_ENV",
			Value: "development"},
		{
			Name:  "REDIS_URL",
			Value: fmt.Sprintf("redis://%v:%v", redisMeta.ClusterUrl, redisMeta.Port),
		},
		{
			Name:  "DATABASE_HOST",
			Value: postgresMeta.ClusterUrl,
		},
		{
			Name:  "DATABASE_PORT",
			Value: fmt.Sprintf("%v", postgresMeta.Port),
		},
		{
			Name: "DATABASE_USERNAME",
			ValueFrom: core.ValueFrom{
				SecretKeyRef: core.SecretKeyRef{
					Key:  "username",
					Name: "postgres-creds-secret",
				},
			},
		},
		{
			Name: "DATABASE_PASSWORD",
			ValueFrom: core.ValueFrom{
				SecretKeyRef: core.SecretKeyRef{
					Key:  "password",
					Name: "postgres-creds-secret",
				},
			},
		},
		{
			Name:  "DATABASE_NAME",
			Value: "dawarich-development",
		},
		{
			Name:  "APPLICATION_PROTOCOL",
			Value: "http",
		},
		{
			Name:  "PROMETHEUS_EXPORTER_ENABLED",
			Value: "false",
		},
		{
			Name:  "SELF_HOSTED",
			Value: "true",
		},
		{
			Name:  "STORE_GEODATA",
			Value: "true",
		},
		{
			Name:  "APPLICATION_HOSTS",
			Value: "dawarich.cluster,.netbird.selfhosted",
		},
	}

	dawarichEnv := []core.Env{
		{
			Name:  "MIN_MINUTES_SPENT_IN_CITY",
			Value: "60",
		},

		{
			Name:  "TIME_ZONE",
			Value: "Europe/Vienna",
		},
	}
	dawarichEnv = append(dawarichEnv, sharedEnv...)

	dawarichSidekiqEnv := []core.Env{
		{
			Name:  "BACKGROUND_PROCESSING_CONCURRENCY",
			Value: "10",
		},
	}
	dawarichSidekiqEnv = append(dawarichSidekiqEnv, sharedEnv...)

	publicVolume := "public-dawarich-volume"
	watchedVolume := "watched-dawarich-volume"
	storageVolume := "storage-dawarich-volume"

	dawarichVolumeMounts := []core.VolumeMount{
		{
			Name:      publicVolume,
			MountPath: "/var/app/public",
		},
		{
			Name:      watchedVolume,
			MountPath: "/var/app/tmp/imports/watched",
		},
		{
			Name:      storageVolume,
			MountPath: "/var/app/storage",
		},
	}

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
							Args: []string{
								"bin/rails",
								"server",
								"'-p'",
								"'3000'",
								"'-b'",
								"'::'",
							},
							Command: []string{
								"web-entrypoint.sh",
							},
							Ports: []core.Port{
								{
									ContainerPort: generatorMeta.Port,
									Name:          generatorMeta.Name,
								},
							},
							Env:          dawarichEnv,
							VolumeMounts: dawarichVolumeMounts,
						},
						{
							Name:  fmt.Sprintf("%v-sidekiq", generatorMeta.Name),
							Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
							Command: []string{
								"sidekiq-entrypoint.sh",
							},
							Ports: []core.Port{
								{
									ContainerPort: generatorMeta.Port,
									Name:          generatorMeta.Name,
								},
							},
							Env:          dawarichSidekiqEnv,
							VolumeMounts: dawarichVolumeMounts,
						},
					},
					Volumes: []core.Volume{
						{
							Name: publicVolume,
							PersistentVolumeClaim: core.PVCVolumeSource{
								ClaimName: publicpvcName,
							},
						},
						{
							Name: watchedVolume,
							PersistentVolumeClaim: core.PVCVolumeSource{
								ClaimName: watchedpvcName,
							},
						},
						{
							Name: storageVolume,
							PersistentVolumeClaim: core.PVCVolumeSource{
								ClaimName: storagepvcName,
							},
						},
					},
				},
			},
		},
	)
}
