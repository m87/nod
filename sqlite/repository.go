package sqlite

import (
	"github.com/m87/nod"
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

	err = db.AutoMigrate(&nod.Node{}, &nod.Tag{}, &nod.NodeTag{})
	if err != nil {
		panic("failed to migrate database")
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	return db
}

