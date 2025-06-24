package tower

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestF64Neg(t *testing.T) {
	for _ = range 100 {
		x := RandomF64()
		assert.True(t, NewF64().Set(x).Neg().Add(x).Equal(F64Zero()))
	}
}

func TestF64Inv(t *testing.T) {
	assert.Panics(t, func() {
		NewF64().Set(F64Zero()).Inv()
	})

	for _ = range 100 {
		x := RandomF64()
		if x.Equal(F64Zero()) {
			continue
		}

		assert.True(t, NewF64().Set(x).Inv().Mul(x).Equal(F64One()))
	}
}
