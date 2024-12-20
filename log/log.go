package log

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

var defaultLogger *logrus.Entry = logrus.StandardLogger().WithField("go.version", runtime.Version())

// Logger provides a leveled-logging interface.
type Logger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})

	// Leveled methods, from logrus
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	WithError(err error) *logrus.Entry

	// TODO:
	// SetOutput sets output destination, it might be useful to suppress logging in tests.
	// SetOutput(io.Writer) error
}

// Fields is used as argument in WithFields method/func
type Fields = logrus.Fields

// Default is a default logger instance
func Default() Logger {
	return defaultLogger
}

func Print(args ...interface{})                 { defaultLogger.Print(args...) }
func Printf(format string, args ...interface{}) { defaultLogger.Printf(format, args...) }
func Println(args ...interface{})               { defaultLogger.Println(args...) }

func Fatal(args ...interface{})                 { defaultLogger.Fatal(args...) }
func Fatalf(format string, args ...interface{}) { defaultLogger.Fatalf(format, args...) }
func Fatalln(args ...interface{})               { defaultLogger.Fatalln(args...) }

func Panic(args ...interface{})                 { defaultLogger.Panic(args...) }
func Panicf(format string, args ...interface{}) { defaultLogger.Panicf(format, args...) }
func Panicln(args ...interface{})               { defaultLogger.Panicln(args...) }

func Debug(args ...interface{})                 { defaultLogger.Debug(args...) }
func Debugf(format string, args ...interface{}) { defaultLogger.Debugf(format, args...) }
func Debugln(args ...interface{})               { defaultLogger.Debugln(args...) }

func Error(args ...interface{})                 { defaultLogger.Error(args...) }
func Errorf(format string, args ...interface{}) { defaultLogger.Errorf(format, args...) }
func Errorln(args ...interface{})               { defaultLogger.Errorln(args...) }

func Info(args ...interface{})                 { defaultLogger.Info(args...) }
func Infof(format string, args ...interface{}) { defaultLogger.Infof(format, args...) }
func Infoln(args ...interface{})               { defaultLogger.Infoln(args...) }

func Warn(args ...interface{})                 { defaultLogger.Warn(args...) }
func Warnf(format string, args ...interface{}) { defaultLogger.Warnf(format, args...) }
func Warnln(args ...interface{})               { defaultLogger.Warnln(args...) }

// An entry is the final or intermediate logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
type Entry struct {
	*logrus.Entry
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func WithError(err error) *Entry {
	var entry Entry

	entry.Entry = defaultLogger.WithError(err)
	return &entry
}

// Add a map of fields to the Entry.
func WithFields(fields Fields) *Entry {
	var entry Entry

	entry.Entry = defaultLogger.WithFields(logrus.Fields(fields))
	return &entry
}

// Add a map of fields to the Entry.
func (e *Entry) WithFields(fields Fields) *Entry {
	e.Entry = e.Entry.WithFields(logrus.Fields(fields))
	return e
}
