// package function

package filelog

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"maps"
	"os"
	"sync"
)

var (
	_default *filelog
	one      = new(sync.Once)
)

const (
	Ldebug Level = 1 << iota
	Linfo
	Lwarn
	Lerror
)

type Level int

func levelString(level Level) string {
	switch level {
	case Ldebug:
		return "debug"
	case Linfo:
		return "info"
	case Lwarn:
		return "warn"
	case Lerror:
		return "error"
	}
	return ""
}
func makedir(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		if _, ok := err.(*fs.PathError); ok {
			os.MkdirAll(dir, 0766)
		} else {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
	}
}

func newLogger() *filelog {
	return &filelog{
		LogDIR:       "filelogs",
		DefalutFlags: log.LstdFlags | log.Lshortfile,
		logs:         make(map[Level]*log.Logger),
		buffWriters:  make(map[string]*bufio.Writer),
		fileLevels: map[string]Level{
			"debug": Ldebug,
			"info":  Linfo,
			"warn":  Lwarn,
			"error": Lerror,
		},
	}
}

func New() *filelog {
	return newLogger()
}

func GetDefault() *filelog {
	one.Do(func() {
		_default = newLogger()
	})

	return _default
}

func Debug(args ...any) {
	GetDefault().Debug(args...)
}

func Info(args ...any) {
	GetDefault().Info(args...)
}

func Warn(args ...any) {
	GetDefault().Warn(args...)
}

func Error(args ...any) {
	GetDefault().Error(args...)
}
func Debugf(format string, args ...any) {
	GetDefault().Debugf(format, args...)
}

func Infof(format string, args ...any) {
	GetDefault().Infof(format, args...)
}

func Warnf(format string, args ...any) {
	GetDefault().Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	GetDefault().Errorf(format, args...)
}

func Print(args ...any) {
	Info(args...)
}

func Printf(format string, args ...any) {
	Infof(format, args...)
}

func Close() (err error) {
	return GetDefault().Close()
}

// config

func SetFileLevels(logfile string, level Level) {
	GetDefault().fileLevels[logfile] = level
}
func GetFileLevels() map[string]Level {
	return maps.Clone(GetDefault().fileLevels)
}

func SetFlags(level Level, flags int) {
	GetDefault().Log(level, "setflags", flags)
	GetDefault().logs[level].SetFlags(flags)
}
