package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"testing"
)

func TestNewClient(t *testing.T) {
	var _ Client = NewClient(aws.Config{})
}
