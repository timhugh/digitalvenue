package util

import "fmt"

func NewOptional[Value any](value Value) Optional[Value] {
	return Optional[Value]{
		value:    value,
		hasValue: true,
	}
}

func EmptyOptional[Value any]() Optional[Value] {
	return Optional[Value]{
		hasValue: false,
	}
}

type Optional[Value any] struct {
	value    Value
	hasValue bool
}

func (o Optional[Value]) HasValue() bool {
	return o.hasValue
}

func (o Optional[Value]) Empty() bool {
	return !o.HasValue()
}

func (o Optional[Value]) Get() (Value, error) {
	if o.Empty() {
		return o.value, fmt.Errorf("Optional has no value")
	}
	return o.value, nil
}

func (o *Optional[Value]) Set(value Value) {
	o.value = value
	o.hasValue = true
}

func (o *Optional[Value]) Reset() {
	o.hasValue = false
}

func (o Optional[Value]) OrElse(elseValue Value) Value {
	if o.hasValue {
		return o.value
	}
	return elseValue
}
