package aws

/* for vgo install sdk using
go get -u github.com/goaws/goaws-sdk-go
this will install outside of the defaulted GOPATH
*/

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jcxplorer/cwlogger"
	goutils "github.com/polyglotDataNerd/zib-Go-utils/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	s3_region = "us-west-2"
	CharSet   = "UTF-8"
)

type S3Obj struct {
	Bucket string
	Key    string
}

type CloudWatch struct {
	LogGroup  string
	Retention int
}

/*dynamo db table struct*/
type DDBT struct {
	HashKey   string
	Attribute string
}

type DDB interface {
	DDBScanGetItems()
	DDBGetItem()
}

/*
Variadic function parameters are like *kwargs in python that have flexible params:
ags ...string
*/
func SessionGenerator(args ...string) *session.Session {
	/*default profile to create session*/
	if args == nil {
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(s3_region)},
		)
		//checks if sessions passes valid creds
		_, credError := sess.Config.Credentials.Get()
		if credError != nil {
			goutils.Error.Fatalln("session not valid", credError.Error())
		}
		return sess
	}
	/*
			pass explicit profile and region name
			 args[0] = profile name
		     args[1] = region
	*/
	if len(args[0]) < 10 {
		sess, _ := session.NewSessionWithOptions(
			session.Options{
				Config:  aws.Config{Region: aws.String(args[1])},
				Profile: args[0],
			})
		//checks if sessions passes valid creds
		_, credError := sess.Config.Credentials.Get()
		if credError != nil {
			goutils.Error.Fatalln("session not valid", credError.Error())
		}
		return sess
	}
	/*
			pass explicit access, secret key and region name
			 args[0] = access key
		     args[1] = secret key
			 args[2] = region
	*/
	if len(args[0]) > 10 {
		creds := credentials.NewStaticCredentials(args[0], args[1], "")
		sess, _ := session.NewSession(&aws.Config{
			Region:      aws.String(args[2]),
			Credentials: creds,
		},
		)
		//checks if sessions passes valid creds
		_, credError := sess.Config.Credentials.Get()
		if credError != nil {
			goutils.Error.Fatalln("session not valid", credError.Error())
		}
		return sess
	}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(s3_region)},
	)

	return sess
}

func (obj S3Obj) S3ReadObj(sess *session.Session) (string, error) {
	s3cli := s3.New(sess)

	getObject, err := s3cli.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(obj.Bucket),
			Key:    aws.String(obj.Key),})

	if err != nil {
		goutils.Error.Fatalln("no s3 Object")
		return "no s3 Object", err
	}

	defer getObject.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, getObject.Body); err != nil {
		return "object malformed", err
	}
	return string(buf.Bytes()), nil
}

func (obj S3Obj) S3ReadObjGzip(sess *session.Session) (string, error) {
	//sess := SessionGenerator()
	s3cli := s3.New(sess)

	getObject, err := s3cli.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(obj.Bucket),
			Key:    aws.String(obj.Key),})

	if err != nil {
		goutils.Error.Fatalln(err.Error())
		return "no s3 Object", err
	}
	/*closes bytestream*/
	defer getObject.Body.Close()

	/*put original object in byte buffer*/
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, getObject.Body); err != nil {
		return "object malformed", err
	}

	/*put compressed byte object in gzip Reader*/
	gzipO, err := gzip.NewReader(bytes.NewBuffer(buf.Bytes()))
	/*closes bytestream*/
	defer gzipO.Close()

	data, err := ioutil.ReadAll(gzipO)
	if err != nil {
		return "object not GZIP format", err
	}
	return string(data), nil

}

