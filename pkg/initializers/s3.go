package initializers

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var S3Session *session.Session

func InitializeS3() {
	endpoint := os.Getenv("BUCKET_ENDPOINT_IP") + ":" + os.Getenv("BUCKET_ENDPOINT_PORT")
	useSSL := os.Getenv("BUCKET_ENDPOINT_USE_SSL") == "true"

	S3Session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Endpoint: aws.String((func() string {
			if useSSL {
				return "https://" + endpoint
			}
			return "http://" + endpoint
		})()),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("BUCKET_ACCESS_KEY"),
			os.Getenv("BUCKET_SECRET_KEY"),
			"",
		),
		S3ForcePathStyle: aws.Bool(true),
	}))
}
