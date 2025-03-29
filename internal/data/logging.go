package data

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
)

type dbtx interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

var _ dbtx = &loggingDB{}

// loggingDB implements the dbtx interface using slog.Logger. The interfaces matches
// the db.DBTX interface used in SQLC generated code. All transactions are down
// with objects that implement this interface allow complete logging of all
// transactions. By default the logging level will be slog.LevelInfo.
type loggingDB struct {
	DB     *sql.DB
	Logger *slog.Logger
	Level  slog.Level
}

// ExecContext implements the dbtx interface.
func (d loggingDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	d.log("sql exec", query, args)
	return d.DB.ExecContext(ctx, query, args...)
}

// PrepareContext implements the dbtx interface.
func (d loggingDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	d.log("sql prepare", query, nil)
	return d.DB.PrepareContext(ctx, query)
}

// QueryContext implements the dbtx interface.
func (d loggingDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	d.log("sql query", query, args)
	return d.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext implements the dbtx interface.
func (d loggingDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	d.log("sql query", query, args)
	return d.DB.QueryRowContext(ctx, query, args...)
}

func (d *loggingDB) log(msg, query string, args []any) {
	tokens := strings.Split(query, ":")
	if len(tokens) > 2 {
		query = strings.TrimSpace(tokens[1])
	}

	d.Logger.Log(context.Background(), d.Level, msg, slog.String("query", query), slog.Any("args", args))
}
