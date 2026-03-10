package nucleotype

// Resolver flattens name indirections and treats "primitive" as a terminal.
// Symbols is a map of name -> TypeNode (e.g., from parsed declarations).
func Resolve(t TypeNode, symbols map[string]TypeNode) TypeNode {
	switch v := t.(type) {
	case NameType:
		if v.Name == "primitive" {
			return v
		}
		if target, ok := symbols[v.Name]; ok {
			return Resolve(target, symbols)
		}
		return v
	case OptionalType:
		return OptionalType{Base: Resolve(v.Base, symbols)}
	case ListType:
		return ListType{Base: Resolve(v.Base, symbols)}
	case UnionType:
		opts := make([]TypeNode, len(v.Options))
		for i, o := range v.Options {
			opts[i] = Resolve(o, symbols)
		}
		return UnionType{Options: opts}
	case GenericType:
		args := make([]TypeNode, len(v.Args))
		for i, a := range v.Args {
			args[i] = Resolve(a, symbols)
		}
		return GenericType{Name: v.Name, Args: args}
	case ObjectType:
		fields := make([]Field, len(v.Fields))
		for i, f := range v.Fields {
			fields[i] = Field{Name: f.Name, Type: Resolve(f.Type, symbols)}
		}
		return ObjectType{Fields: fields}
	default:
		return t
	}
}
