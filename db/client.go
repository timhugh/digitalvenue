package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

var (
	ErrInitFailed    = errors.New("initialization failed")
	ErrNoResults     = errors.New("no results found")
	ErrQueryFailed   = errors.New("query failed")
	ErrMappingFailed = errors.New("mapping failed")
)

type Request struct {
	Query      string
	Args       []any
	ResultChan chan Result
}

type Result struct {
	Data  []map[string]any
	Error error
}

func (r Result) Unwrap(out any) error {
	if r.Error != nil {
		return r.Error
	}

	if len(r.Data) == 0 {
		return fmt.Errorf("%w", ErrNoResults)
	}

	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot unwrap non-pointer type %w", ErrMappingFailed)
	}

	elemVal := outVal.Elem()
	switch elemVal.Kind() {
	case reflect.Struct:
		return mapstructure.Decode(r.Data[0], &out)
	case reflect.Slice:
		return mapstructure.Decode(r.Data, &out)
	default:
		return fmt.Errorf("Result.Unwrap does not support type %s %w", elemVal.Kind(), ErrMappingFailed)
	}
}

type Client interface {
	ExecuteQuery(ctx context.Context, query string, args ...any) Result
}
