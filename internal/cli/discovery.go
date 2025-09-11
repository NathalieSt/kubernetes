package cli

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/rivo/tview"
)

func discoverGenerators(generatorsLocation string, logOutputView *tview.TextView) {
	err := filepath.WalkDir(generatorsLocation, func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path, d.Name(), "directory?", d.IsDir())

		if !d.IsDir() {
			if d.Name() == "main.go" {
				info, err := d.Info()
				logToOutput(logOutputView, fmt.Sprintf("Generator found: %v %v\n", info.Sys(), err))
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("impossible to walk directories: %s", err)
	}
}
