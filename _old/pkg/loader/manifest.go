package loader

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func loadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	if err := expectKeys("genome.yaml", raw, []string{"project", "schema_version", "traits", "expression"}); err != nil {
		return nil, err
	}
	mv := &Manifest{}

	schemaVersion, ok := raw["schema_version"].(string)
	if !ok || schemaVersion == "" {
		return nil, fmt.Errorf("genome.yaml.schema_version must be a non-empty string")
	}
	mv.SchemaVersion = schemaVersion

	projectRaw, ok := raw["project"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("genome.yaml.project must be an object")
	}
	if err := expectKeys("project", projectRaw, []string{"name", "description", "type"}); err != nil {
		return nil, err
	}
	name, ok := projectRaw["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("project.name must be a non-empty string")
	}
	project := Project{
		Name: name,
	}
	if d, ok := projectRaw["description"].(string); ok {
		project.Description = d
	}
	if t, ok := projectRaw["type"].(string); ok {
		project.Type = t
	}
	mv.Project = project

	if traitsRaw, ok := raw["traits"]; ok {
		traits, err := toStringList("traits", traitsRaw)
		if err != nil {
			return nil, err
		}
		mv.Traits = traits
	}

	if exprRaw, ok := raw["expression"]; ok {
		if expr, ok := exprRaw.(map[string]any); ok {
			mv.Expression = expr
		} else {
			return nil, fmt.Errorf("expression must be an object when present")
		}
	}

	return mv, nil
}
