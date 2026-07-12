// Zerolog Config
//
// Logging levels guide:
// TRACE (-1): for tracing the code execution path.
// DEBUG (0): messages useful for troubleshooting the program.
// INFO (1): messages describing the normal operation of an application.
// WARNING (2): for logging events that need may need to be checked later.
// ERROR (3): error messages for a specific operation.
// FATAL (4): severe errors where the application cannot recover. os.Exit(1) is called after the message is logged.
// PANIC (5): similar to FATAL, but panic() is called instead.
package logger

import (
	"adeeb_huma/config"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var log zerolog.Logger

func Init() {
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return pkgerrors.MarshalStack(errors.WithStack(err))
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var output io.Writer
	var lvl zerolog.Level
	if config.APP_ENV == "dev" {
		output = zerolog.ConsoleWriter{
			Out:           os.Stdout,
			TimeFormat:    time.RFC3339,
			FieldsExclude: []string{"stack"},
			FormatExtra: func(evt map[string]interface{}, buf *bytes.Buffer) error {
				if stack, ok := evt["stack"]; ok {
					buf.WriteString("\nstacktrace:\n")
					for _, item := range stack.([]interface{}) {
						stackItem := item.(map[string]interface{})
						str := fmt.Sprintf("%v at %v:%v \n", stackItem["source"], stackItem["func"], stackItem["line"])
						buf.WriteString(str)
					}
				}
				return nil
			},
		}
		lvl = zerolog.DebugLevel
	} else {
		output = zerolog.MultiLevelWriter(os.Stderr)
		lvl = zerolog.TraceLevel
	}

	log = zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Logger()
	// Sample(zerolog.LevelSampler{
	// TraceSampler: &zerolog.BasicSampler{N: 10},  // Log 1 out of every 10 TRACE messages
	// DebugSampler: &zerolog.BasicSampler{N: 5},   // Log 1 out of every 5 DEBUG messages
	// InfoSampler:  &zerolog.BasicSampler{N: 2},   // Log 1 out of every 2 INFO messages
	// })

}

func Get() zerolog.Logger {
	return log
}

// functions to shorten basic usage of zerolog,
// Instead of using l := logger.Get()
// then using it like: l.Info().Msg("...")
// You can just call it directly: logger.Info().Msg("...")

func Trace() *zerolog.Event {
	return log.Trace()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Panic() *zerolog.Event {
	return log.Panic()
}
