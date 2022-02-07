package test

//go test -c -o tests

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	utils "github.com/polyglotDataNerd/poly-Go-utils/aws"
	"github.com/polyglotDataNerd/poly-Go-utils/helpers"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSessionGenerator(t *testing.T) {
	/* gets fixture from testdata folder */
	parentDir, _ := helpers.GetTestDir()
	credPath := fmt.Sprintf("%s%s", parentDir, "/credentials")

	/* helper to mock AWS client */
	mockServer := MockServer()
	t.Logf("mock server URL: %s", mockServer.URL)
	conf := utils.Settings{AWSConfig: &aws.Config{
		Endpoint: aws.String(mockServer.URL),
		Region:   aws.String("us-east-1"),
		/* name of test profile */
		Credentials: credentials.NewSharedCredentials(credPath, "testing"),
	}}

	testcases := []struct {
		session    *session.Session
		descriptor string
		testNumber int
	}{
		{
			session:    conf.SessionGenerator(),
			descriptor: "nil argurment uses default credentials",
			testNumber: 1,
		},
		{
			session:    conf.SessionGenerator(""),
			descriptor: "empty argurment uses default credentials",
			testNumber: 2,
		},
		{
			session:    conf.SessionGenerator("testing"),
			descriptor: "uses a profile",
			testNumber: 3,
		},
		{
			session:    conf.SessionGenerator("AKIA2L3E6U3OMSEI7H5L", "TEST"),
			descriptor: "uses an acess key and secret key",
			testNumber: 4,
		},
	}

	for _, tc := range testcases {
		testAccessKey := "AKIA2L3E6U3OMSEI7H5L"
		creds, _ := tc.session.Config.Credentials.Get()

		if tc.testNumber == 4 {
			secretKey := "TEST"
			msg := fmt.Sprintf("test number %d: %s and testing secret key: %s", tc.testNumber, tc.descriptor, creds.SecretAccessKey)
			log.Info.Printf(msg)
			assert.Equal(t, secretKey, creds.SecretAccessKey, msg)
		} else if tc.testNumber == 3 {
			msg := fmt.Sprintf("test number %d: %s and testing access/secret keys: %s:%s", tc.testNumber, tc.descriptor, creds.AccessKeyID, creds.SecretAccessKey)
			log.Info.Printf(msg)
			assert.Equal(t, testAccessKey, creds.AccessKeyID, msg)
		} else {
			msg := fmt.Sprintf("test number %d: %s and testing access key: %s", tc.testNumber, tc.descriptor, creds.AccessKeyID)
			log.Info.Println(msg)
			assert.Equal(t, testAccessKey, creds.AccessKeyID, msg)
		}

	}
}

func TestS3ReadObj(t *testing.T) {
	/* gets fixture from testdata folder */
	parentDir, _ := helpers.GetTestDir()
	fixturePath := fmt.Sprintf("%s%s", parentDir, "/s3/")
	cli := S3MockDocker()
	objectTextTest := "This is the test body"
	testGzipText := "\"name\"\t\"level\"\t\"city\"\t\"county\"\t\"state\"\t\"country\"\t\"population\"\t\"Latitude\"\t\"Longitude\"\t\"aggregate\"\t\"timezone\"\t\"cases\"\t\"US_Confirmed_County\"\t\"deaths\"\t\"US_Deaths_County\"\t\"recovered\"\t\"US_Recovered_County\"\t\"active\"\t\"US_Active_County\"\t\"tested\"\t\"hospitalized\"\t\"discharged\"\t\"last_updated\"\t\"icu\"\t\"hospitalized_current\"\t\"icu_current\"\n\"\"\t\"\"\t\"\"\t\"\"\t\"\"\t\"Zimbabwe\"\t\"\"\t\"\"\t\"\"\t\"\"\t\"\"\t\"\"\t\"32952.0\"\t\"\"\t\"1178.0\"\t\"\"\t\"24872.0\"\t\"\"\t\"6902.0\"\t\"\"\t\"\"\t\"\"\t\"2021-01-30\"\t\"\"\t\"\"\t\"\"\n\"\"\t\"\"\t\"\"\t\"\"\t\"\"\t\"Zimbabwe\"\t\"\"\t\"\"\t\"\"\t\"\"\t\"\"\t\"\"\t\"605.0\"\t\"\"\t\"7.0\"\t\"\"\t\"166.0\"\t\"\"\t\"432.0\"\t\"\"\t\"\"\t\"\"\t\"2020-07-02\"\t\"\"\t\"\"\t\"\""
	out, errC := cli.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("poly-test")})
	log.Info.Println(out.GoString())
	if errC != nil {
		log.Error.Println(errC)
	}

	testcases := []struct {
		utils      utils.S3Obj
		descriptor string
		testNumber int
	}{
		{
			utils: utils.S3Obj{
				Bucket: "poly-test",
				Key:    "testing/test.csv",
			},
			descriptor: "test read object",
			testNumber: 1,
		},
		{
			utils: utils.S3Obj{
				Bucket: "poly-test",
				Key:    "testing/test.gzip",
			},
			descriptor: "test read object gzip from testdata dir",
			testNumber: 2,
		},
	}

	for _, tc := range testcases {
		if tc.testNumber == 1 {
			input := s3.PutObjectInput{
				Body:                 bytes.NewReader([]byte("This is the test body")),
				Bucket:               aws.String("poly-test"),
				Key:                  aws.String("testing/test.csv"),
				ServerSideEncryption: aws.String("AES256"),
				StorageClass:         aws.String("STANDARD"),
			}
			result, err := cli.PutObject(&input)
			if err != nil {
				log.Error.Println(err)
			}
			log.Info.Println(result)
			s3Session := utils.Settings{AWSConfig: &cli.Config}
			testText, _ := tc.utils.S3ReadObj(s3Session.SessionGenerator())
			msg := fmt.Sprintf("test number %d: S3ReadObj method validates ObjectContent behavior output passed, textbody: %s", tc.testNumber, testText)
			log.Info.Println(msg)
			assert.Equal(t, objectTextTest, testText, msg)
		} else if tc.testNumber == 2 {
			upFile, errF := os.OpenFile(fmt.Sprintf("%s%s", fixturePath, "test.gzip"), os.O_RDWR, 0644)
			if errF != nil {
				log.Error.Panic(errF)
			}
			defer upFile.Close()
			fileInfo, _ := upFile.Stat()
			buffer := make([]byte, fileInfo.Size())
			upFile.Read(buffer)
			input := s3.PutObjectInput{
				Body:                 bytes.NewReader(buffer),
				Bucket:               aws.String("poly-test"),
				Key:                  aws.String("testing/test.gzip"),
				ServerSideEncryption: aws.String("AES256"),
				StorageClass:         aws.String("STANDARD"),
			}
			result, err := cli.PutObject(&input)
			if err != nil {
				log.Error.Println(err)
			}
			log.Info.Println(result)
			s3Session := utils.Settings{AWSConfig: &cli.Config}
			testText, _ := tc.utils.S3ReadObjGzip(s3Session.SessionGenerator())
			msg := fmt.Sprintf("test number %d: S3ReadObjGzip method validates ObjectContent behavior output passed, gzip body: %s", tc.testNumber, testText)
			log.Info.Println(msg)
			assert.Equal(t, testGzipText, testText, msg)
		}
	}
}
