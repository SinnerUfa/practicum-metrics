// mlog - package for connecting loggers via interface slog
package mlog

import (
	"io"
	// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
	slog "golang.org/x/exp/slog" // slog "log/slog"
	"os"
	"path/filepath"
	"sync"
	// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
	// "go.uber.org/zap"
	// zapslog "go.uber.org/zap/exp/zapslog"
	// zapcore "go.uber.org/zap/zapcore"
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

func New(l LoggerType, lv slog.Level) *slog.Logger {
	once.Do(func() {
		o := []io.Writer{os.Stdout}
		f, err := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err == nil {
			o = append(o, f)
		}
		var handler slog.Handler
		switch l {
		case SlogType:
			handler = slogHandler(io.MultiWriter(o...), lv)
			// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20
		// case ZapType:
		// 	handler = zapHandler(io.MultiWriter(o...), lv)
		default:
			handler = slogHandler(io.MultiWriter(o...), lv)
		}
		instance = slog.New(handler)
	})

	return instance
}

func slogHandler(w io.Writer, lv slog.Level) *slog.JSONHandler {
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
		Level:       lv,
	}
	return slog.NewJSONHandler(w, options)
}

// временно закомментировано, в связи с необходимостью перехода на более версию go 1.20

// func zapHandler(w io.Writer, lv slog.Level) *zapslog.Handler {
// 	encoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()) // need fine encoder tuning

// 	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
// 		return lvl >= zapConvertSlogLevel(lv)
// 	})

// 	core := zapcore.NewCore(encoder, zapcore.AddSync(w), highPriority)

// 	return zapslog.NewHandler(core, &zapslog.HandlerOptions{AddSource: true})
// }

// func zapConvertSlogLevel(l slog.Level) zapcore.Level {
// 	switch {
// 	case l >= slog.LevelError:
// 		return zapcore.ErrorLevel
// 	case l >= slog.LevelWarn:
// 		return zapcore.WarnLevel
// 	case l >= slog.LevelInfo:
// 		return zapcore.InfoLevel
// 	default:
// 		return zapcore.DebugLevel
// 	}
// }
