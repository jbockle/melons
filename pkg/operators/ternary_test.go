package operators

import (
	"testing"
)

func TestIf(t *testing.T) {
	var thenValue int = 1
	var elseValue int = 2

	t.Run("IfTrue", func(t *testing.T) {
		result := If(true, thenValue, elseValue)

		if result != 1 {
			t.Errorf("If() = %v, want %v", result, thenValue)
		}
	})

	t.Run("IfFalse", func(t *testing.T) {
		result := If(false, thenValue, elseValue)

		if result != 2 {
			t.Errorf("If() = %v, want %v", result, elseValue)
		}
	})
}
