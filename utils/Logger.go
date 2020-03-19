package utils

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Strings strings.Builder
)

/*
https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html

Softlinks (symlinks) stdout to /var/tmp/utils.log
ln -s /dev/stdout /var/tmp/utils.log
unlink /var/tmp/utils.log

*/
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	file *os.File,
) {
	flag.Parse()

	Trace = log.New(io.MultiWriter(traceHandle, file),
		"TRACE: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	Info = log.New(io.MultiWriter(infoHandle, file),
		"INFO: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	Warning = log.New(io.MultiWriter(warningHandle, file),
		"WARNING: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	Error = log.New(io.MultiWriter(errorHandle, file),
		"ERROR: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
}

func init() {

	file, err := os.OpenFile("/var/tmp/utils.log", os.O_TRUNC|os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Error.Print(err)
	}
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, file)
	Info.Println("Logger Initialized")
}
