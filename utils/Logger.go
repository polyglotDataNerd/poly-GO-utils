package utils

import (
	"flag"
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html
func Init(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	logpath := "/var/tmp/utils.log"
	file, err := os.OpenFile(logpath, os.O_TRUNC|os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// defer file.Close()
	if err != nil {
		Error.Print(err)
	}
	flag.Parse()
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	TFile := Trace
	TFile.SetOutput(file)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	IFile := Info
	IFile.SetOutput(file)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	Warning.SetOutput(file)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	Error.SetOutput(file)
}
func init() {
	Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}
