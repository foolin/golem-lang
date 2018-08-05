// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package graph

import (
	"fmt"
)

type (
	// Node is a node in a Directed Acyclic Graph
	Node struct {
		// The name of this Node
		Name string
		// The names of this Node's children
		Children []string
	}

	entryStatus int

	entry struct {
		node   Node
		status entryStatus
	}
)

const (
	unmarked entryStatus = iota
	temporary
	permanent
)

func visit(e *entry, entryMap map[string]*entry, result []*entry) ([]*entry, error) {

	// if the entry is already permanent, then return
	if e.status == permanent {
		return result, nil
	}

	// if the entry is temporary, then we have found a cyclic graph
	if e.status == temporary {
		return nil, fmt.Errorf("Graph cycle detected on node '%s'", e.node.Name)
	}

	// mark temporarily
	e.status = temporary

	// visit each child
	for _, cn := range e.node.Children {

		ce, ok := entryMap[cn]
		if !ok {
			return nil, fmt.Errorf("Cannot find node '%s'", cn)
		}

		var err error
		result, err = visit(ce, entryMap, result)
		if err != nil {
			return nil, err
		}
	}

	// mark permanently
	e.status = permanent

	// add to result
	return append(result, e), nil
}

// TopologicalSort does a 'topological sort' on the Nodes in a
// Directed Acyclic Graph.  See https://en.wikipedia.org/wiki/Topological_sorting.
// Note: the resulting sort order will not necessarily be unique.
func TopologicalSort(nodes []Node) ([]Node, error) {

	// create map of entries
	entryMap := make(map[string]*entry)
	for _, n := range nodes {
		if _, ok := entryMap[n.Name]; ok {
			return nil, fmt.Errorf("Duplicate node name: '%s'", n.Name)
		}
		entryMap[n.Name] = &entry{node: n, status: unmarked}
	}

	// visit unmarked nodes
	result := []*entry{}
	for _, n := range nodes {

		e := entryMap[n.Name]
		if e.status == unmarked {
			var err error
			result, err = visit(e, entryMap, result)
			if err != nil {
				return nil, err
			}
		}
	}

	// done
	sorted := []Node{}
	for i := len(result) - 1; i >= 0; i-- {
		sorted = append(sorted, result[i].node)
	}
	return sorted, nil
}
