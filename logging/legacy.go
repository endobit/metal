package logging

import (
	"fmt"
	"log/slog"
	"os"
)

type legacy interface {
	Fatalf(string, ...any)
	Printf(string, ...any)
}

var _ legacy = &Legacy{}

// Legacy implements the legacy interface using slog.Logger. This is a common
// pattern in packages, for example goose.Logger matches this interface.
type Legacy struct {
	Logger *slog.Logger
}

// Fatalf implements the legacy logger interface.
func (l Legacy) Fatalf(format string, v ...any) {
	l.Logger.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Printf implements the legacy logger interface.
func (l Legacy) Printf(format string, v ...any) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}
