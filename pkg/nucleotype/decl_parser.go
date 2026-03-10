package nucleotype

import (
	"bufio"
	"fmt"
	"strings"
)

// Decl represents a single declaration (with optional export).
type Decl struct {
	Export bool
	Name   string
	Type   TypeNode
}

// ParseDecls parses declaration files with syntax:
//
//	[export ]name = <type expression>
//
// Type expressions can span multiple lines; parsing stops at end of file.
func ParseDecls(src string) ([]Decl, error) {
	var decls []Decl
	sc := bufio.NewScanner(strings.NewReader(src))
	var buf strings.Builder
	var export bool
	var name string
	readingType := false
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !readingType {
			// header line: [export ]name =
			export = false
			if strings.HasPrefix(line, "export ") {
				export = true
				line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("expected '=' in declaration line: %q", line)
			}
			name = strings.TrimSpace(parts[0])
			buf.Reset()
			buf.WriteString(strings.TrimSpace(parts[1]))
			if balanced(buf.String()) {
				readingType = false
				t, err := Parse(buf.String())
				if err != nil {
					return nil, err
				}
				decls = append(decls, Decl{Export: export, Name: name, Type: t})
			} else {
				readingType = true
			}
		} else {
			// continue accumulating type until balanced
			if buf.Len() > 0 {
				buf.WriteString("\n")
			}
			buf.WriteString(line)
			if balanced(buf.String()) {
				readingType = false
				t, err := Parse(buf.String())
				if err != nil {
					return nil, err
				}
				decls = append(decls, Decl{Export: export, Name: name, Type: t})
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if readingType {
		return nil, fmt.Errorf("unterminated type for %s", name)
	}
	return decls, nil
}

// balanced checks simple delimiter balance for (), {}, <>.
func balanced(s string) bool {
	var stack []rune
	for _, r := range s {
		switch r {
		case '(', '{', '<':
			stack = append(stack, r)
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '>':
			if len(stack) == 0 || stack[len(stack)-1] != '<' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}
