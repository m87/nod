package main

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

	node := &nod.Node{
		Core: nod.NodeCore{Name: "example", Kind: "folder"},
	}
	repo.Save(node)

	q := repo.Query().NameEquals("example")
	found, _ := q.First()
	log.Info("Found node", "name", found.Core.Name)
	}























}	log.Info("Found node", "name", found.Core.Name)	found, _ := q.First()	q := repo.Query().NameEquals("example")	repo.Save(node)	}		Core: nod.NodeCore{Name: "example", Kind: "folder"},	node := &nod.Node{	repo, _ := nod.NewRepository(":memory:", log, mappers)	mappers := nod.NewMapperRegistry()	log := slog.Default()	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})func main() {)	"log/slog"	"gorm.io/gorm"	"gorm.io/driver/sqlite"	"github.com/m87/nod"import (