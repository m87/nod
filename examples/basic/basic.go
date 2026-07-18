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
		Core: nod.NodeCore{Name: "example", Kind: "folder"},
	}
	if _, err := repo.Nodes().SaveNode(node); err != nil {
		panic(err)
	}

	found, err := nod.NewNodeQuery(repo).
		Where(nod.NodeFields.Name.Equals("example")).
		FindAll()
	if err != nil {
		panic(err)
	}
	if len(found) == 0 {
		panic("node not found")
	}

	slog.Info("Found node", "name", found[0].Core.Name)
}
