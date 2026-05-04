package s3

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client *s3.Client
	bucket string
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
