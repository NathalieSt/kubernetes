package main

import (
	"fmt"
	"kubernetes/internal/pkg/utils"
	"kubernetes/pkg/schema/generator"
	"kubernetes/pkg/schema/k8s/apps"
	"kubernetes/pkg/schema/k8s/authorization"
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"os"
)

func createVectorManifests(generatorMeta generator.GeneratorMeta) map[string][]byte {
	namespace := utils.ManifestConfig{
		Filename:  "namespace.yaml",
		Manifests: utils.GenerateNamespace(generatorMeta.Namespace),
	}

	serviceAccountName := generatorMeta.Name
	serviceAccount := utils.ManifestConfig{
		Filename: "serviceaccount.yaml",
		Manifests: []any{
			core.NewServiceAccount(meta.ObjectMeta{
				Name: serviceAccountName,
			}),
		},
	}

	config, err := os.ReadFile("./agent.yaml")
	if err != nil {
		fmt.Printf("Error while reading config.yaml")
	}

	configmapName := generatorMeta.Name
	configmap := utils.ManifestConfig{
		Filename: "configmap.yaml",
		Manifests: []any{
			core.NewConfigMap(meta.ObjectMeta{
				Name: configmapName,
			}, map[string]string{
				"agent.yaml": string(config),
			}),
		},
	}

	rbac := utils.ManifestConfig{
		Filename: "rbac.yaml",
		Manifests: []any{
			authorization.NewClusterRole(meta.ObjectMeta{
				Name: generatorMeta.Name,
			}, []authorization.Rule{
				{
					APIGroups: []string{},
					Resources: []string{"namespaces", "nodes", "pods"},
					Verbs:     []string{"list", "watch"},
				},
			}),
			authorization.NewClusterRoleBinding(meta.ObjectMeta{
				Name: generatorMeta.Name,
			}, authorization.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     generatorMeta.Name,
			}, []authorization.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      serviceAccountName,
					Namespace: generatorMeta.Namespace,
				},
			}),
		},
	}

	daemonset := utils.ManifestConfig{
		Filename: "daemonset.yaml",
		Manifests: []any{
			apps.NewDaemonSet(
				meta.ObjectMeta{
					Name: generatorMeta.Name,
					Labels: map[string]string{
						"app.kubernetes.io/name":    generatorMeta.Name,
						"app.kubernetes.io/version": generatorMeta.Docker.Version,
					},
				},
				apps.DaemonSetSpec{
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
							ServiceAccountName: serviceAccountName,
							DNSPolicy:          core.ClusterFirst,
							Containers: []core.Container{
								{
									Name: generatorMeta.Name,
									Args: []string{
										"--config-dir",
										"/etc/vector/",
									},
									Env: []core.Env{
										{
											Name:  "VECTOR_LOG",
											Value: "info",
										},
										{
											Name: "VECTOR_SELF_NODE_NAME",
											ValueFrom: core.ValueFrom{
												FieldRef: core.FieldRef{
													FieldPath: "spec.nodeName",
												},
											},
										},
										{
											Name: "VECTOR_SELF_POD_NAME",
											ValueFrom: core.ValueFrom{
												FieldRef: core.FieldRef{
													FieldPath: "metadata.name",
												},
											},
										},
										{
											Name: "VECTOR_SELF_POD_NAMESPACE",
											ValueFrom: core.ValueFrom{
												FieldRef: core.FieldRef{
													FieldPath: "metadata.namespace",
												},
											},
										},
										{
											Name:  "PROCFS_ROOT",
											Value: "/host/proc",
										},
										{
											Name:  "SYSFS_ROOT",
											Value: "/host/sys",
										},
									},
									Image: fmt.Sprintf("%v:%v", generatorMeta.Docker.Registry, generatorMeta.Docker.Version),
									Ports: []core.Port{
										{
											ContainerPort: 9090,
											Name:          "prom-exporter",
										},
									},
									VolumeMounts: []core.VolumeMount{
										{
											Name:      "data",
											MountPath: "/vector-data-dir",
										},
										{
											Name:      "config",
											MountPath: "/etc/vector/",
										},
										{
											Name:      "var-log",
											MountPath: "/var/log/",
										},
										{
											Name:      "var-lib",
											MountPath: "/var/lib",
										},
										{
											Name:      "procfs",
											MountPath: "/host/proc",
										},
										{
											Name:      "sysfs",
											MountPath: "/host/sys",
										},
									},
								},
							},
							Volumes: []core.Volume{
								{
									Name: "config",
									ConfigMap: core.ConfigMapVolumeSource{
										Name: configmapName,
									},
								},
								{
									Name: "data",
									HostPath: core.HostPath{
										Path: "/var/lib/vector",
									},
								},
								{
									Name: "var-log",
									HostPath: core.HostPath{
										Path: "/var/log/",
									},
								},
								{
									Name: "var-lib",
									HostPath: core.HostPath{
										Path: "/var/lib/",
									},
								},
								{
									Name: "procfs",
									HostPath: core.HostPath{
										Path: "/proc",
									},
								},
								{
									Name: "sysfs",
									HostPath: core.HostPath{
										Path: "/sys",
									},
								},
							},
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
			daemonset.Filename,
			serviceAccount.Filename,
			configmap.Filename,
			rbac.Filename,
		}),
	}

	return utils.MarshalManifests([]utils.ManifestConfig{namespace, kustomization, daemonset, rbac, serviceAccount, configmap})
}
