package store

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/MixinNetwork/safe/governance/config"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemasql []byte

type Database struct {
	db *sql.DB
}

func OpenDatabase() (*Database, error) {
	path := config.AppConfig.Database.Path
	dsn := fmt.Sprintf("file:%s?mode=rwc&_journal_mode=WAL&cache=private", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(schemasql))
	if err != nil {
		return nil, err
	}
	return &Database{
		db: db,
	}, db.Ping()
}

func OpenSQLite3ReadOnlyStore(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?mode=ro", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

func (s *Database) Close() error {
	return s.db.Close()
}

func (s *Database) RunInTransaction(ctx context.Context, fn func(context.Context, *sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Database) Exec(ctx context.Context, query string, args ...any) error {
	_, err := s.db.ExecContext(ctx, query, args)
	return err
}

func (s *Database) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func BuildInsertionSQL(table string, cols []string) string {
	vals := strings.Repeat(", ?", len(cols)-1)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (?%s)", table, strings.Join(cols, ","), vals)
}

type Row interface {
	Scan(dest ...any) error
}
