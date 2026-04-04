# nod

Golang library for managing tree-structured data with support for tags, key-value attributes (KV), content, and transactions, built on GORM. Supports SQLite and PostgreSQL.

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
	mappers := nod.NewMapperRegistry()
	repo, err := sqlite_nod.NewRepository(":memory:", slog.Default(), mappers)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	node := &nod.Node{
		Core: nod.NodeCore{Name: "root", Kind: "folder"},
	}
	id, _ := repo.Save(node)

	found, _ := repo.Query().NodeId(id).First()
	slog.Info("Found node", "name", found.Core.Name)
}
```

For PostgreSQL, use the `postgres` adapter:

```go
import postgres_nod "github.com/m87/nod/postgres"

repo, err := postgres_nod.NewRepository(dsn, slog.Default(), mappers)
```

## Typed Repository

Register a mapper to work with your own domain models:

```go
type MyNode struct {
	Name string
}

type MyMapper struct{}

func (m MyMapper) ToNode(model *MyNode) (*nod.Node, error) {
	return &nod.Node{
		Core: nod.NodeCore{Name: model.Name, Kind: "my-node"},
	}, nil
}
func (m MyMapper) FromNode(node *nod.Node) (*MyNode, error) {
	return &MyNode{Name: node.Core.Name}, nil
}
func (m MyMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == "my-node"
}

func main() {
	registry := nod.NewMapperRegistry()
	nod.RegisterMapper(registry, MyMapper{})
	repo, _ := sqlite_nod.NewRepository(":memory:", slog.Default(), registry)

	typed := nod.As[MyNode](repo)
	typed.Save(&MyNode{Name: "example"})
	found, _ := typed.Query().NameEquals("example").First()
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

Run Postgres adapter tests (requires DSN):

```
export NOD_TEST_POSTGRES_DSN='host=localhost port=5432 user=nod password=nod dbname=nod_test sslmode=disable'
go test ./postgres -v
```

## License

MIT
