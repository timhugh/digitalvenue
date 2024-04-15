package dv_aws

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

func GetImageAttributes(itemType string, image map[string]events.DynamoDBAttributeValue, attrNames ...string) (map[string]string, error) {
	missing := make([]string, 0)

	if image == nil {
		return nil, errors.New("image is nil")
	}

	if itemTypeAttr, ok := image["Type"]; !ok || itemTypeAttr.String() != itemType {
		return nil, errors.Errorf("image is not a %s", itemType)
	}

	attrs := make(map[string]string)
	for _, attrName := range attrNames {
		attr, ok := image[attrName]
		if !ok {
			missing = append(missing, attrName)
			continue
		}
		attrs[attrName] = attr.String()
	}

	if len(missing) > 0 {
		return nil, errors.Errorf("missing attributes: %v", missing)
	}

	return attrs, nil
}
