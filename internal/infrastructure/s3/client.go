// internal/infrastructure/s3/client.go
package s3

import (
	"DataLake/pkg"
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client *s3.Client
	bucket string
}

func New() *Client {
	// Ładowanie konfiguracji AWS z zmiennych środowiskowych
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(pkg.GetEnv("S3_REGION", "us-east-1")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			pkg.GetEnv("S3_ACCESS_KEY", ""),
			pkg.GetEnv("S3_SECRET_KEY", ""),
			"", // Session token
		)),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: pkg.GetEnv("S3_ENDPOINT", ""),
			}, nil
		})),
	)
	if err != nil {
		panic("Błąd konfiguracji S3: " + err.Error())
	}

	return &Client{
		client: s3.NewFromConfig(cfg),
		bucket: pkg.GetEnv("S3_BUCKET", ""),
	}
}

func (c *Client) Put(key string, data []byte) error {
	_, err := c.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	return err
}

func (c *Client) Get(key string) ([]byte, error) {
	out, err := c.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	defer out.Body.Close()

	return io.ReadAll(out.Body)
}

func (c *Client) Exists(key string) (bool, error) {
	_, err := c.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (c *Client) List(prefix string) ([]string, error) {
	out, err := c.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		return nil, err
	}

	var keys []string

	for _, obj := range out.Contents {
		keys = append(keys, *obj.Key)
	}

	return keys, nil
}

func (c *Client) Delete(key string) error {
	_, err := c.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return err
}
