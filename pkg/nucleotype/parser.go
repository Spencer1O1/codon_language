package nucleotype

import (
	"fmt"
	"strings"
	"unicode"
)

// Parse parses a type expression into an AST.
func Parse(input string) (TypeNode, error) {
	p := &parser{src: input}
	p.next()
	t, err := p.parseType()
	if err != nil {
		return nil, err
	}
	if p.tok.typ != tokEOF {
		return nil, fmt.Errorf("unexpected token %q at %d", p.tok.val, p.pos)
	}
	return t, nil
}

// tokenizer
type tokenType int

const (
	tokEOF tokenType = iota
	tokIdent
	tokString
	tokSymbol
)

type token struct {
	typ tokenType
	val string
}

type parser struct {
	src string
	pos int
	tok token
}

func (p *parser) next() {
	p.skipWS()
	if p.pos >= len(p.src) {
		p.tok = token{typ: tokEOF}
		return
	}
	ch := p.src[p.pos]
	switch {
	case isIdentStart(ch):
		start := p.pos
		p.pos++
		for p.pos < len(p.src) && isIdentPart(p.src[p.pos]) {
			p.pos++
		}
		p.tok = token{typ: tokIdent, val: p.src[start:p.pos]}
	case ch == '"':
		p.pos++
		start := p.pos
		for p.pos < len(p.src) && p.src[p.pos] != '"' {
			p.pos++
		}
		p.tok = token{typ: tokString, val: p.src[start:p.pos]}
		if p.pos < len(p.src) && p.src[p.pos] == '"' {
			p.pos++
		}
	default:
		p.pos++
		p.tok = token{typ: tokSymbol, val: string(ch)}
	}
	p.skipWS()
}

func (p *parser) skipWS() {
	for p.pos < len(p.src) {
		ch := p.src[p.pos]
		if unicode.IsSpace(rune(ch)) {
			p.pos++
			continue
		}
		// skip comments starting with #
		if ch == '#' {
			for p.pos < len(p.src) && p.src[p.pos] != '\n' {
				p.pos++
			}
			continue
		}
		break
	}
}

func isIdentStart(b byte) bool {
	return unicode.IsLetter(rune(b)) || b == '_'
}
func isIdentPart(b byte) bool {
	return isIdentStart(b) || unicode.IsDigit(rune(b)) || b == '-'
}

// grammar: TYPE := UNION ('|' UNION)*
func (p *parser) parseType() (TypeNode, error) {
	left, err := p.parsePostfix()
	if err != nil {
		return nil, err
	}
	var opts []TypeNode
	for p.tok.typ == tokSymbol && p.tok.val == "|" {
		opts = append(opts, left)
		p.next()
		right, err := p.parsePostfix()
		if err != nil {
			return nil, err
		}
		left = right
	}
	if len(opts) == 0 {
		return left, nil
	}
	opts = append(opts, left)
	return UnionType{Options: opts}, nil
}

// POSTFIX := PRIMARY ( '?' | '[]' )*
func (p *parser) parsePostfix() (TypeNode, error) {
	base, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		if p.tok.typ == tokSymbol && p.tok.val == "?" {
			base = OptionalType{Base: base}
			p.next()
			continue
		}
		if p.tok.typ == tokSymbol && p.tok.val == "[" && p.peekSymbol("]") {
			// consume []
			p.next() // [
			p.next() // ]
			base = ListType{Base: base}
			continue
		}
		break
	}
	return base, nil
}

// PRIMARY := ident | ident '<' type (, type)* '>' | '{' fields '}' | '(' type ')'
func (p *parser) parsePrimary() (TypeNode, error) {
	switch p.tok.typ {
	case tokIdent:
		name := p.tok.val
		p.next()
		if p.tok.typ == tokSymbol && p.tok.val == "<" {
			p.next() // consume <
			var args []TypeNode
			for {
				t, err := p.parseType()
				if err != nil {
					return nil, err
				}
				args = append(args, t)
				if p.tok.typ == tokSymbol && p.tok.val == ">" {
					p.next()
					break
				}
				if p.tok.typ == tokSymbol && p.tok.val == "," {
					p.next()
					continue
				}
				return nil, fmt.Errorf("expected ',' or '>' in generic at %d", p.pos)
			}
			return GenericType{Name: name, Args: args}, nil
		}
		return NameType{Name: name}, nil
	case tokString:
		val := p.tok.val
		p.next()
		return LiteralType{Value: val}, nil
	case tokSymbol:
		if p.tok.val == "(" {
			p.next()
			t, err := p.parseType()
			if err != nil {
				return nil, err
			}
			if p.tok.typ != tokSymbol || p.tok.val != ")" {
				return nil, fmt.Errorf("expected ')' at %d", p.pos)
			}
			p.next()
			return t, nil
		}
		if p.tok.val == "{" {
			return p.parseObject()
		}
	}
	return nil, fmt.Errorf("unexpected token %q at %d", p.tok.val, p.pos)
}

func (p *parser) parseObject() (TypeNode, error) {
	// assume current token is "{"
	p.next()
	var fields []Field
	for {
		if p.tok.typ == tokSymbol && p.tok.val == "}" {
			p.next()
			break
		}
		if p.tok.typ != tokIdent {
			return nil, fmt.Errorf("expected field name at %d", p.pos)
		}
		name := p.tok.val
		p.next()
		if p.tok.typ != tokSymbol || p.tok.val != ":" {
			return nil, fmt.Errorf("expected ':' after field name at %d", p.pos)
		}
		p.next()
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		fields = append(fields, Field{Name: name, Type: t})
		// optional trailing comma handled by skipWS
		if p.tok.typ == tokSymbol && p.tok.val == "}" {
			continue
		}
		if p.tok.typ == tokSymbol && p.tok.val == "," {
			p.next()
			continue
		}
	}
	return ObjectType{Fields: fields}, nil
}

func (p *parser) peekSymbol(val string) bool {
	saveTok := p.tok
	savePos := p.pos
	p.next()
	ok := p.tok.typ == tokSymbol && p.tok.val == val
	p.tok = saveTok
	p.pos = savePos
	return ok
}

// Pretty renders a TypeNode to a string (for debugging/tests).
func Pretty(t TypeNode) string {
	switch v := t.(type) {
	case NameType:
		return v.Name
	case GenericType:
		var args []string
		for _, a := range v.Args {
			args = append(args, Pretty(a))
		}
		return fmt.Sprintf("%s<%s>", v.Name, strings.Join(args, ", "))
	case LiteralType:
		return fmt.Sprintf("\"%s\"", v.Value)
	case OptionalType:
		return Pretty(v.Base) + "?"
	case ListType:
		return Pretty(v.Base) + "[]"
	case UnionType:
		var opts []string
		for _, o := range v.Options {
			opts = append(opts, Pretty(o))
		}
		return strings.Join(opts, " | ")
	case ObjectType:
		var parts []string
		for _, f := range v.Fields {
			parts = append(parts, fmt.Sprintf("%s: %s", f.Name, Pretty(f.Type)))
		}
		return "{ " + strings.Join(parts, ", ") + " }"
	default:
		return "<unknown>"
	}
}
