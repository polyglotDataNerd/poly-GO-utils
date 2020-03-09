package sgutils

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
	flag.Parse()
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
}
func init() {
	Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}
