package operators

import (
	"testing"
)

func TestNilOrZeroCoalesce(t *testing.T) {
	type foo struct {
		bar string
	}

	t.Run("NilOrZeroCoalesce struct", func(t *testing.T) {
		var zeroValue foo
		var oneValue foo = foo{bar: "bar"}

		result := Coalesce(zeroValue, oneValue)

		if result != oneValue {
			t.Errorf("NilOrZeroCoalesce() = %v, want %v", result, oneValue)
		}
	})

	t.Run("NilOrZeroCoalesce zeroed", func(t *testing.T) {
		var zeroValue int
		var oneValue int = 1

		result := Coalesce(zeroValue, oneValue)

		if result != 1 {
			t.Errorf("NilOrZeroCoalesce() = %v, want %v", result, 1)
		}
	})

	t.Run("NilOrZeroCoalesce zeroed All Nil", func(t *testing.T) {
		var zeroValue int

		result := Coalesce(zeroValue, zeroValue)

		if result != 0 {
			t.Errorf("NilOrZeroCoalesce() = %v, want %v", result, 0)
		}
	})
}
