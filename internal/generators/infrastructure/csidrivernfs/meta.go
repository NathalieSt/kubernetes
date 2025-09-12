package main

import (
	"kubernetes/pkg/schema/generator"
)

var CSIDriverNFS = generator.GeneratorMeta{
	Name:          "csi-driver-nfs",
	Namespace:     "csi-driver-nfs",
	GeneratorType: generator.Infrastructure,
	Helm: &generator.Helm{
		Chart:   "csi-driver-nfs",
		Url:     "https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/charts",
		Version: "4.11.0",
	},
	DependsOnGenerators: []string{},
}
