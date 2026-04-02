# nod

Golang library for managing tree-structured data with support for tags, key-value attributes (KV), content, and transactions, built on GORM/SQLite.

## Installation

```
go get github.com/m87/nod
```

## Quick Start

```go
import (
	"github.com/m87/nod"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log/slog"
)

func main() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	log := slog.Default()
	mappers := nod.NewMapperRegistry()
	repo, _ := nod.NewRepository(":memory:", log, mappers)
	// ...
}
```

## Usage Example

```go
// Create a new node
node := &nod.Node{
	Core: nod.NodeCore{Name: "root", Kind: "folder"},
}
repo.Save(node)

// Query
q := repo.Query().NameEquals("root")
found, _ := q.First()
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
