package tower

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestF16Neg(t *testing.T) {
	for _ = range 100 {
		x := RandomF16()
		assert.True(t, NewF16().Set(x).Neg().Add(x).Equal(F16Zero()))
	}
}

func TestF16Inv(t *testing.T) {
	assert.Panics(t, func() {
		NewF16().Set(F16Zero()).Inv()
	})

	for _ = range 100 {
		x := RandomF16()
		if x.Equal(F16Zero()) {
			continue
		}

		assert.True(t, NewF16().Set(x).Inv().Mul(x).Equal(F16One()))
	}
}
