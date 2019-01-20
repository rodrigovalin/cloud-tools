package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const AwsRegion = "us-east-1"

func isS3File(filename string) bool {
	return strings.HasPrefix(filename, "https://s3.amazonaws.com")
}

func writeS3File(filename string, filecontents []byte) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return err
	}

	bucket, filename := explodeS3URL(filename)

	retry := 3

	for retry > 0 {
		fmt.Printf("Trying to push %s into bucket %s\n", filename, bucket)

		uploader := s3manager.NewUploader(sess)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
			Body:   bytes.NewReader(filecontents),
		})
		if err != nil {
			fmt.Println(err)
			retry -= 1
		} else {
			break
		}
	}

	return err
}

// explodeS3URL returns bucket and filename from S3 URL.
func explodeS3URL(s3url string) (string, string) {
	u, err := url.Parse(s3url)
	if err != nil {
		return "", ""
	}

	splittedBucket := strings.SplitN(u.Path, "/", 2)
	return splittedBucket[0], splittedBucket[1]
}

func listVersionsForFile(filename string) (versions []string, err error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion),
	})
	if err != nil {
		return
	}

	bucket, filename := explodeS3URL(filename)

	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(filename),
	}

	svc := s3.New(sess)
	result, err := svc.ListObjectVersions(input)
	if err != nil {
		return
	}

	for _, version := range result.Versions {
		versions = append(versions, s3VersionToString(version))
	}

	return
}

func s3VersionToString(version *s3.ObjectVersion) string {
	result := version.LastModified.String()
	if *version.IsLatest {
		result = fmt.Sprintf("%s (latest)", result)
	}

	return result
}
