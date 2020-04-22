package log

import (
	"fmt"
	"os"
)

type Logger interface {
	Info(m string)
	Infof(format string, v ...interface{})

	Debug(m string)
	Debugf(format string, v ...interface{})

	Warning(m string)
	Warningf(format string, v ...interface{})

	Error(m string)
	Errorf(format string, v ...interface{})

	Fatal(m string)
	Fatalf(format string, v ...interface{})
}

type StdLogger struct{}

func (l StdLogger) Info(m string) {
	fmt.Println(" INFO", m)
}

func (l StdLogger) Infof(format string, v ...interface{})    { l.Info(fmt.Sprintf(format, v...)) }
func (l StdLogger) Debug(m string)                           { fmt.Println("DEBUG", m) }
func (l StdLogger) Debugf(format string, v ...interface{})   { l.Debug(fmt.Sprintf(format, v...)) }
func (l StdLogger) Warning(m string)                         { fmt.Println("WARNING ", m) }
func (l StdLogger) Warningf(format string, v ...interface{}) { l.Warning(fmt.Sprintf(format, v...)) }
func (l StdLogger) Error(m string)                           { fmt.Println("ERROR ", m) }
func (l StdLogger) Errorf(format string, v ...interface{})   { l.Error(fmt.Sprintf(format, v...)) }
func (l StdLogger) Fatal(m string) {
	fmt.Println("FATAL ", m)
	os.Exit(1)
}
func (l StdLogger) Fatalf(format string, v ...interface{}) { l.Fatal(fmt.Sprintf(format, v...)) }

type QuietLogger struct{}

func (l QuietLogger) Info(string)                     {}
func (l QuietLogger) Infof(string, ...interface{})    {}
func (l QuietLogger) Debug(string)                    {}
func (l QuietLogger) Debugf(string, ...interface{})   {}
func (l QuietLogger) Error(string)                    {}
func (l QuietLogger) Errorf(string, ...interface{})   {}
func (l QuietLogger) Warning(string)                  {}
func (l QuietLogger) Warningf(string, ...interface{}) {}
func (l QuietLogger) Fatal(string)                    { os.Exit(1) }
func (l QuietLogger) Fatalf(string, ...interface{})   { os.Exit(1) }

type TestLogger string

func (l *TestLogger) m(m string)                          { *l = TestLogger(m) }
func (l *TestLogger) mf(f string, v ...interface{})       { l.m(fmt.Sprintf(f, v...)) }
func (l *TestLogger) Info(m string)                       { l.m(m) }
func (l *TestLogger) Infof(f string, v ...interface{})    { l.mf(f, v...) }
func (l *TestLogger) Debug(m string)                      { l.m(m) }
func (l *TestLogger) Debugf(f string, v ...interface{})   { l.mf(f, v...) }
func (l *TestLogger) Warning(m string)                    { l.m(m) }
func (l *TestLogger) Warningf(f string, v ...interface{}) { l.mf(f, v...) }
func (l *TestLogger) Error(m string)                      { l.m(m) }
func (l *TestLogger) Errorf(f string, v ...interface{})   { l.mf(f, v...) }
func (l *TestLogger) Fatal(m string)                      { l.m(m) }
func (l *TestLogger) Fatalf(f string, v ...interface{})   { l.mf(f, v...) }