func (obj S3Obj) S3ReadObjGZIPDir(sess *session.Session) (map[string]string, error) {
	objectMap := make(map[string]string)
	s3cli := s3.New(sess)

	s3resp, _ := s3cli.ListObjects(
		&s3.ListObjectsInput{
			Bucket: aws.String(obj.Bucket),
			Prefix: aws.String(obj.Key),
		})

	for _, item := range s3resp.Contents {
		getObject, err := s3cli.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(obj.Bucket),
				Key:    aws.String(*item.Key),})

		if err != nil {
			goutils.Error.Fatalln("no s3 Object", err)
		}
		/*closes bytestream*/
		defer getObject.Body.Close()

		/*put original object in byte buffer*/
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, getObject.Body); err != nil {
		}

		/*put compressed byte object in gzip Reader*/
		gzipO, err := gzip.NewReader(bytes.NewBuffer(buf.Bytes()))
		/*closes bytestream*/
		defer gzipO.Close()

		data, err := ioutil.ReadAll(gzipO)
		if err != nil {
			goutils.Error.Fatalln("object not GZIP format", err)
		}
		objectMap[*item.Key] = string(data)
	}
	return objectMap, nil
}

func (obj S3Obj) S3ReadObjDir(sess *session.Session) (map[string]string, error) {
	objectMap := make(map[string]string)
	s3cli := s3.New(sess)

	s3resp, _ := s3cli.ListObjects(
		&s3.ListObjectsInput{
			Bucket: aws.String(obj.Bucket),
			Prefix: aws.String(obj.Key),
		})

	for _, item := range s3resp.Contents {
		getObject, err := s3cli.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(obj.Bucket),
				Key:    aws.String(*item.Key),})

		if err != nil {
			goutils.Error.Fatalln("no s3 Object", err)
		}
		/*closes bytestream*/
		defer getObject.Body.Close()

		/*put original object in byte buffer*/
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, getObject.Body); err != nil {
		}
		objectMap[filepath.Base(*item.Key)] = string(buf.Bytes())
	}
	return objectMap, nil
}

/* string implementation */
func (obj S3Obj) S3WriteGzip(builder string, sess *session.Session) {
	s3cli := s3.New(sess)

	/*put original object in byte buffer*/
	var b bytes.Buffer
	/*create new gzip writer*/
	gz := gzip.NewWriter(&b)
	/*converts string to bytes to write into gzip*/
	if _, byteerr := gz.Write([]byte(builder)); byteerr != nil {
		goutils.Info.Panic("object malformed", byteerr.Error())
	}
	if byteerr := gz.Close(); byteerr != nil {
		goutils.Info.Panic("object malformed", byteerr.Error())
	}

	input := &s3.PutObjectInput{
		Body:                 bytes.NewReader(b.Bytes()),
		Bucket:               aws.String(obj.Bucket),
		Key:                  aws.String(obj.Key),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("STANDARD_IA"),
	}
	result, err := s3cli.PutObject(input)
	if err != nil {
		goutils.Error.Fatalln(err)
	}
	goutils.Info.Println(result)

}

/* string reader implementation */
func (obj S3Obj) S3WriteGzipReader(reader io.Reader, sess *session.Session) {
	s3cli := s3.New(sess)
	/*reads payload*/
	payload, byteerr := ioutil.ReadAll(reader)
	if byteerr != nil {
		goutils.Info.Panic("object malformed", byteerr.Error())
	}
	/*put original object in byte buffer*/
	var b bytes.Buffer
	/*create new gzip writer*/
	gz := gzip.NewWriter(&b)
	defer gz.Close()

	/*converts string to bytes to write into gzip*/
	if _, byteerr := gz.Write(payload); byteerr != nil {
		goutils.Info.Panic("object malformed", byteerr.Error())
	}
	if byteerr := gz.Close(); byteerr != nil {
		goutils.Info.Panic("object malformed", byteerr.Error())
	}
	input := &s3.PutObjectInput{
		Body:                 bytes.NewReader(b.Bytes()),
		Bucket:               aws.String(obj.Bucket),
		Key:                  aws.String(obj.Key),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("STANDARD"),
	}
	result, err := s3cli.PutObject(input)
	if err != nil {
		goutils.Error.Fatalln(err)
	}
	goutils.Info.Println(result)

}

