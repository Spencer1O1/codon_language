package loader

import "fmt"

func parseRelations(raw any) ([]RelationDefinition, error) {
	if raw == nil {
		return nil, nil
	}
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("relations must be a list")
	}
	out := make([]RelationDefinition, 0, len(list))
	for i, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("relations[%d] must be an object", i)
		}
		var rd RelationDefinition
		from, ok := m["from"].(string)
		if !ok || from == "" {
			return nil, fmt.Errorf("relations[%d].from must be a string", i)
		}
		if err := validateIdentifier("entity", from); err != nil {
			return nil, fmt.Errorf("relations[%d].from: %w", i, err)
		}
		rd.From = from

		to, ok := m["to"].(string)
		if !ok || to == "" {
			return nil, fmt.Errorf("relations[%d].to must be a string", i)
		}
		if err := validateIdentifier("entity", to); err != nil {
			return nil, fmt.Errorf("relations[%d].to: %w", i, err)
		}
		rd.To = to

		t, ok := m["type"].(string)
		if !ok || t == "" {
			return nil, fmt.Errorf("relations[%d].type must be a string", i)
		}
		if err := validateRelationType(t); err != nil {
			return nil, fmt.Errorf("relations[%d].type: %w", i, err)
		}
		rd.Type = t

		name, ok := m["name"].(string)
		if !ok || name == "" {
			return nil, fmt.Errorf("relations[%d].name must be a string", i)
		}
		rd.Name = name

		out = append(out, rd)
	}
	return out, nil
}

func parseReferences(raw any) ([]ReferenceDefinition, error) {
	if raw == nil {
		return nil, nil
	}
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("references must be a list")
	}
	out := make([]ReferenceDefinition, 0, len(list))
	for i, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("references[%d] must be an object", i)
		}
		var rd ReferenceDefinition
		from, ok := m["from"].(string)
		if !ok || from == "" {
			return nil, fmt.Errorf("references[%d].from must be a string", i)
		}
		if err := validateIdentifier("entity", from); err != nil {
			return nil, fmt.Errorf("references[%d].from: %w", i, err)
		}
		rd.From = from

		to, ok := m["to"].(string)
		if !ok || to == "" {
			return nil, fmt.Errorf("references[%d].to must be a string", i)
		}
		if err := validateEntityReference(to); err != nil {
			return nil, fmt.Errorf("references[%d].to: %w", i, err)
		}
		rd.To = to

		t, ok := m["type"].(string)
		if !ok || t == "" {
			return nil, fmt.Errorf("references[%d].type must be a string", i)
		}
		if err := validateRelationType(t); err != nil {
			return nil, fmt.Errorf("references[%d].type: %w", i, err)
		}
		rd.Type = t

		name, ok := m["name"].(string)
		if !ok || name == "" {
			return nil, fmt.Errorf("references[%d].name must be a string", i)
		}
		rd.Name = name

		refVal := "id"
		if refRaw, ok := m["reference"]; ok {
			if s, ok := refRaw.(string); ok && s != "" {
				refVal = s
			} else {
				return nil, fmt.Errorf("references[%d].reference must be a string when present", i)
			}
		}
		rd.Reference = refVal

		out = append(out, rd)
	}
	return out, nil
}
