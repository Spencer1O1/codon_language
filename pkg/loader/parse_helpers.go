package loader

import (
	"fmt"
)

func parseEntities(raw any) ([]ComposedEntity, error) {
	if raw == nil {
		return nil, nil
	}
	asMap, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("entities must be a mapping of identifier -> definition")
	}
	var out []ComposedEntity
	for name, defRaw := range asMap {
		if err := validateIdentifier("entity", name); err != nil {
			return nil, err
		}
		defMap := map[string]any{}
		if defRaw != nil {
			if m, ok := defRaw.(map[string]any); ok {
				defMap = m
			} else {
				return nil, fmt.Errorf("entity %q definition must be an object", name)
			}
		}
		fields, err := parseFieldMap(defMap["fields"])
		if err != nil {
			return nil, fmt.Errorf("entity %q: %w", name, err)
		}
		out = append(out, ComposedEntity{
			Name:   name,
			Fields: fields,
		})
	}
	return out, nil
}

func parseCapabilities(raw any) ([]ComposedCapability, error) {
	if raw == nil {
		return nil, nil
	}
	switch v := raw.(type) {
	case []any:
		var out []ComposedCapability
		for i, item := range v {
			name, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("capabilities[%d] must be a string or structured map", i)
			}
			if err := validateIdentifier("capability", name); err != nil {
				return nil, err
			}
			out = append(out, ComposedCapability{Name: name})
		}
		return out, nil
	case map[string]any:
		var out []ComposedCapability
		for name, rawDef := range v {
			if err := validateIdentifier("capability", name); err != nil {
				return nil, err
			}
			defMap := map[string]any{}
			if rawDef != nil {
				m, ok := rawDef.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("capability %q definition must be an object", name)
				}
				defMap = m
			}
			cap := ComposedCapability{Name: name}
			if desc, ok := defMap["description"].(string); ok {
				cap.Description = desc
			}
			inputs, err := parseFieldMap(defMap["inputs"])
			if err != nil {
				return nil, fmt.Errorf("capability %q inputs: %w", name, err)
			}
			cap.Inputs = inputs

			outputs, err := parseFieldMap(defMap["outputs"])
			if err != nil {
				return nil, fmt.Errorf("capability %q outputs: %w", name, err)
			}
			cap.Outputs = outputs

			if effectsRaw, ok := defMap["effects"]; ok {
				effects, err := toStringList("effects", effectsRaw)
				if err != nil {
					return nil, fmt.Errorf("capability %q: %w", name, err)
				}
				cap.Effects = effects
			}
			out = append(out, cap)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("capabilities must be a list of identifiers or a map of structured definitions")
	}
}

func parseFieldMap(raw any) (map[string]FieldDefinition, error) {
	if raw == nil {
		return nil, nil
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("fields must be an object")
	}
	out := make(map[string]FieldDefinition, len(m))
	for name, defRaw := range m {
		fd, err := parseFieldDefinition(defRaw)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", name, err)
		}
		out[name] = fd
	}
	return out, nil
}

func parseFieldDefinition(raw any) (FieldDefinition, error) {
	// shorthand: string
	if t, ok := raw.(string); ok {
		if err := validateFieldType(t); err != nil {
			return FieldDefinition{}, err
		}
		return FieldDefinition{Type: t}, nil
	}

	m, ok := raw.(map[string]any)
	if !ok {
		return FieldDefinition{}, fmt.Errorf("field definition must be string shorthand or structured object")
	}

	fd := FieldDefinition{}

	typeVal, ok := m["type"].(string)
	if !ok || typeVal == "" {
		return FieldDefinition{}, fmt.Errorf("structured field must include type")
	}
	if err := validateFieldType(typeVal); err != nil {
		return FieldDefinition{}, err
	}
	fd.Type = typeVal

	if opt, ok := m["optional"].(bool); ok {
		fd.Optional = opt
	}
	if uniq, ok := m["unique"].(bool); ok {
		fd.Unique = uniq
	}
	if defv, ok := m["default"]; ok {
		fd.Default = defv
	}
	if valsRaw, ok := m["values"]; ok {
		values, err := toStringList("values", valsRaw)
		if err != nil {
			return FieldDefinition{}, err
		}
		fd.Values = values
	}
	if ref, ok := m["reference"].(string); ok {
		if err := validateEntityReference(ref); err != nil {
			return FieldDefinition{}, err
		}
		fd.Reference = ref
	}

	// Constraints from spec:
	if fd.Type == "enum" && len(fd.Values) == 0 {
		return FieldDefinition{}, fmt.Errorf("enum fields must define values")
	}
	if fd.Type != "enum" && len(fd.Values) > 0 {
		return FieldDefinition{}, fmt.Errorf("values allowed only for enum fields")
	}
	if fd.Type == "reference" && fd.Reference == "" {
		return FieldDefinition{}, fmt.Errorf("reference fields must provide reference")
	}
	if fd.Type != "reference" && fd.Reference != "" {
		return FieldDefinition{}, fmt.Errorf("reference may only be present when type is reference")
	}

	return fd, nil
}
