package storage

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pradeepbgs/envy/internal/config"
)


type R2 struct {
	client *s3.Client
	bucket string
}

func New(cfg *config.Config) *R2 {
	c := s3.New(s3.Options{
		BaseEndpoint: aws.String(cfg.R2Endpoint),
		Region: "auto",
		Credentials: credentials.NewStaticCredentialsProvider(cfg.AccessKey,cfg.SecretKey,""),
	})

	return &R2{
		client: c,
		bucket: cfg.Bucket,
	}
}

func (r *R2) Upload (ctx context.Context, key string, data []byte) error {
	_,err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key: aws.String(key),
		Body: bytes.NewReader(data),
	})
	return err
}

func (r *R2) Download (ctx context.Context, key string) ([]byte, error) {
	out , err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key: aws.String(key),
	})
	if err != nil {
		return nil,err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}

func (r *R2) List(ctx context.Context) ([]string, error) {
	out,err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &r.bucket,
	})
	if err != nil {
		return nil,err
	}

	var keys []string
	
	for _, obj := range out.Contents {
		if strings.HasSuffix(*obj.Key, ".enc") {
			keys = append(keys, *obj.Key)
		}
	}

	return keys, nil
}

func (r *R2) Delete(ctx context.Context, key string) error {
	_,err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &r.bucket,
		Key: aws.String(key),
	})
	return err
}