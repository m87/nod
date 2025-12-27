package nod

import (
	nod "github.com/m87/nod/core"
	"github.com/m87/nod/tags"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func NewRepository(path string) *nod.Repository {
	db := initDB(path)

	return &nod.Repository{
		Node: &nod.NodeRepository{DB: db},
	}
}

func initDB(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&nod.Node{}, &tags.Tag{}, &tags.NodeTag{})
	if err != nil {
		panic("failed to migrate database")
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	return db
}

