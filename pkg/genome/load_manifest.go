package genome

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func unmarshalStrict(data []byte, out any) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	return dec.Decode(out)
}

// LoadManifest loads and validates genome.yaml from rootDir.
func LoadManifest(rootDir string) (*Manifest, error) {
	manifestPath := filepath.Join(rootDir, "genome.yaml")

	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read manifest %q: %w", manifestPath, err)
	}

	var manifest Manifest
	if err := unmarshalStrict(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest %q: %w", manifestPath, err)
	}

	if manifest.Project.Name == "" {
		return nil, fmt.Errorf("genome missing project.name")
	}

	if len(manifest.Modules) == 0 {
		return nil, fmt.Errorf("genome missing modules")
	}

	seenPaths := make(map[string]string, len(manifest.Modules))
	for moduleName, relPath := range manifest.Modules {
		if firstModule, exists := seenPaths[relPath]; exists {
			return nil, fmt.Errorf(
				"multiple modules reference the same module file path %q: %q and %q",
				relPath, firstModule, moduleName,
			)
		}
		seenPaths[relPath] = moduleName
	}

	return &manifest, nil
}
