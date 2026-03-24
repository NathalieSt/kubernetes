package main

import (
	"flag"
	"fmt"
	"io/fs"
	"kubernetes/internal/pkg/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func findGeneratorPaths(generatorsDir string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(generatorsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "main.go" {
			paths = append(paths, strings.TrimSuffix(path, "main.go"))
		}
		return nil
	})
	return paths, err
}

func main() {
	rootDir := flag.String("root", "", "Root directory of the generator project")
	flag.Parse()

	if *rootDir == "" {
		var err error
		*rootDir, err = utils.FindRoot()
		if err != nil {
			fmt.Println("❌ Could not find root directory (no go.mod found)")
			os.Exit(1)
		}
	}

	generatorsDir := filepath.Join(*rootDir, "internal/generators")
	paths, err := findGeneratorPaths(generatorsDir)
	if err != nil {
		fmt.Printf("❌ Failed to walk generators directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d potential generator locations\n", len(paths))

	var wg sync.WaitGroup
	var mu sync.Mutex
	failures := []string{}

	for _, path := range paths {
		wg.Go(func() {
			meta, err := utils.GetGeneratorMeta(*rootDir, path)
			if err != nil || meta == nil {
				return
			}

			fmt.Printf("▶ Running generator: %v\n", meta.Name)
			rootFlag := fmt.Sprintf("--root=%v", *rootDir)
			_, err = utils.RunGeneratorMain(path, []string{rootFlag})
			if err != nil {
				mu.Lock()
				failures = append(failures, meta.Name)
				mu.Unlock()
				fmt.Printf("❌ Failed: %v\n", meta.Name)
			} else {
				fmt.Printf("✅ Done: %v\n", meta.Name)
			}
		})
	}
	wg.Wait()

	if len(failures) > 0 {
		fmt.Printf("\n❌ %d generator(s) failed: %v\n", len(failures), failures)
		os.Exit(1)
	}

	fmt.Println("\n✅ All generators completed successfully")
}
