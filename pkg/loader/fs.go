package loader

import (
	"io/fs"
	"os"
	"path/filepath"
)

// LoadFS variant to ease testing with an fs.FS (e.g., embed).
func LoadFS(fsys fs.FS, root string) (cg *ComposedGenome, err error) {
	tmpDir, err := os.MkdirTemp("", "codon-fs-*")
	if err != nil {
		return nil, err
	}
	defer func() {
		if rmErr := os.RemoveAll(tmpDir); rmErr != nil && err == nil {
			// propagate cleanup failure only if no earlier error
			err = rmErr
		}
	}()
	if err = copyFS(fsys, root, tmpDir); err != nil {
		return nil, err
	}
	return Load(tmpDir)
}

func copyFS(fsys fs.FS, root, dst string) error {
	return fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}
