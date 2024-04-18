package dv_s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/timhugh/digitalvenue/util/core"
)

const (
	qrCodeBucketNameKey = "S3_QR_CODE_BUCKET_NAME"

	contentTypeFormat = "image/%s"
	bucketURLFormat   = "https://%s.s3.amazonaws.com/%s"
)

var (
	contentEncoding = "base64"
)

type S3QRStorage struct {
	s3Client Client
	bucket   string
}

func NewS3QRStorage(s3Client Client) (*S3QRStorage, error) {
	qrBucket, err := core.RequireEnv(qrCodeBucketNameKey)
	if err != nil {
		return nil, err
	}

	return &S3QRStorage{
		s3Client: s3Client,
		bucket:   qrBucket,
	}, nil
}

func (s *S3QRStorage) Save(qrCode *core.QRCode) (string, error) {
	contentType := fmt.Sprintf(contentTypeFormat, qrCode.FileType)
	key := fmt.Sprintf("%s/%s/%s.png", qrCode.TenantID, qrCode.OrderID, qrCode.OrderItemID)
	putObjectInput := s3.PutObjectInput{
		Bucket:          &s.bucket,
		Key:             &key,
		Body:            bytes.NewBuffer(qrCode.Image),
		ACL:             types.ObjectCannedACLPublicRead,
		ContentEncoding: &contentEncoding,
		ContentType:     &contentType,
	}

	_, err := s.s3Client.PutObject(context.Background(), &putObjectInput)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(bucketURLFormat, s.bucket, key), nil
}
