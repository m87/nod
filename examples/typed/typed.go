package main

import (
	"log/slog"

	"github.com/m87/nod"
	sqlite_nod "github.com/m87/nod/sqlite"
)

type MyNode struct {
	name string
}

type MyMapper struct{}

func (m MyMapper) ToNode(node *MyNode) (*nod.Node, error) {
	return &nod.Node{
		Core: nod.NodeCore{Name: node.name, Kind: "my-node"},
	}, nil
}

func (m MyMapper) FromNode(node *nod.Node) (*MyNode, error) {
	return &MyNode{name: node.Core.Name}, nil
}

func (m MyMapper) IsApplicable(node *nod.Node) bool {
	return true
}

func main() {
	registry := nod.NewMapperRegistry()
	nod.RegisterMapper(registry, MyMapper{})
	repo, _ := sqlite_nod.NewRepository(":memory:", slog.Default(), registry)
	typedRepo := nod.As[MyNode](repo)

	node := &MyNode{name: "example"}
	typedRepo.Save(node)

	q := typedRepo.Query().NameEquals("example")
	found, _ := q.First()
	slog.Info("Found node", "name", found.name)
}
