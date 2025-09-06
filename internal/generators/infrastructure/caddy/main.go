package main

import "kubernetes/internal/pkg/utils"

func main() {
	//fmt.Println("✅ Finding project root")
	//rootDir, err := utils.FindRoot()
	//if err != nil {
	//	fmt.Println("❌ An error occurred while finding the project root")
	//	fmt.Println("Error: " + err.Error())
	//}
	defer utils.Timer()()
	getCaddyConfigMap()

	//utils.RunGenerator(utils.GeneratorConfig{
	//		Meta:            Forgejo,
	//		OutputDir:       filepath.Join(rootDir, "/cluster/apps/forgejo/"),
	//		CreateManifests: createForgejoManifests,
	//	})
}
