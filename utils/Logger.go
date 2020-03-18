package utils

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/jcxplorer/cwlogger"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	//"time"
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
	errorHandle io.Writer,
	file *os.File) {

	flag.Parse()
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	Trace.SetOutput(file)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	Info.SetOutput(file)

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
	file, err := os.OpenFile("/var/tmp/utils.log", os.O_TRUNC|os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Error.Print(err)
	}
	Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, file)
	output, _ := ioutil.ReadFile("/var/tmp/utils.log")
	CloudWatchPut("yelp-parser", 30).Log(time.Now(), string(output))
}

func CloudWatchPut(logGroup string, retention int) *cwlogger.Logger {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	logger, err := cwlogger.New(&cwlogger.Config{
		Client:       cloudwatchlogs.New(sess),
		LogGroupName: logGroup,
		Retention:    retention,
	})
	if err != nil {
		fmt.Errorf("%v", err)
	}
	return logger
}
