package filelog

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"time"
)

type writer func([]byte) (int, error)

func (w writer) Write(b []byte) (int, error) {
	return w(b)
}

type filelog struct {
	LogDIR       string
	DefalutFlags int
	logs         map[Level]*log.Logger

	buffWriters map[string]*bufio.Writer
	fileLevels  map[string]Level
}

func (flg *filelog) Log(level Level, args ...any) {
	instance := flg.logs[level]
	if instance == nil {
		classifyWriter := func(b []byte) (n int, err error) {
			for k, v := range flg.fileLevels {
				if v&level != level {
					continue
				}

				bufWriter := flg.buffWriters[k]
				if bufWriter == nil {
					bufWriter = bufio.NewWriter(writer(
						func(bytes []byte) (n int, err error) {
							makedir(flg.LogDIR)
							name := fmt.Sprintf("%s/%s_%s.log", flg.LogDIR, time.Now().Format("2006-01-02"), k)
							f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
							if err != nil {
								return 0, err
							}

							n, err = f.Write(bytes)
							return
						},
					))
					flg.buffWriters[k] = bufWriter
				}

				n, err = bufWriter.Write(b)
			}

			return
		}

		flg.logs[level] = log.New(writer(classifyWriter), levelString(level)+" ", flg.DefalutFlags)
		instance = flg.logs[level]
	}

	instance.Println(args...)
}

func (flg *filelog) Close() (err error) {
	for _, v := range flg.buffWriters {
		if e := v.Flush(); e != nil {
			err = e
		}
	}
	return
}

func (flg *filelog) Debug(args ...any) {
	flg.Log(Ldebug, args...)
}

func (flg *filelog) Info(args ...any) {
	flg.Log(Linfo, args...)
}

func (flg *filelog) Warn(args ...any) {
	flg.Log(Lwarn, args...)
}

func (flg *filelog) Error(args ...any) {
	flg.Log(Lerror, args...)
}
func (flg *filelog) Debugf(format string, args ...any) {
	flg.Log(Ldebug, fmt.Sprintf(format, args...))
}

func (flg *filelog) Infof(format string, args ...any) {
	flg.Log(Linfo, fmt.Sprintf(format, args...))
}

func (flg *filelog) Warnf(format string, args ...any) {
	flg.Log(Lwarn, fmt.Sprintf(format, args...))
}

func (flg *filelog) Errorf(format string, args ...any) {
	flg.Log(Lerror, fmt.Sprintf(format, args...))
}

func (flg *filelog) Println(args ...any) {
	flg.Info(args...)
}

func (flg *filelog) Printf(format string, args ...any) {
	flg.Infof(format, args...)
}

// config

func (flg *filelog) SetFileLevels(logfile string, level Level) {
	flg.fileLevels[logfile] = level
}
func (flg *filelog) GetFileLevels() map[string]Level {
	return maps.Clone(flg.fileLevels)
}

func (flg *filelog) SetFlags(level Level, flags int) {
	flg.Log(level, "setflags", flags)
	flg.logs[level].SetFlags(flags)
}
