// Package nod provides a tree-structured data management library for relational databases.
//
// It supports tags, key-value attributes (KV), rich text content, transactions,
// and typed domain model mapping via [NodeMapper]. The core storage is abstracted
// behind GORM, with adapters for SQLite (nod/sqlite) and PostgreSQL (nod/postgres).
//
// Basic usage:
//
//	repo, _ := sqlite.NewRepository(":memory:", slog.Default(), nod.NewMapperRegistry())
//	repo.Save(&nod.Node{Core: nod.NodeCore{Name: "root", Kind: "folder"}})
//	found, _ := repo.Query().NameEquals("root").First()
//
// For typed domain models, register a [NodeMapper] and use [As] to obtain a [TypedRepository]:
//
//	nod.RegisterMapper(registry, myMapper{})
//	typed := nod.As[MyModel](repo)
//	typed.Save(&MyModel{Name: "example"})
package nod
