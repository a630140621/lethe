package flag

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
)

type StringEnum struct {
	value   *string
	allowed []string
	d       *string
}

var _ pflag.Value = &StringEnum{}

func NewStringEnum(values ...string) StringEnum {
	return StringEnum{
		allowed: values,
		value:   nil,
		d:       nil,
	}
}

func (enum *StringEnum) validate(t string) error {
	for _, v := range enum.allowed {
		if t == v {
			return nil
		}
	}

	return errors.New(enum.Usage())
}

func (enum StringEnum) WithDefault(d string) StringEnum {
	enum.d = &d
	return enum
}

func (enum StringEnum) FlagValue() pflag.Value {
	return &enum
}

func (enum *StringEnum) String() string {
	if enum.value != nil {
		return *enum.value
	}

	if enum.d != nil {
		return *enum.d
	}

	return ""
}

func (enum *StringEnum) Set(v string) error {
	if err := enum.validate(v); err != nil {
		return err
	}

	enum.value = &v
	return nil
}

func (enum *StringEnum) Type() string {
	return "enum"
}

func (enum *StringEnum) Usage() string {
	return fmt.Sprintf("should be one of %v", enum.allowed)
}

func (enum *StringEnum) Is(t string) bool {
	return enum.String() == t
}
