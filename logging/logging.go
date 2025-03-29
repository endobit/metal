// Package logging provides tools for using slog.Logger.
package logging

import (
	"fmt"
	"log/slog"
	"os"
)

// Options is the configuration for the logger.
type Options struct {
	Level string
	JSON  bool
}

var (
	DefaultLevel = "info"
	DefaultJSON  = false
)

//go:generate go tool "github.com/dmarkham/enumer" -type level -linecomment -text
type level int

const (
	levelDebug level = level(slog.LevelDebug) // debug
	levelInfo  level = level(slog.LevelInfo)  // info
	levelWarn  level = level(slog.LevelWarn)  // warn
	levelError level = level(slog.LevelError) // error
)

type flagSet interface {
	StringVar(*string, string, string, string)
	BoolVar(*bool, string, bool, string)
}

// NewOptions returns Options configured from the flagSet.
func NewOptions(f flagSet) *Options {
	var opt Options

	f.StringVar(&opt.Level, "log-level", DefaultLevel, "set the log level")
	f.BoolVar(&opt.JSON, "log-json", DefaultJSON, "output logs in JSON format")

	return &opt
}

// NewLogger returns a new logger configured from the Options.
func (o Options) NewLogger() (*slog.Logger, error) {
	l, err := levelString(o.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q", o.Level)
	}

	opts := slog.HandlerOptions{
		Level: slog.Level(l),
	}

	var handler slog.Handler

	if o.JSON {
		handler = slog.NewJSONHandler(os.Stdout, &opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, &opts)
	}

	logger := slog.New(handler)

	return logger, nil
}
