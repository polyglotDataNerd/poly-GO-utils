package utils

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jcxplorer/cwlogger"
	"io"
	"log"
	"os"
	"time"

	//"time"
)

type CloudWatch struct {
	LogGroup  string
	Retention int
}

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

	loggerClient := CloudWatch{
		LogGroup:  "yelp-parser",
		Retention: 24,
	}
	logs := loggerClient.CloudWatchPut()

	flag.Parse()
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.LstdFlags|log.Lshortfile)
	infoHandle.Write(buff.Bytes())
	logs.Log(time.Now(), buff.String())

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

func (c *CloudWatch) CloudWatchPut() *cwlogger.Logger {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	logger, err := cwlogger.New(&cwlogger.Config{
		Client:       cloudwatchlogs.New(sess),
		LogGroupName: c.LogGroup,
		Retention:    c.Retention,
	})
	if err != nil {
		fmt.Errorf("%v", err)
	}
	return logger
}