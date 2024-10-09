package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"time"

	"milton/pkg/db/orm"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDB struct {
	db  *sql.DB
	Orm *orm.Queries
}

type SqliteOpts struct {
	ConnString string
	DbName     string
}

//go:embed sqlc/schema.sql
var ddl string

func New(opts *SqliteOpts) *SqliteDB {
	ctx := context.Background()

	if opts.ConnString == "" {
		opts.ConnString = os.Getenv("DB_URL")
	}
	if opts.DbName == "" {
		opts.DbName = "database"
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", opts.ConnString))
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(2)

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

	orm := orm.New(db)
	return &SqliteDB{db: db, Orm: orm}
}

func str(v string) sql.NullString {
	return sql.NullString{String: v, Valid: true}
}

func t(s string) sql.NullTime {
	v, _ := time.Parse("2006-01-02", s)
	return sql.NullTime{Time: v, Valid: true}

}
