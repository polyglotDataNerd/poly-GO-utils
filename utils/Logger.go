package utils

import (
	"bytes"
	"flag"
	"github.com/polyglotDataNerd/zib-Go-utils/aws"
	"io"
	"log"
	"os"
	"time"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	buff    bytes.Buffer
)

// https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html
func Init(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	loggerClient := aws.CloudWatch{
		LogGroup:  "yelp-parser",
		Retention: 24,
	}.CloudWatchPut()

	flag.Parse()
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	infoHandle.Write(buff.Bytes())
	loggerClient.Log(time.Now(), buff.String())

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	warningHandle.Write(buff.Bytes())

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	errorHandle.Write(buff.Bytes())
}
func init() {
	//var buff bytes.Buffer
	//file, err := os.OpenFile("/var/tmp/utils.log", os.O_TRUNC|os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	Error.Print(err)
	//}
	Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}
