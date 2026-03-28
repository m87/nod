package nod_test

import (
	"testing"
	"github.com/m87/nod"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log/slog"
)

func TestRepository_SaveAndQuery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	log := slog.Default()
	mappers := nod.NewMapperRegistry()
	repo, _ := nod.NewRepository(":memory:", log, mappers)

	node := &nod.Node{
		Core: nod.NodeCore{Name: "test", Kind: "testkind"},
	}
	_, err = repo.Save(node)
	if err != nil {
		t.Fatalf("failed to save node: %v", err)
	}

	q := repo.Query().NameEquals("test")
	found, err := q.First()
	if err != nil {
		t.Fatalf("failed to query node: %v", err)
	}
	if found.Core.Name != "test" {
		t.Errorf("expected name 'test', got %s", found.Core.Name)
	}
}

































}	}		t.Errorf("expected name 'test', got %s", found.Core.Name)	if found.Core.Name != "test" {	}		t.Fatalf("failed to query node: %v", err)	if err != nil {	found, err := q.First()	q := repo.Query().NameEquals("test")	}		t.Fatalf("failed to save node: %v", err)	if err != nil {	_, err = repo.Save(node)	}		Core: nod.NodeCore{Name: "test", Kind: "testkind"},	node := &nod.Node{	repo := nod.NewRepository(db, log, mappers)	mappers := nod.NewMapperRegistry()	log := slog.Default()	}		t.Fatalf("failed to open db: %v", err)	if err != nil {	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})func TestRepository_SaveAndQuery(t *testing.T) {)	"log/slog"	"gorm.io/gorm"	"gorm.io/driver/sqlite"	"github.com/m87/nod"	"testing"