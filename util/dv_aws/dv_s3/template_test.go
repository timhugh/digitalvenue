package dv_s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/util/test"
	"io"
	"os"
	"testing"
)

const (
	templateBucketName = "test-template-bucket"
)

func initTemplateTest(t *testing.T) (*TemplateStore, Client) {
	mock.SetUp(t)
	client := mock.Mock[Client]()

	err := os.Setenv("S3_TEMPLATES_BUCKET_NAME", templateBucketName)
	if err != nil {
		t.Fatal(err)
	}

	store, err := NewS3TemplateStore(client)
	if err != nil {
		t.Fatal(err)
	}

	return store, client
}

func TestS3TemplateStore_Get(t *testing.T) {
	store, client := initTemplateTest(t)

	getObjectOutput := &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader([]byte(test.TemplateBody))),
	}
	getObjectInputCaptor := mock.Captor[*s3.GetObjectInput]()
	mock.When(client.GetObject(mock.AnyContext(), getObjectInputCaptor.Capture())).ThenReturn(getObjectOutput, nil)

	template, err := store.Get(test.TenantID, test.TemplateKey)
	if err != nil {
		t.Fatal(err)
	}

	expectedInput := &s3.GetObjectInput{
		Bucket: aws.String(templateBucketName),
		Key:    aws.String(test.TenantID + "/" + test.TemplateKey),
	}
	if err := test.Diff(expectedInput, getObjectInputCaptor.Last()); err != nil {
		t.Error(err)
	}

	if err := test.Diff(test.NewTemplate(), template); err != nil {
		t.Error(err)
	}
}

func TestS3TemplateStore_Save(t *testing.T) {
	store, client := initTemplateTest(t)

	putObjectInputCaptor := mock.Captor[*s3.PutObjectInput]()
	mock.When(client.PutObject(mock.AnyContext(), putObjectInputCaptor.Capture())).ThenReturn(nil, nil)

	err := store.Save(test.NewTemplate())
	if err != nil {
		t.Fatal(err)
	}

	expectedInput := &s3.PutObjectInput{
		Bucket: aws.String(templateBucketName),
		Key:    aws.String(test.TenantID + "/" + test.TemplateKey),
		Body:   bytes.NewReader([]byte(test.TemplateBody)),
	}
	if err := test.Diff(expectedInput, putObjectInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestS3TemplateStore_GetsBucketNameFromEnv(t *testing.T) {
	_, client := initTemplateTest(t)
	is := is.New(t)

	err := os.Setenv("S3_TEMPLATES_BUCKET_NAME", templateBucketName)
	is.NoErr(err)

	_, err = NewS3TemplateStore(client)
	is.NoErr(err)

	err = os.Unsetenv("S3_TEMPLATES_BUCKET_NAME")
	is.NoErr(err)

	_, err = NewS3TemplateStore(client)
	if err == nil {
		t.Error("expected error")
	}
}
