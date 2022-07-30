package store

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3StorageOptions struct {
	AccessKeyId     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Prefix          string
}

type S3Storage struct {
	s3     *s3.S3
	bucket string
	prefix string
}

func NewS3Storage(options S3StorageOptions) *S3Storage {
	config := aws.Config{}

	if options.AccessKeyId != "" {
		config.Credentials = credentials.NewStaticCredentials(
			options.AccessKeyId,
			options.SecretAccessKey,
			"",
		)
	}

	if options.Region != "" {
		config.Region = aws.String(options.Region)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigDisable,
		Config:            config,
	}))

	return &S3Storage{
		s3:     s3.New(sess),
		bucket: options.Bucket,
		prefix: options.Prefix,
	}
}

func (s *S3Storage) buildKey(key string) *string {
	return aws.String(strings.TrimSuffix(s.prefix, "/") + "/" + strings.TrimPrefix(key, "/"))
}

func (s *S3Storage) put(filename string, key string) (*s3.PutObjectOutput, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}

	return s.s3.PutObject(&s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    s.buildKey(key),
		Body:   file,
	})
}

func (s *S3Storage) get(filename string, key string) error {
	o, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    s.buildKey(key),
	})

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(o.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, buf.Bytes(), 0644)
}

func (s *S3Storage) Name() string {
	return "S3"
}

func (s *S3Storage) Store(secrets []string) error {
	for _, secret := range secrets {
		fmt.Printf("%-30s -> s3://%s/%s\n", secret, s.bucket, *s.buildKey(secret))
		_, err := s.put(secret, secret)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *S3Storage) Load(secrets []string) error {
	for _, secret := range secrets {
		fmt.Printf("%-50s -> %-30s \n", fmt.Sprintf("s3://%s/%s", s.bucket, *s.buildKey(secret)), secret)
		err := s.get(secret, secret)
		if err != nil {
			return err
		}
	}
	return nil
}
