package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func setupDatabase(t *testing.T) (*Store, func(*testing.T)) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite memory db: %v", err)
	}
	teardown := func(t *testing.T) {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close db: %v", err)
		}
	}

	_, file, _, _ := runtime.Caller(0) // Figure out the path to the migrations/ directory
	migrations := filepath.Join(filepath.Dir(file), "migrations")

	if err := goose.SetDialect("sqlite"); err != nil {
		t.Fatalf("goose dialect error: %v", err)
	}

	if err := goose.Up(db, migrations); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	return NewStore(slog.New(slog.DiscardHandler), db), teardown
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(msg, append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\texp: %#v\n\tgot: %#v\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
