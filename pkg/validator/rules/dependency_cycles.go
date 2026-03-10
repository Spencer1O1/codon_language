package rules

import (
	"fmt"
	"strings"

	"github.com/Spencer1O1/codon-language/pkg/loader"
	"github.com/Spencer1O1/codon-language/pkg/validator/core"
)

func init() {
	core.Register(checkDependencyCycles)
}

func checkDependencyCycles(genome *loader.ComposedGenome, res *core.Result) {
	graph := make(map[string][]string)
	for _, g := range genome.Genes {
		key := g.Chromosome + "." + g.Name
		graph[key] = append(graph[key], g.Dependencies...)
	}
	visited := make(map[string]bool)
	stack := make(map[string]bool)

	var dfs func(string, []string)
	dfs = func(node string, path []string) {
		if stack[node] {
			// found cycle
			cycle := append(path, node)
			res.Add("dependencies", fmt.Sprintf("dependency cycle detected: %s", strings.Join(cycle, " -> ")))
			return
		}
		if visited[node] {
			return
		}
		visited[node] = true
		stack[node] = true
		for _, neigh := range graph[node] {
			dfs(neigh, append(path, node))
		}
		stack[node] = false
	}

	for node := range graph {
		if !visited[node] {
			dfs(node, nil)
		}
	}
}
