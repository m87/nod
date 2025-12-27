package sqlite

import (
	"github.com/m87/nod/core"
	"github.com/m87/nod/tags"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func NewRepository(path string) *core.Repository {
	db := initDB(path)

	return &core.Repository{
		Node: &core.NodeRepository{DB: db},
	}
}

func initDB(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&core.Node{}, &tags.Tag{}, &tags.NodeTag{})
	if err != nil {
		panic("failed to migrate database")
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	return db
}

