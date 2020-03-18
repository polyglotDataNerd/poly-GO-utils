package utils

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
)

type loggers struct {
	Trace     *log.Logger
	Info      *log.Logger
	Warning   *log.Logger
	Error     *log.Logger
	buff      bytes.Buffer
	outLogger *log.Logger
}

// https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html
func (l *loggers) Init(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	buff bytes.Buffer) {

	flag.Parse()
	l.Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	traceHandle.Write(buff.Bytes())

	l.Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	infoHandle.Write(buff.Bytes())

	l.Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	warningHandle.Write(buff.Bytes())

	l.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	errorHandle.Write(buff.Bytes())
}
func init() {
	l := loggers{}
	file, err := os.OpenFile("/var/tmp/utils.log", os.O_TRUNC|os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.Error.Print(err)
	}
	l.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, l.buff)
	l.outLogger.Printf("%s", l.buff.String())
	file.Write(l.buff.Bytes())
}
