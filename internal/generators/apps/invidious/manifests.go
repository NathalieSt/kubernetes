package main

import (
	"fmt"
	"kubernetes/internal/generators"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

func createInvidiousManifests(generatorMeta generator.GeneratorMeta, rootDir string, relativeDir string) (map[string][]byte, error) {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	cachePVCName := "cache-pvc"
	cachePVC := utils.ManifestConfig{
		Filename: "cache-pvc.yaml",
		Manifests: []any{
			core.NewPersistentVolumeClaim(meta.ObjectMeta{
				Name: cachePVCName,
			}, core.PersistentVolumeClaimSpec{
				AccessModes: []string{"ReadWriteMany"},
				Resources: core.VolumeResourceRequirements{Requests: map[string]string{
					"storage": "10Gi",
				}},
				StorageClassName: generators.NFSRemoteClass,
			},
			),
		},
	}

	/*
			TODO: do it like this
			        - command: ['/bin/sh']
		          args:
		            - -c
		            - |
		              export INVIDIOUS_CONFIG=$(echo "$INVIDIOUS_CONFIG" | sed \
		                -e "s/__dbname/$INVIDIOUS_DB_DBNAME/" \
		                -e "s/__user/$INVIDIOUS_DB_USER/" \
		                -e "s/__password/$INVIDIOUS_DB_PASSWORD/" \
		                -e "s/__host/$INVIDIOUS_DB_HOST/" \
		                -e "s/__hmac_key/$INVIDIOUS_HMAC_KEY/")
		              exec /invidious/invidious
		          env:
		            - name: INVIDIOUS_CONFIG
		              value: |
		                db:
		                  dbname: __dbname
		                  user: __user
		                  password: __password
		                  host: __host
		                  port: 5432
		                check_tables: true
		                hmac_key: __hmac_key
		                channel_threads: 4
		                feed_threads: 4
		                pool_size: 2000
		                captcha_enabled: false
		                disable_proxy: false
		                default_user_preferences:
		                  local: true
		                  quality: dash
		                  quality_dash: auto

	*/

	cachePVCVolume := "cache-pvc-volume"
	deployment := utils.ManifestConfig{
		Filename: "deployment.yaml",
		Manifests: []any{
			apps.NewDeployment(
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
									Env: []core.Env{
										core.Env{
											Name: "INVIDIOUS_CONFIG",
											Value: `
db:
  dbname: 1234
  user: test
  password: test
  host: test.svc.cluster.local
  port: 5432
check_tables: true
hmac_key: 1234567890123456
channel_threads: 4
feed_threads: 4
pool_size: 2000
captcha_enabled: false
disable_proxy: false
default_user_preferences:
local: true
quality: dash
quality_dash: auto
											`,
										},
									},
								},
								{
									Name:  "invidious-companion",
									Image: "quay.io/invidious/invidious-companion:latest",
									Env:   []core.Env{},
									VolumeMounts: []core.VolumeMount{
										{
											MountPath: "/var/tmp/youtubei.js",
											Name:      cachePVCVolume,
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: cachePVCVolume,
									PersistentVolumeClaim: core.PVCVolumeSource{
										ClaimName: cachePVCName,
									},
								},
							},
						},
					},
				},
			),
		},
	}

	service := utils.ManifestConfig{
		Filename: "service.yaml",
		Manifests: []any{
			core.NewService(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				}, core.ServiceSpec{
					Selector: map[string]string{
						"app.kubernetes.io/name": generatorMeta.Name,
					},
					Ports: []core.ServicePort{
						{
							Name:       fmt.Sprintf("http-%v", generatorMeta.Name),
							Port:       generatorMeta.Port,
							TargetPort: generatorMeta.Port,
							Protocol:   core.TCP,
						},
					},
				},
			),
		},
	}

	kustomization := utils.ManifestConfig{
		Filename: "kustomization.yaml",
		Manifests: utils.GenerateKustomization(generatorMeta.Name, []string{
			namespace.Filename,
			service.Filename,
			cachePVC.Filename,
			deployment.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, cachePVC, deployment, service}), nil
}
