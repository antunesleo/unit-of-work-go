package core

import (
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
)

func NewDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "uowgo",
		Database: "uowgo",
	})
	db.AddQueryHook(pgdebug.DebugHook{Verbose: true})
	return db
}
