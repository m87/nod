# nod

Golang library for managing tree-structured data with support for tags, key-value attributes (KV), content, and transactions, built on GORM. Includes an SQLite adapter.

## Installation

```
go get github.com/m87/nod
```

## Quick Start

```go
package main

import (
	"log/slog"

	"github.com/m87/nod"
	sqlite_nod "github.com/m87/nod/sqlite"
)

func main() {
	repo, err := sqlite_nod.NewRepositoryInMemory(slog.Default(), nil)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	node := &nod.Node{
		Core: nod.NodeCore{Name: "root", Kind: "folder"},
	}
	if _, err := repo.Nodes().SaveNode(node); err != nil {
		panic(err)
	}

	found, err := nod.NewNodeQuery(repo).
		Where(nod.NodeFields.Name.Equals("root")).
		FindAll()
	if err != nil {
		panic(err)
	}

	slog.Info("Found node", "name", found[0].Core.Name)
}
```

## Typed Nodes

Register an adapter to work with your own domain models:

```go
type MyNode struct {
	Name string
}

type MyAdapter struct{}

func (m MyAdapter) ToNode(model *MyNode) (*nod.Node, error) {
	return &nod.Node{
		Core: nod.NodeCore{Name: model.Name, Kind: "my-node"},
	}, nil
}
func (m MyAdapter) FromNode(node *nod.Node) (*MyNode, error) {
	return &MyNode{Name: node.Core.Name}, nil
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
	id, err := nodes.SaveNode(&MyNode{Name: "example"})
	if err != nil {
		panic(err)
	}
	found, err := nodes.GetNode(id)
	if err != nil {
		panic(err)
	}
	slog.Info("Found node", "name", found.Name)
}
```

## Documentation

- [GoDoc](https://pkg.go.dev/github.com/m87/nod)

## Tests

Run tests:

```
go test ./...
```

Run SQLite adapter tests:

```
go test ./sqlite -v
```

## License

MIT
