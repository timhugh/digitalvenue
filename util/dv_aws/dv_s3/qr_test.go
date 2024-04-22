package dv_s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/util/test"
	"os"
	"testing"
)

func initTest(t *testing.T) *is.I {
	is := is.New(t)
	mock.SetUp(t)

	err := os.Setenv("S3_QR_CODE_BUCKET_NAME", test.QRCodeBucket)
	is.NoErr(err)

	return is
}

func TestS3QRStorage_Save(t *testing.T) {
	is := initTest(t)

	client := mock.Mock[Client]()
	putObjectInputCaptor := mock.Captor[*s3.PutObjectInput]()
	mock.WhenDouble(client.PutObject(mock.Any[context.Context](), putObjectInputCaptor.Capture())).ThenReturn(nil, nil)

	storage, err := NewS3QRStorage(client)
	is.NoErr(err)

	qr := test.NewQRCode()
	url, err := storage.Save(qr)
	is.NoErr(err)
	is.Equal(fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s/%s.png", test.QRCodeBucket, qr.TenantID, qr.OrderID, qr.OrderItemID), url)

	bucket := test.QRCodeBucket
	key := fmt.Sprintf("%s/%s/%s.png", qr.TenantID, qr.OrderID, qr.OrderItemID)
	contentType := "image/" + qr.FileType
	expectedPutObjectInput := &s3.PutObjectInput{
		Bucket:          &bucket,
		Key:             &key,
		Body:            bytes.NewReader(qr.Image),
		ACL:             types.ObjectCannedACLPublicRead,
		ContentEncoding: &contentEncoding,
		ContentType:     &contentType,
	}
	if err := test.Diff(expectedPutObjectInput, putObjectInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}
