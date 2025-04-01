package optional_test

import (
	"testing"

	"github.com/timhugh/digitalvenue/util/optional"
)

func TestOptional(t *testing.T) {
	t.Run("holds values", func(t *testing.T) {
		message := "Hello!"
		opt := optional.New(message)
		value, err := opt.Get()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if value != message {
			t.Errorf("Expected %v, got %v", message, value)
		}
	})

	t.Run("can be set", func(t *testing.T) {
		message := "Hello!"
		opt := optional.Empty[string]()
		opt.Set(message)
		value, err := opt.Get()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if value != message {
			t.Errorf("Expected %v, got %v", message, value)
		}
	})

	t.Run("can be empty", func(t *testing.T) {
		opt := optional.Empty[string]()
		if !opt.Empty() {
			t.Errorf("Expected empty optional, got %v", opt)
		}
		if opt.HasValue() {
			t.Errorf("Expected empty optional, got %v", opt)
		}
		_, err := opt.Get()
		if err == nil {
			t.Error("Expected error, got none")
		}
	})

	t.Run("can be reset", func(t *testing.T) {
		opt := optional.New("Hello!")
		opt.Reset()
		if !opt.Empty() {
			t.Errorf("Expected empty optional, got %v", opt)
		}
		if opt.HasValue() {
			t.Errorf("Expected empty optional, got %v", opt)
		}
		_, err := opt.Get()
		if err == nil {
			t.Error("Expected error, got none")
		}
	})

	t.Run("chaining defaults", func(t *testing.T) {
		optWithValue := optional.New("foo")
		emptyOpt := optional.Empty[string]()

		if optWithValue.OrElse("bar") != "foo" {
			t.Error("Expected present value 'foo', got else value 'bar'")
		}

		emptyValue := emptyOpt.OrElse("bar")
		if emptyValue != "bar" {
			t.Errorf("Expected else value 'bar', got %s", emptyValue)
		}
	})
}
