package genome

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadGenome loads genome.yaml and all referenced module files from rootDir.
func LoadGenome(rootDir string) (*Genome, error) {
	manifest, err := LoadManifest(rootDir)
	if err != nil {
		return nil, err
	}

	g := &Genome{
		Manifest: *manifest,
		Modules:  make(map[string]Module, len(manifest.Modules)),
	}

	for moduleName, relPath := range manifest.Modules {
		modulePath := filepath.Join(rootDir, relPath)

		moduleBytes, err := os.ReadFile(modulePath)
		if err != nil {
			return nil, fmt.Errorf("read module %q at %q: %w", moduleName, modulePath, err)
		}

		var mf ModuleFile
		if err := unmarshalStrict(moduleBytes, &mf); err != nil {
			return nil, fmt.Errorf("parse module %q at %q: %w", moduleName, modulePath, err)
		}

		if len(mf.Module) != 1 {
			return nil, fmt.Errorf("module file %q must define exactly one module", modulePath)
		}

		var actualName string
		var mod Module
		for name, m := range mf.Module {
			actualName = name
			mod = m
		}

		if actualName != moduleName {
			return nil, fmt.Errorf(
				"module registry mismatch: manifest says %q but file %q defines %q",
				moduleName, modulePath, actualName,
			)
		}

		g.Modules[moduleName] = mod
	}

	return g, nil
}
