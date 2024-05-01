package dv_s3

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"io"
	"path"
)

type TemplateStore struct {
	bucketName string
	s3Client   Client
}

func NewS3TemplateStore(
	s3Client Client,
) (*TemplateStore, error) {
	bucketName, err := core.RequireEnv("S3_TEMPLATES_BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	return &TemplateStore{
		bucketName: bucketName,
		s3Client:   s3Client,
	}, nil
}

func (s *TemplateStore) Save(template *core.Template) error {
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(buildFullTemplateKey(template.TenantID, template.Key)),
		Body:   bytes.NewReader([]byte(template.Body)),
	}

	_, err := s.s3Client.PutObject(context.TODO(), putObjectInput)
	if err != nil {
		return errors.Wrap(err, "failed to put template to s3")
	}

	return nil
}

func (s *TemplateStore) Get(tenantID string, templateKey string) (*core.Template, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(buildFullTemplateKey(tenantID, templateKey)),
	}

	getObjectOutput, err := s.s3Client.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve template object from s3")
	}
	defer getObjectOutput.Body.Close()

	body, err := io.ReadAll(getObjectOutput.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read template body")
	}

	return &core.Template{
		TenantID: tenantID,
		Key:      templateKey,
		Body:     string(body),
	}, nil
}

func buildFullTemplateKey(tenantID string, templateKey string) string {
	return path.Join(tenantID, templateKey)
}
