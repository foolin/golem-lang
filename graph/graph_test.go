// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package graph

import (
	"reflect"
	"testing"
)

func names(nodes []Node) []string {
	names := []string{}
	for _, n := range nodes {
		names = append(names, n.Name)
	}
	return names
}

func ok(t *testing.T, sorted []Node, err error, expect []string) {

	if sorted == nil {
		t.Error("sorted == nil")
		return
	}

	if err != nil {
		t.Error("err != nil")
	}

	if !reflect.DeepEqual(names(sorted), expect) {
		t.Error(names(sorted), " != ", expect)
	}
}

func fail(t *testing.T, sorted []Node, err error, expect string) {

	if sorted != nil {
		t.Error("sorted != nil")
	}

	if err == nil {
		t.Error("err == nil")
		return
	}

	if err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
	}
}

func TestGraph(t *testing.T) {

	nodes := []Node{
		Node{"a", nil},
		Node{"a", nil},
	}
	sorted, err := TopologicalSort(nodes)
	fail(t, sorted, err, "Duplicate node name: 'a'")

	nodes = []Node{
		Node{Name: "b", Children: []string{"c"}},
		Node{Name: "a", Children: []string{"d"}},
		Node{Name: "c"},
	}
	sorted, err = TopologicalSort(nodes)
	fail(t, sorted, err, "Cannot find node 'd'")

	nodes = []Node{
		Node{Name: "b", Children: []string{"c"}},
		Node{Name: "a", Children: []string{"b"}},
		Node{Name: "c", Children: []string{"a"}},
	}
	sorted, err = TopologicalSort(nodes)
	fail(t, sorted, err, "Graph cycle detected on node 'b'")

	nodes = []Node{
		Node{Name: "b", Children: []string{"c"}},
		Node{Name: "a", Children: []string{"b"}},
		Node{Name: "c"},
	}
	sorted, err = TopologicalSort(nodes)
	ok(t, sorted, err, []string{"a", "b", "c"})

	nodes = []Node{
		Node{Name: "f"},
		Node{Name: "g", Children: []string{"h", "i", "f"}},
		Node{Name: "c"},
		Node{Name: "b", Children: []string{"c", "d", "e"}},
		Node{Name: "d"},
		Node{Name: "a", Children: []string{"b"}},
		Node{Name: "e", Children: []string{"f", "g", "d"}},
		Node{Name: "h"},
		Node{Name: "i"},
	}
	sorted, err = TopologicalSort(nodes)
	ok(t, sorted, err, []string{"a", "b", "e", "d", "c", "g", "i", "h", "f"})

	nodes = []Node{
		Node{Name: "a", Children: []string{"b"}},
		Node{Name: "b", Children: []string{"c", "d"}},
		Node{Name: "c"},
		Node{Name: "d", Children: []string{"e", "f"}},
		Node{Name: "e"},
		Node{Name: "f", Children: []string{"a"}},
	}
	sorted, err = TopologicalSort(nodes)
	fail(t, sorted, err, "Graph cycle detected on node 'a'")
}
