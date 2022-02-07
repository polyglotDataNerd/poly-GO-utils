package test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/localstack"
	"github.com/polyglotDataNerd/poly-Go-utils/helpers"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"net/http"
	"net/http/httptest"
)

func MockServer() *httptest.Server {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	return mock
}

func MockClient(cfgs ...*aws.Config) *client.Client {
	/* gets fixture from testdata folder */
	parentDir, _ := helpers.GetTestDir()
	credPath := fmt.Sprintf("%s%s", parentDir, "/credentials")

	mockServer := MockServer()
	sess := session.Session{Config: &aws.Config{
		Endpoint: aws.String(mockServer.URL),
		Region:   aws.String("us-east-1"),
		/* name of test profile */
		Credentials: credentials.NewSharedCredentials(credPath, "testing")}}
	c := sess.ClientConfig("Mock", cfgs...)

	svc := client.New(
		*c.Config,
		metadata.ClientInfo{
			ServiceName:   "Mock",
			SigningRegion: c.SigningRegion,
			Endpoint:      c.Endpoint,
			APIVersion:    "2015-12-08",
			JSONVersion:   "1.1",
			TargetPrefix:  "MockServer",
		},
		c.Handlers,
	)
	return svc
}

func S3Mock() (*s3.S3, *gnomock.Container) {
	///* gets fixture from testdata folder */
	parentDir, _ := helpers.GetTestDir()
	credPath := fmt.Sprintf("%s%s", parentDir, "/credentials")
	s3Dir := fmt.Sprintf("%s%s", parentDir, "/s3")

	p := localstack.Preset(
		localstack.WithServices(localstack.S3),
		localstack.WithS3Files(s3Dir),
		localstack.WithVersion("latest"),
	)
	log.Info.Println(p.Image())
	c, err := gnomock.Start(p)

	if err != nil {
		log.Error.Panic(err)
	}
	s3Endpoint := fmt.Sprintf("http://%s/", c.Address(localstack.APIPort))
	log.Info.Printf("endpoint: %s", s3Endpoint)
	conf := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String(s3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewSharedCredentials(credPath, "testing"),
	}

	sess, serr := session.NewSession(conf)
	if serr != nil {
		log.Error.Panic(serr)
	}

	return s3.New(sess), c

}

func S3MockDocker() *s3.S3 {
	/* uses localstack container in docker compose, container needs to be up first for the hardcoded endpoint to be active */
	/* gets fixture from testdata folder */
	parentDir, _ := helpers.GetTestDir()
	credPath := fmt.Sprintf("%s%s", parentDir, "/credentials")
	conf := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("http://localstack:4566"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewSharedCredentials(credPath, "testing"),
	}

	sess, serr := session.NewSession(conf)
	if serr != nil {
		log.Error.Panic(serr)
	}

	return s3.New(sess)

}
