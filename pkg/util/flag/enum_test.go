package flag

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func ExampleStringEnum() {
	// 应该提到外部作用域, 为了测试所以放在内部
	var format = NewStringEnum("yaml", "json").WithDefault("json")
	pflag.Var(&format, "format", format.Usage())
}

func TestStringEnum(t *testing.T) {
	t.Run("without default value", func(t *testing.T) {
		format := NewStringEnum("yaml", "json")
		assert.Nil(t, format.d)
		assert.Nil(t, format.value)
		assert.Equal(t, format.String(), "")
	})

	t.Run("with default value", func(t *testing.T) {
		format := NewStringEnum("yaml", "json").WithDefault("json")
		assert.Equal(t, *format.d, "json")
		assert.Nil(t, format.value)
		assert.Equal(t, format.String(), "json")
	})

	t.Run("set value", func(t *testing.T) {
		format := NewStringEnum("yaml", "json").WithDefault("json")
		assert.Error(t, format.Set("not one of enum"))
		assert.NoError(t, format.Set("json"))
		assert.Equal(t, format.String(), "json")
		assert.NoError(t, format.Set("yaml"))
		assert.Equal(t, format.String(), "yaml")
	})
}
