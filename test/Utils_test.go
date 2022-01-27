package test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	utils "github.com/polyglotDataNerd/poly-Go-utils/aws"
	"github.com/polyglotDataNerd/poly-Go-utils/helpers"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"github.com/stretchr/testify/assert"
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
