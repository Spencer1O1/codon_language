package loader

import "fmt"

func toStringList(label string, raw any) ([]string, error) {
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a list of strings", label)
	}
	out := make([]string, 0, len(list))
	for i, item := range list {
		s, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%d] must be a string", label, i)
		}
		out = append(out, s)
	}
	return out, nil
}

func toStringListOptional(label string, raw any) ([]string, error) {
	if raw == nil {
		return nil, nil
	}
	return toStringList(label, raw)
}
