package uws1

import (
	"fmt"
	"sort"
	"strings"
)

// detectDependencyCycles walks the dependsOn graph and reports any cycles.
// Parallel-group dependencies fan out to each member of the group (excluding
// the node itself) so that indirect self-dependencies surface too. Unknown
// dependency targets are ignored here because validateDependencyList already
// flags them.
func detectDependencyCycles(idx *documentIndex, result *ValidationResult) {
	if len(idx.dependencies) == 0 {
		return
	}

	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make(map[string]int)
	seen := make(map[string]bool)

	neighbors := func(node string) []string {
		deps := idx.dependencies[node]
		if len(deps) == 0 {
			return nil
		}
		out := make([]string, 0, len(deps))
		for _, d := range deps {
			if idx.operations[d] != nil || idx.workflows[d] != nil || idx.steps[d] != nil {
				out = append(out, d)
				continue
			}
			if members, ok := idx.parallelGroupMembers[d]; ok {
				for _, m := range members {
					if m != node {
						out = append(out, m)
					}
				}
			}
		}
		return out
	}

	sources := make([]string, 0, len(idx.dependencies))
	for k := range idx.dependencies {
		sources = append(sources, k)
	}
	sort.Strings(sources)
	for _, s := range sources {
		if color[s] != white {
			continue
		}
		type frame struct {
			node      string
			neighbors []string
			next      int
		}
		stack := []frame{{node: s, neighbors: neighbors(s)}}
		path := []string{s}
		color[s] = gray
		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			if top.next >= len(top.neighbors) {
				color[top.node] = black
				stack = stack[:len(stack)-1]
				path = path[:len(path)-1]
				continue
			}
			nb := top.neighbors[top.next]
			top.next++
			switch color[nb] {
			case white:
				color[nb] = gray
				stack = append(stack, frame{node: nb, neighbors: neighbors(nb)})
				path = append(path, nb)
			case gray:
				start := -1
				for i, n := range path {
					if n == nb {
						start = i
						break
					}
				}
				if start < 0 {
					continue
				}
				cycle := append([]string(nil), path[start:]...)
				key := canonicalCycleKey(cycle)
				if seen[key] {
					continue
				}
				seen[key] = true
				cycle = append(cycle, nb)
				result.addError("dependsOn", fmt.Sprintf("cycle detected: %s", strings.Join(cycle, " -> ")))
			}
		}
	}
}

func canonicalCycleKey(cycle []string) string {
	if len(cycle) == 0 {
		return ""
	}
	minIdx := 0
	for i, n := range cycle {
		if n < cycle[minIdx] {
			minIdx = i
		}
	}
	rotated := make([]string, 0, len(cycle))
	rotated = append(rotated, cycle[minIdx:]...)
	rotated = append(rotated, cycle[:minIdx]...)
	return strings.Join(rotated, "->")
}

// detectWorkflowRecursion rejects direct and indirect workflow-reference
// cycles. Workflow calls are synchronous, so even a cycle without dependsOn
// edges would otherwise re-enter an in-flight workflow and deadlock.
func detectWorkflowRecursion(idx *documentIndex, result *ValidationResult) {
	const (
		unvisited = iota
		visiting
		visited
	)
	state := make(map[string]int, len(idx.workflows))
	stack := make([]string, 0, len(idx.workflows))
	reported := make(map[string]bool)

	var visit func(string)
	visit = func(workflowID string) {
		state[workflowID] = visiting
		stack = append(stack, workflowID)
		for _, call := range idx.workflowCalls[workflowID] {
			if idx.workflows[call.target] == nil {
				continue
			}
			switch state[call.target] {
			case unvisited:
				visit(call.target)
			case visiting:
				start := 0
				for i, id := range stack {
					if id == call.target {
						start = i
						break
					}
				}
				cycle := append([]string(nil), stack[start:]...)
				key := canonicalCycleKey(cycle)
				if reported[key] {
					continue
				}
				reported[key] = true
				cycle = append(cycle, call.target)
				result.addError(call.path, fmt.Sprintf("recursive workflow call detected: %s", strings.Join(cycle, " -> ")))
			}
		}
		stack = stack[:len(stack)-1]
		state[workflowID] = visited
	}

	workflowIDs := make([]string, 0, len(idx.workflows))
	for id := range idx.workflows {
		workflowIDs = append(workflowIDs, id)
	}
	sort.Strings(workflowIDs)
	for _, id := range workflowIDs {
		if state[id] == unvisited {
			visit(id)
		}
	}
}
