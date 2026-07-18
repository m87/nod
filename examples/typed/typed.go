package main

import (
	"log/slog"

	"github.com/m87/nod"
	sqlite_nod "github.com/m87/nod/sqlite"
)

type MyNode struct {
	name string
}

type MyAdapter struct{}

func (m MyAdapter) ToNode(node *MyNode) (*nod.Node, error) {
	return &nod.Node{
		Core: nod.NodeCore{Name: node.name, Kind: "my-node"},
	}, nil
}

func (m MyAdapter) FromNode(node *nod.Node) (*MyNode, error) {
	return &MyNode{name: node.Core.Name}, nil
}

func (m MyAdapter) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == "my-node"
}

func main() {
	registry := nod.NewAdapterRegistry()
	if err := nod.RegisterNodeAdapter(registry, MyAdapter{}); err != nil {
		panic(err)
	}

	repo, err := sqlite_nod.NewRepositoryInMemory(slog.Default(), registry)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	nodes := nod.Nodes[MyNode](repo)

	node := &MyNode{name: "example"}
	id, err := nodes.SaveNode(node)
	if err != nil {
		panic(err)
	}

	found, err := nod.NewTypedNodeQuery[MyNode](repo).
		Where(nod.NodeFields.Id.Equals(id)).
		FindAll()
	if err != nil {
		panic(err)
	}
	if len(found) == 0 {
		panic("node not found")
	}

	slog.Info("Found node", "name", found[0].name)
}
