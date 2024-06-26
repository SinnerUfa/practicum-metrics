// mlog - package for connecting loggers via interface slog
// TODO: implement or find a more informative handler
package mlog

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
)

const (
	outFileName string = "out.log"
)

type LoggerType int

const (
	SlogType LoggerType = iota
	ZapType
)

func New(l LoggerType) *slog.Logger {
	once.Do(func() {
		o := []io.Writer{os.Stdout}
		f, err := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err == nil {
			o = append(o, f)
		}
		var handler slog.Handler
		switch l {
		case SlogType:
			handler = slogHandler(io.MultiWriter(o...))
		case ZapType:
			handler = zapHandler(io.MultiWriter(o...))
		default:
			handler = slogHandler(io.MultiWriter(o...))
		}

		instance = slog.New(handler)
	})

	return instance
}

func slogHandler(w io.Writer) *slog.JSONHandler {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.MessageKey:
			if a.Value.String() == "" {
				return slog.Attr{}
			}
		case slog.TimeKey:
			timeString := a.Value.Time().Format("15:05:05.00")
			return slog.Attr{
				Key:   slog.TimeKey,
				Value: slog.StringValue(timeString),
			}
		case slog.SourceKey:
			source := a.Value.Any().(*slog.Source)
			source = &slog.Source{
				File:     filepath.Base(filepath.Dir(source.File)) + "/" + filepath.Base(source.File),
				Function: filepath.Base(source.Function),
				Line:     source.Line,
			}

			return slog.Any(slog.SourceKey, source)
		}
		return a
	}
	options := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replace,
	}
	return slog.NewJSONHandler(w, options)
}

func zapHandler(w io.Writer) *slog.JSONHandler {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.MessageKey:
			if a.Value.String() == "" {
				return slog.Attr{}
			}
		case slog.TimeKey:
			timeString := a.Value.Time().Format("15:05:05.00")
			return slog.Attr{
				Key:   slog.TimeKey,
				Value: slog.StringValue(timeString),
			}
		case slog.SourceKey:
			source := a.Value.Any().(*slog.Source)
			source = &slog.Source{
				File:     filepath.Base(filepath.Dir(source.File)) + "/" + filepath.Base(source.File),
				Function: filepath.Base(source.Function),
				Line:     source.Line,
			}

			return slog.Any(slog.SourceKey, source)
		}
		return a
	}
	options := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replace,
	}
	return slog.NewJSONHandler(w, options)
}
