package sqlite

import (
	"github.com/m87/nod"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func NewRepository(path string) *nod.Repository {
	db := initDB(path)

	return &nod.Repository{
		Db:   db,
		Node: &nod.NodeRepository{DB: db},
	}
}

func initDB(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Exec("PRAGMA foreign_keys = ON;")

	err = db.AutoMigrate(&nod.NodeCore{}, &nod.Tag{}, &nod.NodeTag{})
	if err != nil {
		panic("failed to migrate database")
	}


	return db
}


