package main

import (
	"log/slog"

	"github.com/m87/nod"
	sqlite_nod "github.com/m87/nod/sqlite"
)

func main() {
	repo, _ := sqlite_nod.NewRepository(":memory:", slog.Default(), nod.NewMapperRegistry())

	node := &nod.Node{
		Core: nod.NodeCore{Name: "example", Kind: "folder"},
	}
	repo.Save(node)

	q := repo.Query().NameEquals("example")
	found, _ := q.First()
	slog.Info("Found node", "name", found.Core.Name)
}
