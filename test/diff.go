package test

import (
	"errors"
	"github.com/go-test/deep"
)

func Diff(a interface{}, b interface{}) error {
	var errorString string
	if diff := deep.Equal(a, b); diff != nil {
		for _, d := range diff {
			errorString += d + "\n"
		}
	}
	if errorString != "" {
		return errors.New(errorString)
	}
	return nil
}
