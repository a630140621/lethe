package generate

import (
	"errors"
	"fmt"
)

type format string

var (
	formatJSON format = "json"
	formatYAML format = "yaml"
)

func (f format) Is(t format) bool {
	return f == t
}

func (f format) String() string {
	if string(f) == "" {
		// default is yaml
		return string(formatYAML)
	}

	return string(f)
}

func (f *format) Set(v string) error {
	switch format(v) {
	case formatJSON, formatYAML:
		conv := format(v)
		*f = conv
		return nil
	default:
		return errors.New(f.Usage())
	}
}

func (f *format) Type() string {
	return "string"
}

func (f *format) Usage() string {
	return fmt.Sprintf("format should be one of %v, default is yaml", []format{formatJSON, formatYAML})
}
