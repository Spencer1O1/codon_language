package nucleotype

// Resolve flattens name indirections and treats "primitive" as a terminal.
// It is recursion-safe: self-references return the name as-is rather than looping.
func Resolve(t TypeNode, symbols map[string]TypeNode) TypeNode {
	return resolve(t, symbols, map[string]bool{})
}

func resolve(t TypeNode, symbols map[string]TypeNode, seen map[string]bool) TypeNode {
	switch v := t.(type) {
	case NameType:
		if v.Name == "primitive" {
			return v
		}
		if target, ok := symbols[v.Name]; ok {
			if seen[v.Name] {
				return v // break cycles; leave as reference
			}
			seen[v.Name] = true
			resolved := resolve(target, symbols, seen)
			delete(seen, v.Name)
			return resolved
		}
		return v
	case OptionalType:
		return OptionalType{Base: resolve(v.Base, symbols, seen)}
	case ListType:
		return ListType{Base: resolve(v.Base, symbols, seen)}
	case UnionType:
		opts := make([]TypeNode, len(v.Options))
		for i, o := range v.Options {
			opts[i] = resolve(o, symbols, seen)
		}
		return UnionType{Options: opts}
	case GenericType:
		args := make([]TypeNode, len(v.Args))
		for i, a := range v.Args {
			args[i] = resolve(a, symbols, seen)
		}
		return GenericType{Name: v.Name, Args: args}
	case ObjectType:
		fields := make([]Field, len(v.Fields))
		for i, f := range v.Fields {
			fields[i] = Field{Name: f.Name, Type: resolve(f.Type, symbols, seen)}
		}
		return ObjectType{Fields: fields}
	default:
		return t
	}
}
