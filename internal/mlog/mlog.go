package mlog

import (
	"io"
	"log"
	"os"
	"sync"
)

const (
	LogFatal   = "FATAL: "
	LogPanic   = "PANIC: "
	LogWarning = "WARNING: "
	LogInfo    = "INFORMATION: "
	LogDebug   = "DEBUG: "
)

type Logger interface {
	Fatal(v ...any)
	Panic(v ...any)
	Warning(v ...any)
	Info(v ...any)
	Debug(v ...any)
}

type mlog struct {
	showDebug bool
	errLogger *log.Logger
	outLogger *log.Logger
}

var (
	instance *mlog
	once     sync.Once
)

const (
	outFileName string = "out.log"
	errFileName string = "err.log"
)

func New(showDebug bool) Logger {
	once.Do(func() {
		o := []io.Writer{os.Stdout}
		f, err := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err == nil {
			o = append(o, f)
		}

		e := []io.Writer{os.Stderr}
		f, err = os.OpenFile(errFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err == nil {
			e = append(e, f)
		}
		instance = &mlog{
			showDebug: showDebug,
			outLogger: log.New(io.MultiWriter(o...), "", log.LstdFlags|log.LUTC),
			errLogger: log.New(io.MultiWriter(e...), "", log.LstdFlags|log.LUTC),
		}
	})
	return instance
}

func (m *mlog) Fatal(v ...any) {
	m.errLogger.SetPrefix(LogFatal)
	m.errLogger.Fatal(v...)
}

func (m *mlog) Panic(v ...any) {
	m.errLogger.SetPrefix(LogPanic)
	m.errLogger.Panic(v...)
}
func (m *mlog) Warning(v ...any) {
	m.outLogger.SetPrefix(LogWarning)
	m.outLogger.Println(v...)
}

func (m *mlog) Info(v ...any) {
	m.outLogger.SetPrefix(LogInfo)
	m.outLogger.Println(v...)
}
func (m *mlog) Debug(v ...any) {
	if m.showDebug {
		m.outLogger.SetPrefix(LogDebug)
		m.outLogger.Println(v...)
	}
}
