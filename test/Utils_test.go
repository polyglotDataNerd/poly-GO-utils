package test

import (
	"github.com/aws/aws-sdk-go/aws"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeConfig struct {
	mockConfig *aws.Config
}

func SessionGenerator(t *testing.T) {
	mockServer := MockServer()
	//log.Info.Println(aws.StringValue(&creds.AccessKeyID))
	//conf := FakeConfig{&aws.Config{
	//	Region: aws.String("us-east-1"),
	//	DisableSSL: aws.Bool(true),
	//	Endpoint: aws.String(mockServer.URL),
	//}}
	log.Info.Println(mockServer.URL)

	testAccessKey := "ASIA2PWORN3JBUX4Y4XX"
	assert.Equal(t, testAccessKey, testAccessKey)

}
