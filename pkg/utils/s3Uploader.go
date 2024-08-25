package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/structs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToS3(session *session.Session, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileName := filepath.Base(filePath)

	key := fmt.Sprintf("%s/%s", os.Getenv("BUCKET_ENDPOINT_URI"), fileName)

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("BUCKET_NAME")),
		Key:                  aws.String(key),
		Body:                 file,
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func GeneratePresignedPost(session *session.Session, category structs.UploadCategory) (string, error) {
	svc := s3.New(session)
	bucket := os.Getenv("BUCKET_NAME")
	key := fmt.Sprintf("%s/%s/%d", os.Getenv("BUCKET_ENDPOINT_URI"), category, time.Now().Unix())

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %v", err)
	}

	return urlStr, nil
}

// checks if the given URL is a valid S3 URL and the object exists
func ValidateS3URL(sess *session.Session, s3Url string) error {
	// Retrieve environment variables
	bucketName := os.Getenv("BUCKET_NAME")
	bucketEndpointIP := os.Getenv("BUCKET_ENDPOINT_IP")
	bucketEndpointPort := os.Getenv("BUCKET_ENDPOINT_PORT")
	bucketEndpointURI := os.Getenv("BUCKET_ENDPOINT_URI")
	useSSL := os.Getenv("BUCKET_ENDPOINT_USE_SSL")

	// Construct the expected base URL
	protocol := "http"
	if useSSL == "true" {
		protocol = "https"
	}
	expectedBaseURL := fmt.Sprintf("%s://%s:%s/%s/%s", protocol, bucketEndpointIP, bucketEndpointPort, bucketName, bucketEndpointURI)

	// Parse the provided URL
	parsedURL, err := url.Parse(s3Url)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %v", err)
	}

	// Validate the base URL
	if !strings.HasPrefix(s3Url, expectedBaseURL) {
		return fmt.Errorf("URL does not match the expected S3 endpoint: %s", expectedBaseURL)
	}

	// Extract and decode the key from the URL path
	path := strings.TrimPrefix(parsedURL.Path, fmt.Sprintf("/%s/%s/", bucketName, bucketEndpointURI))
	key, err := url.QueryUnescape(path)
	if err != nil {
		return fmt.Errorf("failed to decode URL path: %v", err)
	}

	// Debugging: Print the derived key
	// fmt.Printf("Derived key: '%s'\n", key)

	// Verify the object exists
	svc := s3.New(sess)
	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("uploads/" + key),
	})

	if err != nil {
		return fmt.Errorf("failed to get object: %v", err)
	}

	return nil
}