/* reader implementation */
func (obj S3Obj) S3UploadGzip(reader io.Reader, sess *session.Session) {
	/*reads payload*/
	payload, byteerr := ioutil.ReadAll(reader)
	if byteerr != nil {
		goutils.Info.Panic("object malformed", byteerr.Error())
	}
	/*put original object in byte buffer*/
	var b bytes.Buffer
	/*create new gzip writer*/
	gz := gzip.NewWriter(&b)

	if _, gzerr := gz.Write(payload); gzerr != nil {
		goutils.Info.Panic("object malformed", gzerr.Error())
	}
	if gzerr := gz.Close(); gzerr != nil {
		goutils.Info.Panic("object malformed", gzerr.Error())
	}
	/*creates client and then does a put*/
	s3cli := s3manager.NewUploader(sess)
	_, err := s3cli.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(obj.Bucket),
		Key:                  aws.String(obj.Key),
		ServerSideEncryption: aws.String("AES256"),
		Body:                 bytes.NewReader(b.Bytes()),
	})

	if err != nil {
		goutils.Error.Fatalln(err.Error())
	}
}

func SESEmail(emailTo string, emailFrom string, subject string, body string) {
	sess := SessionGenerator()
	cli := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(emailTo),
			},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data:    aws.String(subject),
				Charset: aws.String(CharSet),
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Data:    aws.String(body),
					Charset: aws.String(CharSet),
				},
			},
		},
		Source: aws.String(emailFrom),
	}

	goutils.Info.Println("Email Sent to address: " + emailTo)

	result, err := cli.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				goutils.Info.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				goutils.Info.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				goutils.Info.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				goutils.Info.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			goutils.Info.Println(err.Error())
		}
		return
	}
	goutils.Info.Println(result)

}

func SSMParams(params string, index int) (output string) {
	paramArray := strings.Split(params, ",")
	sess := SessionGenerator()
	cli := ssm.New(sess)

	/*SSM pattern*/
	paramInput := ssm.GetParametersInput{Names: aws.StringSlice(paramArray), WithDecryption: aws.Bool(true)}
	req, resp := cli.GetParametersRequest(&paramInput)

	paramsvalue := req.Send()
	if paramsvalue == nil {
		output = aws.StringValue(resp.Parameters[index].Value)
	}
	return
}

func (mapper DDBT) DDBGetQuery(tableName string, index string, indexedProperty string) *map[string]string {
	/*main MAP struct to pass to DynamoDB*/
	baseMap := make(map[string]string)
	sess := SessionGenerator()
	cli := dynamodb.New(sess)

	query := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String(index),
		KeyConditions: map[string]*dynamodb.Condition{
			indexedProperty: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(mapper.Attribute),
					},
				},
			},
		},
	}
	resp, err := cli.Query(query)
	if err != nil {
		fmt.Print(err.Error(), err)
		goutils.Info.Println(err.Error(), err)
		os.Exit(1)
	}
	for _, dmap := range resp.Items {
		baseMap[*dmap["sourceSystemId"].S] = *dmap["uuid"].S
	}
	return &baseMap
}

func (mapper DDBT) DDBScanGetItems(tableName string, key string, attribute string) map[string]string {
	baseMap := make(map[string]string)
	/*main abstract struct to pass to DynamoDB*/
	sess := SessionGenerator()
	cli := dynamodb.New(sess)

	var query = &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	repsonse, err := cli.Scan(query)
	if err != nil {
		goutils.Info.Println(err.Error(), err)
	}
	for _, i := range repsonse.Items {
		/*DDB primary key*/
		keyErr := dynamodbattribute.Unmarshal(i[key], &mapper.HashKey)
		if keyErr != nil {
			fmt.Print(keyErr.Error(), err)
			goutils.Info.Println(keyErr.Error(), err)
			os.Exit(1)
		}
		/*DDB Attribute*/
		valueErr := dynamodbattribute.Unmarshal(i[attribute], &mapper.Attribute)
		if valueErr != nil {
			fmt.Print(valueErr.Error(), err)
			goutils.Info.Println(valueErr.Error(), err)
			os.Exit(1)
		}
		baseMap[mapper.HashKey] = mapper.Attribute
	}
	return baseMap
}

func (c *CloudWatch) CloudWatchPut() *cwlogger.Logger {
	sess := SessionGenerator("default", "us-west-2")
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
