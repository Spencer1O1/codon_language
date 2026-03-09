package rules

import (
	"fmt"
	"strings"

	"github.com/Spencer1O1/codon-language/internal/domain/genome"
	"github.com/Spencer1O1/codon-language/internal/domain/validation"
)

type NoCircularDependenciesRule struct{}

func (NoCircularDependenciesRule) Name() string {
	return "no-circular-dependencies"
}

func (NoCircularDependenciesRule) Validate(g *genome.Genome) []validation.Finding {
	var findings []validation.Finding

	const (
		unvisited = 0
		visiting  = 1
		visited   = 2
	)

	state := make(map[string]int, len(g.Modules))
	var stack []string
	seenCycles := make(map[string]struct{})

	var dfs func(string)
	dfs = func(moduleName string) {
		state[moduleName] = visiting
		stack = append(stack, moduleName)

		mod := g.Modules[moduleName]
		for _, dep := range mod.Dependencies {
			// Let the "dependency exists" rule handle missing modules.
			if _, ok := g.Modules[dep]; !ok {
				continue
			}

			switch state[dep] {
			case unvisited:
				dfs(dep)
			case visiting:
				// Found a cycle. Extract the cycle path from the stack.
				cycleStart := indexOf(stack, dep)
				if cycleStart >= 0 {
					cyclePath := append([]string{}, stack[cycleStart:]...)
					cyclePath = append(cyclePath, dep)

					cycleKey := strings.Join(cyclePath, "->")
					if _, alreadyReported := seenCycles[cycleKey]; !alreadyReported {
						seenCycles[cycleKey] = struct{}{}

						findings = append(findings, validation.Finding{
							Severity: validation.SeverityError,
							Code:     "circular_dependency",
							Path:     fmt.Sprintf("module/%s/dependencies", moduleName),
							Message:  fmt.Sprintf("circular dependency detected: %s", strings.Join(cyclePath, " -> ")),
						})
					}
				}
			case visited:
				// noop
			}
		}

		stack = stack[:len(stack)-1]
		state[moduleName] = visited
	}

	for moduleName := range g.Modules {
		if state[moduleName] == unvisited {
			dfs(moduleName)
		}
	}

	return findings
}

func indexOf(items []string, target string) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return -1
}
